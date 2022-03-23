package tables

import (
	"context"
	"encoding/json"
	redis2 "github.com/quanxiang-cloud/cabin/tailormade/db/redis"
	"github.com/quanxiang-cloud/form/internal/models"
	"github.com/quanxiang-cloud/form/internal/models/redis"
	"github.com/quanxiang-cloud/form/internal/service/types"
	"github.com/quanxiang-cloud/form/pkg/misc/config"
	"github.com/quanxiang-cloud/form/pkg/misc/utils"
	"reflect"
)

const (
	_properties     = "properties"
	xComponent      = "x-component"
	xComponentProps = "x-component-props"
	items           = "items"
)

type component struct {
	tableRelationRepo models.TableRelationRepo
	next              Guidance
	serialRepo        models.SerialRepo
}

func (c *component) Do(ctx context.Context, bus *Bus) (*DoResponse, error) {
	i := bus.Schema[_properties]

	asMap, err := getAsMap(i)
	if err != nil {
		return nil, err
	}
	c.subDo(ctx, asMap, &base{
		appID:   bus.AppID,
		tableID: bus.TableID,
	})
	return c.next.Do(ctx, bus)
}

type base struct {
	appID     string
	tableID   string
	fieldName string

	fieldValue types.M
	components string
}

func (c *component) subDo(ctx context.Context, properties types.M, bus *base) {
	// 判断是否是 数据组件

	for fieldName, fieldValue := range properties {
		isLayout := isLayoutComponent(fieldValue)
		// 判断是否是布局组件
		if isLayout {
			//
			v, err := getAsMap(fieldValue)
			if err != nil {
				continue
			}

			toMap, err := getMapToMap(v, _properties)
			if err != nil {
				continue
			}
			c.subDo(ctx, toMap, bus)
		}
		asMap, err := getAsMap(fieldValue)
		if err != nil {
			continue
		}
		components := getMapToString(asMap, xComponent)
		bus.components = components
		bus.fieldName = fieldName
		bus.fieldValue = asMap
		if components == "Serial" {
			c.doSerial(ctx, bus)
		}
		if components == "SubTable" {

			c.doRelation(ctx, bus)
		}
	}
}

func (c *component) doRelation(ctx context.Context, bus *base) error {

	// TODO  add repo

	//
	cp := &ComponentProp{}
	c1, ok := bus.fieldValue[xComponentProps]
	if !ok {
		return nil
	}
	err := genComponent(c1, cp)
	if err != nil {
		return err
	}
	// 解决item  判断 子表单 是否带有流水号 递归
	toMap, err := getMapToMap(bus.fieldValue, items)
	if err != nil {
		return nil
	}
	mapToMap, err := getMapToMap(toMap, _properties)
	if err != nil {
		return nil
	}

	bases := &base{
		appID:   cp.AppID,
		tableID: cp.TableID,
	}
	c.subDo(ctx, mapToMap, bases)
	return nil
}

func (c *component) doSerial(ctx context.Context, bus *base) error {
	fieldName := bus.fieldName
	cp := &ComponentProp{}

	c1, ok := bus.fieldValue[xComponentProps]
	if !ok {
		return nil
	}
	err := genComponent(c1, cp)
	if err != nil {
		return err
	}

	// 解析模版
	serial, template := utils.ParseTemplate(cp.Template)
	serialData, err := json.Marshal(serial)
	if err != nil {
		return err
	}
	// 判断流水号是否第一次创建
	oldSerialStr := c.serialRepo.Get(ctx, bus.appID, bus.tableID, fieldName, redis.Serials)
	if oldSerialStr == "" {
		if err := c.serialRepo.Create(ctx, bus.appID, bus.tableID, fieldName, map[string]interface{}{
			redis.Serials:  serialData,
			redis.Template: template,
		}); err != nil {
			return err
		}
		return nil
	}

	// 判断是否修改初始位
	if err = utils.CheckSerial(&serial, oldSerialStr); err != nil {
		return err
	}

	if serialData, err = json.Marshal(serial); err != nil {
		return err
	}
	if err := c.serialRepo.Create(ctx, bus.tableID, bus.tableID, fieldName, map[string]interface{}{
		redis.Serials:  serialData,
		redis.Template: template,
	}); err != nil {
		return err
	}

	return nil

}

// ComponentProp schema中 x-component-props 结构
type ComponentProp struct {
	AppID   string   `json:"appID"`
	TableID string   `json:"tableID"`
	Columns []string `json:"columns"`
	// 'sub_table | foreign_table'
	Subordination string `json:"subordination"`
	// 关联表schema
	AssociatedTable interface{}            `json:"associatedTable"`
	Multiple        bool                   `json:"multiple"`
	FieldName       string                 `json:"fieldName"`
	AggType         string                 `json:"aggType"`
	Conditions      map[string]interface{} `json:"condition"`
	FilterConfig    map[string]interface{} `json:"filterConfig"`
	Template        string                 `json:"template"`
}

func genComponent(c interface{}, cp *ComponentProp) error {
	cb, cbErr := json.Marshal(c)
	if cbErr != nil {
		return cbErr
	}
	err := json.Unmarshal(cb, cp)
	if err != nil {
		return err
	}
	return nil
}

func isLayoutComponent(value interface{}) bool {
	switch reflect.TypeOf(value).Kind() {
	case reflect.Map:
		v := reflect.ValueOf(value)
		if value := v.MapIndex(reflect.ValueOf("x-internal")); value.IsValid() {
			if value.CanInterface() {
				return isLayoutComponent(value.Interface())
			}
		}
		if value := v.MapIndex(reflect.ValueOf("isLayoutComponent")); value.IsValid() {
			if _, ok := value.Interface().(bool); ok {
				return value.Interface().(bool)
			}
		}
	default:
		return false
	}
	return false
}

func newComponent(conf *config.Config) (Guidance, error) {
	swagger, err := newRegisterSwagger(conf)
	if err != nil {
		return nil, err
	}
	redisClient, err := redis2.NewClient(conf.Redis)
	if err != nil {
		return nil, err
	}

	return &component{
		next:       swagger,
		serialRepo: redis.NewSerialRepo(redisClient),
	}, nil
}
