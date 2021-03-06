package util

import (
	"fmt"
	"reflect"

	"github.com/go-openapi/spec"
	error2 "github.com/quanxiang-cloud/cabin/error"
	"github.com/quanxiang-cloud/form/internal/models"
	"github.com/quanxiang-cloud/form/pkg/misc/code"
)

const (
	datetime      = "datetime"
	labelValue    = "label-value"
	_id           = "_id"
	_createdAt    = "created_at"
	_creatorID    = "creator_id"
	_creatorName  = "creator_name"
	_updatedAt    = "updated_at"
	_modifierID   = "modifier_id"
	_modifierName = "modifier_name"
)

func GetMapToMap(schema map[string]interface{}, key string) (map[string]interface{}, error) {
	value, ok := schema[key]
	if !ok {
		return nil, error2.New(code.ErrParameter)
	}
	if v, ok := value.(map[string]interface{}); ok {
		return v, nil
	}
	return nil, error2.New(code.ErrParameter)
}

// GetMapToString GetMapToString.
func GetMapToString(schema map[string]interface{}, key string) string {
	value, ok := schema[key]
	if !ok {
		return ""
	}
	if v, ok := value.(string); ok {
		return v
	}
	return ""
}

// GetAsMap GetAsMap.
func GetAsMap(v interface{}) (map[string]interface{}, error) {
	if m, ok := v.(map[string]interface{}); ok {
		return m, nil
	}
	return nil, error2.New(code.ErrParameter)
}

func GetMapToBool(schema map[string]interface{}, key string) (bool, error) {
	value, ok := schema[key]
	if !ok {
		return false, error2.New(code.ErrParameter)
	}
	if v, ok := value.(bool); ok {
		return v, nil
	}
	return false, error2.New(code.ErrParameter)
}

func Convert1(schema map[string]interface{}) (s models.SchemaProperties, total int64, err error) {
	s = make(models.SchemaProperties, 0)
	total = 0

	for key, value := range schema {
		v, err := GetAsMap(value)
		if err != nil {
			continue
		}
		if component, ok := v["x-component"]; ok && (component == "SubTable" || component == "AssociatedRecords") {
			continue
		}

		// 2、 判断是否是布局组件
		isLayout := IsLayoutComponent(value)
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
		schemaProps := models.SchemaProps{}
		for k1, v1 := range v {
			switch k1 {
			case "type":
				schemaProps.Type = v1.(string)
				if v1 == datetime || v1 == labelValue {
					schemaProps.Type = "string"
				}
				if v1 == "array" {
					if _, ok := v["items"]; !ok {
						schemaProps.Items = &models.SchemaProps{
							Type: "string",
						}
						continue
					}
					if _, ok := v["items"].(map[string]interface{}); !ok {
						return nil, 0, error2.New(code.ErrItemConvert)
					}
				}
			case "length":
				t, _ := v1.(int)
				schemaProps.Length = t
			case "title":
				t, _ := v1.(string)
				schemaProps.Title = t
			case "required":
				t, _ := v1.(bool)
				schemaProps.Required = t
			case "properties":
				if p, ok := v1.(map[string]interface{}); ok {
					s2, _, _ := Convert1(p)
					schemaProps.Properties = s2
				}
			case "items":
				if item, ok := v1.(map[string]interface{}); ok {
					types := GetMapToString(item, "type")
					items := &models.SchemaProps{
						Type: types,
					}
					if types == "object" {
						toMap, err := GetMapToMap(item, "properties")
						if err != nil {
							return nil, 0, error2.New(code.ErrItemConvert)
						}
						s2, _, _ := Convert1(toMap)
						items.Properties = s2
					}
					schemaProps.Items = items
				}
			default:
				continue
			}
		}
		s[key] = schemaProps
	}

	return s, total, nil
}

func GetSpecSchema(properties models.SchemaProperties) (spec.SchemaProperties, []string) {
	if properties == nil {
		return nil, nil
	}
	pr := make(spec.SchemaProperties, 0)
	require := make([]string, 0)
	for key, value := range properties {
		if value.Required {
			require = append(require, key)
		}
		d := spec.Schema{
			SchemaProps: spec.SchemaProps{
				Type:  []string{value.Type},
				Title: value.Title,
			},
		}
		if value.Items != nil {
			schema, requires := GetSpecSchema(value.Items.Properties)
			d.Items = &spec.SchemaOrArray{
				Schema: &spec.Schema{
					SchemaProps: spec.SchemaProps{
						Type:       []string{value.Items.Type},
						Properties: schema,
						Required:   requires,
					},
				},
			}
		}
		if value.Properties != nil {
			schema, requires := GetSpecSchema(value.Properties)
			d.SchemaProps.Properties = schema
			d.SchemaProps.Required = requires
		}
		pr[key] = d
	}
	return pr, require
}

func IsLayoutComponent(value interface{}) bool {
	switch reflect.TypeOf(value).Kind() {
	case reflect.Map:
		v := reflect.ValueOf(value)
		if value := v.MapIndex(reflect.ValueOf("x-internal")); value.IsValid() {
			if value.CanInterface() {
				return IsLayoutComponent(value.Interface())
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

func FilterSystem(sourceStruct spec.SchemaProperties, dst spec.SchemaProperties) {
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

func GetSummary(tableName, operate string) string {
	return fmt.Sprintf("%s(%s)", tableName, operate)
}
