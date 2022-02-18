package swagger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"git.internal.yunify.com/qxp/misc/error2"
	"github.com/quanxiang-cloud/form/pkg/misc/code"
	"github.com/quanxiang-cloud/form/pkg/misc/config"
	"reflect"
	"strings"
	"text/template"
)

const (
	datetime      = "datetime"
	labelValue    = "label-value"
	url           = "/api/v1/form/%s/home/form/%s/%s"
	_id           = "_id"
	_createdAt    = "created_at"
	_creatorID    = "creator_id"
	_creatorName  = "creator_name"
	_updatedAt    = "updated_at"
	_modifierID   = "modifier_id"
	_modifierName = "modifier_name"
	_create       = "create"
	_delete       = "delete"
	_update       = "update"
	_search       = "search"
	_get          = "get"
)

// Schema  Schema
type Schema map[string]interface{}

// GenSwagger GenSwagger
func GenSwagger(conf *config.Config, schema Schema, tableName, appID, tableID string) (string, error) {
	template, err := template.ParseFiles(conf.SwaggerPath)
	if err != nil {
		return "", err
	}
	dstSchema := make(map[string]interface{})
	filterSystem(schema, dstSchema)
	schemas, err := json.Marshal(schema)

	if err != nil {
		return "", err
	}
	filterSchemas, err := json.Marshal(dstSchema)
	if err != nil {
		return "", err
	}
	sourceStruct := NewSourceStruct(string(schemas), string(filterSchemas), tableName, appID, tableID)
	buffer := new(bytes.Buffer)
	err = template.Execute(buffer, sourceStruct)
	if err != nil {
		return "", err
	}
	return buffer.String(), nil
}

func filterSystem(sourceStruct Schema, dst Schema) {
	if sourceStruct == nil || dst == nil {
		return
	}
	for key, value := range sourceStruct {
		switch key {
		case _id, _createdAt, _creatorID, _creatorName, _updatedAt, _modifierID, _modifierName:
		default:
			dst[key] = value
		}
	}
}

// NewSourceStruct NewSourceStruct
func NewSourceStruct(properties, filterProperties, tableName, appID, tableID string) *SourceStruct {
	return &SourceStruct{
		Properties:       properties,
		FilterProperties: filterProperties,
		PathCreate:       fmt.Sprintf(url, appID, tableID, _create),
		OperateIDCreate:  GetInnerXName(tableID, _create),
		PathDelete:       fmt.Sprintf(url, appID, tableID, _delete),
		OperateIDDelete:  GetInnerXName(tableID, _delete),
		PathUpdate:       fmt.Sprintf(url, appID, tableID, _update),
		OperateIDUpdate:  GetInnerXName(tableID, _update),
		PathSearch:       fmt.Sprintf(url, appID, tableID, _search),
		OperateIDSearch:  GetInnerXName(tableID, _search),
		PathGet:          fmt.Sprintf(url, appID, tableID, _get),
		OperateIDGet:     GetInnerXName(tableID, _get),
		CreateSummary:    getSummary(tableName, "新增"),
		DeleteSummary:    getSummary(tableName, "删除"),
		SearchSummary:    getSummary(tableName, "查询多条"),
		GetSummary:       getSummary(tableName, "查询单条"),
		UpdateSummary:    getSummary(tableName, "更新"),
	}
}

func getSummary(tableName, operate string) string {
	return fmt.Sprintf("%s(%s)", tableName, operate)
}

//SourceStruct SourceStruct
type SourceStruct struct {
	Properties       string
	FilterProperties string
	PathCreate       string
	OperateIDCreate  string
	CreateSummary    string
	DeleteSummary    string
	SearchSummary    string
	GetSummary       string
	UpdateSummary    string
	PathDelete       string
	OperateIDDelete  string
	PathSearch       string
	OperateIDSearch  string
	PathUpdate       string
	OperateIDUpdate  string
	PathGet          string
	OperateIDGet     string
}

// SwagValue SwagValue
type SwagValue map[string]interface{}

// Convert1  Convert1
func Convert1(schema Schema) (s map[string]interface{}, total int64, err error) {
	s = make(Schema, 0)
	total = 0
	for key, value := range schema {
		if v, ok := value.(map[string]interface{}); ok {
			temp := make(Schema, 0)
			// 1、 判断x-component  是SubTable （子表单），AssociatedData  （关联数据）直接放行 ，
			if component, ok := v["x-component"]; ok && (component == "SubTable" || component == "AssociatedRecords") {
				continue
			}
			// 2、 判断是否是布局组件
			isLayout := isLayoutComponent(value)
			if isLayout {
				if p, ok := v["properties"]; ok {
					if p1, ok := p.(map[string]interface{}); ok {
						s2, t, err := Convert1(p1)
						if err != nil {
							return nil, 0, err
						}
						for key, value := range s2 {
							s[key] = value
						}
						total = t + total
						continue
					}
				}
			}
			switch key {
			case _id, _createdAt, _creatorID, _creatorName, _updatedAt, _modifierID, _modifierName:
			default:
				total = total + 1
			}
			for k1, v1 := range v {
				switch k1 {
				case "type":
					temp[k1] = v1
					if v1 == datetime || v1 == labelValue {
						temp[k1] = "string"
					}
					if v1 == "array" {
						if _, ok := v["items"]; !ok {
							temp["items"] = SwagValue{
								"type": "string",
							}
							continue
						}
						if _, ok := v["items"].(map[string]interface{}); !ok {
							return nil, 0, error2.NewError(code.ErrItemConvert)
						}
					}
				case "length", "title", "not_null":
					temp[k1] = v1
				case "properties":
					if p, ok := v1.(map[string]interface{}); ok {
						s2, _, _ := Convert1(p)
						temp[k1] = s2
					}
				case "items":
					if item, ok := v1.(map[string]interface{}); ok {
						temp[k1] = item
					}
				default:
					continue
				}
			}
			s[key] = temp
		} else {

			return nil, 0, error2.NewError(code.ErrValueConvert)

		}
	}
	return
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

// GenXName GenXName
func GenXName(appID, tableID, tag, content string) string {

	return fmt.Sprintf("/system/app/%s/raw/inner/%s/%s/%s.r", appID, NameSpace, content, GetInnerXName(tableID, tag))
}

// GetInnerXName GetInnerXName
func GetInnerXName(tableID, tag string) string {
	tableIDs := strings.Split(tableID, "_")
	return fmt.Sprintf("%s_%s", tableIDs[len(tableIDs)-1], tag)
}
