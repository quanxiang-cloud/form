package swagger

import (
	"fmt"
	"github.com/go-openapi/spec"
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
	Service       = "form"
)

func filterSystem(sourceStruct spec.SchemaProperties, dst spec.SchemaProperties) {
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

func getSummary(tableName, operate string) string {
	return fmt.Sprintf("%s(%s)", tableName, operate)
}

// Convert1  Convert1
//func Convert1(schema Schema) (s Schema, total int64, err error) {
//	s = make(Schema, 0)
//	total = 0
//	for key, value := range schema {
//		if v, ok := value.(map[string]interface{}); ok {
//			temp := make(Schema, 0)
//			// 1、 判断x-component  是SubTable （子表单），AssociatedData  （关联数据）直接放行 ，
//			if component, ok := v["x-component"]; ok && (component == "SubTable" || component == "AssociatedRecords") {
//				continue
//			}
//			// 2、 判断是否是布局组件
//			isLayout := IsLayoutComponent(value)
//			if isLayout {
//				if p, ok := v["properties"]; ok {
//					if p1, ok := p.(map[string]interface{}); ok {
//						s2, t, err := Convert1(p1)
//						if err != nil {
//							return nil, 0, err
//						}
//						for key, value := range s2 {
//							s[key] = value
//						}
//						total = t + total
//						continue
//					}
//				}
//			}
//			switch key {
//			case _id, _createdAt, _creatorID, _creatorName, _updatedAt, _modifierID, _modifierName:
//			default:
//				total = total + 1
//			}
//			for k1, v1 := range v {
//				switch k1 {
//				case "type":
//					temp[k1] = v1
//					if v1 == datetime || v1 == labelValue {
//						temp[k1] = "string"
//					}
//					if v1 == "array" {
//						if _, ok := v["items"]; !ok {
//							//temp["items"] = SwagValue{
//							//	"type": "string",
//							//}
//							//continue
//						}
//						if _, ok := v["items"].(map[string]interface{}); !ok {
//							return nil, 0, error2.NewError(code.ErrItemConvert)
//						}
//					}
//				case "length", "title", "not_null":
//					temp[k1] = v1
//				case "properties":
//					if p, ok := v1.(map[string]interface{}); ok {
//						s2, _, _ := Convert1(p)
//						temp[k1] = s2
//					}
//				case "items":
//					if item, ok := v1.(map[string]interface{}); ok {
//						temp[k1] = item
//					}
//				default:
//					continue
//				}
//			}
//			s[key] = temp
//		} else {
//
//			return nil, 0, error2.NewError(code.ErrValueConvert)
//
//		}
//	}
//	return
//}
