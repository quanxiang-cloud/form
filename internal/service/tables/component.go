package tables

import (
	"context"
	"encoding/json"
	redis2 "github.com/quanxiang-cloud/cabin/tailormade/db/redis"

	id2 "github.com/quanxiang-cloud/cabin/id"
	"github.com/quanxiang-cloud/form/internal/models"
	"github.com/quanxiang-cloud/form/internal/models/mysql"
	"github.com/quanxiang-cloud/form/internal/models/redis"
	"github.com/quanxiang-cloud/form/internal/service"
	"github.com/quanxiang-cloud/form/internal/service/tables/util"
	"github.com/quanxiang-cloud/form/internal/service/types"
	"github.com/quanxiang-cloud/form/pkg/misc/config"
	"github.com/quanxiang-cloud/form/pkg/misc/utils"
	"gorm.io/gorm"
)

const (
	_properties     = "properties"
	xComponent      = "x-component"
	xComponentProps = "x-component-props"
	items           = "items"
)

type component struct {
	db                *gorm.DB
	tableRelationRepo models.TableRelationRepo
	next              Guidance
	serialRepo        models.SerialRepo
}

func (c *component) Do(ctx context.Context, bus *Bus) (*DoResponse, error) {
	i := bus.Schema[_properties]

	asMap, err := util.GetAsMap(i)
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
		isLayout := util.IsLayoutComponent(fieldValue)
		// 判断是否是布局组件
		if isLayout {
			v, err := util.GetAsMap(fieldValue)
			if err != nil {
				continue
			}

			toMap, err := util.GetMapToMap(v, _properties)
			if err != nil {
				continue
			}
			c.subDo(ctx, toMap, bus)
		}
		asMap, err := util.GetAsMap(fieldValue)
		if err != nil {
			continue
		}
		components := util.GetMapToString(asMap, xComponent)
		bus.components = components
		bus.fieldName = fieldName
		bus.fieldValue = asMap
		if components == "Serial" {
			c.doSerial(ctx, bus)
		}
		if components == "SubTable" || components == "AssociatedRecords" {
			c.doRelation(ctx, bus)
		}
	}
}

func (c *component) doRelation(ctx context.Context, bus *base) error {
	cp := &ComponentProp{}
	c1, ok := bus.fieldValue[xComponentProps]
	if !ok {
		return nil
	}
	err := genComponent(c1, cp)
	if err != nil {
		return err
	}
	tables := &models.TableRelation{
		ID:         id2.StringUUID(),
		AppID:      bus.appID,
		TableID:    bus.tableID,
		FieldName:  bus.fieldName,
		SubTableID: cp.TableID,
		Filter:     cp.Columns,
	}
	tables.SubTableType = cp.Subordination
	if bus.components == "AssociatedRecords" {
		tables.SubTableType = "associated_records"
	}
	err = c.addRepo(tables)
	if err != nil {
		return err
	}

	// 解决item  判断 子表单 是否带有流水号 递归
	toMap, err := util.GetMapToMap(bus.fieldValue, items)
	if err != nil {
		return nil
	}
	mapToMap, err := util.GetMapToMap(toMap, _properties)
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

func (c *component) addRepo(table *models.TableRelation) error {
	relation, err := c.tableRelationRepo.Get(c.db, table.TableID, table.FieldName)
	if err != nil {
		return err
	}
	if relation.ID == "" { // create
		return c.tableRelationRepo.BatchCreate(c.db, table)
	}
	return c.tableRelationRepo.Update(c.db, table.TableID, table.FieldName, table)
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

// ComponentProp schema中 x-component-props 结构.
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

func newComponent(conf *config.Config) (Guidance, error) {
	db, err := service.CreateMysqlConn(conf)
	if err != nil {
		return nil, err
	}
	swagger, err := newRegisterSwagger(conf)
	if err != nil {
		return nil, err
	}
	redisClient, err := redis2.NewClient(conf.Redis)
	if err != nil {
		return nil, err
	}

	return &component{
		db:                db,
		tableRelationRepo: mysql.NewTableRelationRepo(),
		next:              swagger,
		serialRepo:        redis.NewSerialRepo(redisClient),
	}, nil
}
