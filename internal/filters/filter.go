package filters

import (
	"fmt"
	"github.com/quanxiang-cloud/form/internal/models"
	"reflect"
)

// JSONFilter2 json field filter,inputJSON JUST match (map[string]interface,[]map[string]interface)
func JSONFilter2(inputJSON interface{}, requiredFields map[string]interface{}) {
	switch reflect.TypeOf(inputJSON).Kind() {
	case reflect.Map:

		v := reflect.ValueOf(inputJSON)
		iter := v.MapRange()
		for iter.Next() {
			if _, ok := requiredFields[iter.Key().String()]; ok {
				// TODO
				switch reflect.TypeOf(requiredFields[iter.Key().String()]).Kind() {
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Float32, reflect.Float64:
				default:
					JSONFilter2(iter.Value().Interface(), requiredFields[iter.Key().String()].(map[string]interface{}))
					continue
				}
			} else {
				// TODO delete
				v.SetMapIndex(iter.Key(), reflect.Value{})
			}
		}
	case reflect.Slice, reflect.Array:
		of := reflect.ValueOf(inputJSON)
		for i := 0; i < of.Len(); i++ {
			JSONFilter2(of.Index(i).Interface(), requiredFields)
			continue
		}
	case reflect.Ptr:
		JSONFilter2(reflect.ValueOf(inputJSON).Elem().Interface(), requiredFields)
	default:

	}
}

const (
	str              = "string"
	obj              = "object"
	arr              = "array"
	num              = "number"
	dateTime         = "datetime"
	boolean          = "boolean"
	decimal          = "decimal"
	editPermission   = 2
	maxPermissionNum = 2 //字段权限位置由后端控制位数，当前是后两位，ps: 0b0011
	minPermissionNum = 1 //字段权限位置由后端控制位数，当前是后两位，ps: 0b0001
)

// DealSchemaToFilterType 将schema处理成过滤器需要的格式
func DealSchemaToFilterType(schema models.Schema) map[string]interface{} {
	out := make(map[string]interface{})
	if schema.Types == obj && schema.XInternal.Permission != 0 {
		for k := range schema.Properties {
			if schema.Properties[k].XInternal.Permission != 0 {
				if schema.Properties[k].Types == obj || schema.Properties[k].Types == arr {
					out[k] = schema.Properties[k].XInternal.Permission
					res := DealSchemaToFilterType(schema.Properties[k])
					if res != nil {
						for k1, v1 := range res {
							out[k1] = v1
						}
					}
					continue
				} else {
					//对字段权限进行处理，只留后两位数据
					f := schema.Properties[k].XInternal.Permission
					permission := getPermissionFromSchemaPermission(int(f), maxPermissionNum)
					if permission != 0 {
						out[k] = permission
					}
					continue
				}
			}
		}
		return out
	}
	if schema.Types == arr && schema.XInternal.Permission != 0 {
		filterType := DealSchemaToFilterType(*schema.Item)
		if len(filterType) > 0 {
			return filterType
		}
		return nil
	}
	return nil
}

func getPermissionFromSchemaPermission(permission int, n int) int {
	sprintf := fmt.Sprintf("%b", permission)
	var per = 0
	var mask = 1
	if len(sprintf) <= n {
		return permission
	}
	for i := n; i < len(sprintf); i++ {
		per = per + (mask << uint(i))
	}
	return (permission | per) ^ per
}

// SchemaFilterToNewSchema2 将全量schema过滤成不同权限组需要的
func SchemaFilterToNewSchema2(oldSchema interface{}, filter map[string]interface{}) {
	switch reflect.TypeOf(oldSchema).Kind() {
	case reflect.Map:
		v := reflect.ValueOf(oldSchema)
		if value := v.MapIndex(reflect.ValueOf("type")); value.IsValid() {
			if value.Elem().String() == arr {
				if itemValue := v.MapIndex(reflect.ValueOf("item")); itemValue.IsValid() {
					SchemaFilterToNewSchema2(itemValue.Interface(), filter)
				}
			}
			if value.Elem().String() == obj {
				if propertiesValue := v.MapIndex(reflect.ValueOf("properties")); propertiesValue.IsValid() {
					schemaFilter2(propertiesValue.Interface(), filter)
				}
			}
		}
	case reflect.Slice, reflect.Array:
		of := reflect.ValueOf(oldSchema)
		for i := 0; i < of.Len(); i++ {
			if value := of.Index(i); value.IsValid() {
				SchemaFilterToNewSchema2(of.Index(i).Interface(), filter)
			}
			continue
		}
	case reflect.Ptr:
		if reflect.ValueOf(oldSchema).IsValid() {
			SchemaFilterToNewSchema2(reflect.ValueOf(oldSchema).Elem().Interface(), filter)
		}
	default:

	}
}
func schemaFilter2(oldSchema interface{}, filter map[string]interface{}) {
	switch reflect.TypeOf(oldSchema).Kind() {
	case reflect.Map:
		v := reflect.ValueOf(oldSchema)
		iter := v.MapRange()
		for iter.Next() {
			if _, ok := filter[iter.Key().String()]; ok {
				switch reflect.TypeOf(filter[iter.Key().String()]).Kind() {
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Float32, reflect.Float64:
					if iter.Value().IsValid() {
						if !schemaUpdatePermission2(iter.Value().Interface(), filter[iter.Key().String()]) {
							// TODO delete
							v.SetMapIndex(iter.Key(), reflect.Value{})
						}
					}
					continue
				default:
					if iter.Value().IsValid() {
						SchemaFilterToNewSchema2(iter.Value().Interface(), filter)
					}
					continue
				}
			} else {
				if !isLayoutComponent(iter.Value().Interface(), filter) {
					// TODO delete
					v.SetMapIndex(iter.Key(), reflect.Value{})
				}

			}
		}

	default:

	}
}
func isLayoutComponent(oldSchema interface{}, filter map[string]interface{}) bool {
	switch reflect.TypeOf(oldSchema).Kind() {
	case reflect.Map:
		v := reflect.ValueOf(oldSchema)
		if value := v.MapIndex(reflect.ValueOf("x-internal")); value.IsValid() {
			if isLayoutComponent(value.Interface(), filter) {
				if propertiesValue := v.MapIndex(reflect.ValueOf("properties")); propertiesValue.IsValid() {
					schemaFilter2(propertiesValue.Interface(), filter)
				}
				return true
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

func schemaUpdatePermission2(oldSchema interface{}, permission interface{}) bool {
	switch reflect.TypeOf(oldSchema).Kind() {
	case reflect.Map:
		v := reflect.ValueOf(oldSchema)

		if value := v.MapIndex(reflect.ValueOf("x-internal")); value.IsValid() {
			return schemaUpdatePermission2(value.Interface(), permission)
		}
		if value := v.MapIndex(reflect.ValueOf("permission")); value.IsValid() {
			var oldP = 0
			var setP = 0
			switch reflect.TypeOf(value.Interface()).Kind() {
			case reflect.Int:
				oldP = int(value.Interface().(int))
			case reflect.Int32:
				oldP = int(value.Interface().(int32))
			case reflect.Int64:
				oldP = int(value.Interface().(int64))
			case reflect.Float64:
				oldP = int(value.Interface().(float64))
			case reflect.Float32:
				oldP = int(value.Interface().(float32))
			default:
				return false
			}
			switch reflect.TypeOf(permission).Kind() {
			case reflect.Int:
				setP = permission.(int)
			case reflect.Int32:
				setP = int(permission.(int32))
			case reflect.Int64:
				setP = int(permission.(int64))

			case reflect.Float64:
				setP = int(permission.(float64))
			case reflect.Float32:
				setP = int(permission.(float32))
			}
			basicsPermission := setBasicsPermission(oldP, maxPermissionNum, setP)
			v.SetMapIndex(reflect.ValueOf("permission"), reflect.ValueOf(basicsPermission))
			return true
		}

	default:
	}
	return true
}

func getMaxPermission(n int) int {
	var mask = 1
	if n == 1 {
		return mask
	}
	for i := 0; i < n-1; i++ {
		mask = mask<<1 + 1
	}
	return mask
}

func setBasicsPermission(oldPermission int, n int, setPermission int) int {
	maxPermission := getMaxPermission(n)
	if setPermission > maxPermission {
		return oldPermission
	}
	basic := oldPermission | maxPermission
	var mask = 1
	for i := 0; i < n; i++ {
		basic = basic ^ (mask << uint(i))
	}
	return basic | setPermission
}

const (
	_id           = "_id"
	_createdAt    = "created_at"
	_creatorID    = "creator_id"
	_creatorName  = "creator_name"
	_updatedAt    = "updated_at"
	_modifierID   = "modifier_id"
	_modifierName = "modifier_name"
)

// FilterCheckData 提交数据时检查数据权限
func FilterCheckData(data interface{}, filter map[string]interface{}) bool {
	switch reflect.TypeOf(data).Kind() {
	case reflect.Map:
		v := reflect.ValueOf(data)
		iter := v.MapRange()
		for iter.Next() {
			switch iter.Key().String() {
			case _id, _createdAt, _creatorID, _creatorName, _updatedAt, _modifierID, _modifierName:
				//v.SetMapIndex(iter.Key(), reflect.Value{})
				continue
			default:
				if _, ok := filter[iter.Key().String()]; !ok {
					return false
				}
				switch reflect.TypeOf(filter[iter.Key().String()]).Kind() {
				case reflect.Int8:
					if (filter[iter.Key().String()].(int8) & editPermission) == 0 {
						return false
					}
				case reflect.Int:
					if (filter[iter.Key().String()].(int) & editPermission) == 0 {
						return false
					}
				case reflect.Int16:
					if (filter[iter.Key().String()].(int16) & editPermission) == 0 {
						return false
					}
				case reflect.Int32:
					if (filter[iter.Key().String()].(int32) & editPermission) == 0 {
						return false
					}
				case reflect.Int64:
					if (filter[iter.Key().String()].(int64) & editPermission) == 0 {
						return false
					}
				case reflect.Float32:
					if (int64(filter[iter.Key().String()].(float32)) & editPermission) == 0 {
						return false
					}
				case reflect.Float64:
					if (int64(filter[iter.Key().String()].(float64)) & editPermission) == 0 {
						return false
					}
				default:
					flag := FilterCheckData(iter.Value().Interface(), filter)

					if !flag {
						return false
					}
					continue
				}
			}
		}
		return true
	case reflect.Array, reflect.Slice:
		of := reflect.ValueOf(data)
		for i := 0; i < of.Len(); i++ {
			flag := FilterCheckData(of.Index(i).Interface(), filter)

			if !flag {
				return false
			}
			continue
		}
		return true
	default:
		return false
	}
}

// DefaultSchema 处理初始化的权限
func DefaultSchema(oldSchema interface{}) {
	switch reflect.TypeOf(oldSchema).Kind() {
	case reflect.Map:
		v := reflect.ValueOf(oldSchema)
		if value := v.MapIndex(reflect.ValueOf("type")); value.IsValid() {
			if value.Elem().String() == arr {
				if itemValue := v.MapIndex(reflect.ValueOf("item")); itemValue.IsValid() {
					DefaultSchema(itemValue.Interface())
				}
			}
			if value.Elem().String() == obj {
				if value1 := v.MapIndex(reflect.ValueOf("x-internal")); value1.IsValid() {
					defaultSchemaUpdatePermission2(value1.Interface(), maxPermissionNum)
				}
				if propertiesValue := v.MapIndex(reflect.ValueOf("properties")); propertiesValue.IsValid() {
					defaultSchema(propertiesValue.Interface())
				}
			}
		}
	case reflect.Slice, reflect.Array:
		of := reflect.ValueOf(oldSchema)
		for i := 0; i < of.Len(); i++ {
			if value := of.Index(i); value.IsValid() {
				DefaultSchema(of.Index(i).Interface())
			}
			continue
		}
	case reflect.Ptr:
		if reflect.ValueOf(oldSchema).IsValid() {
			DefaultSchema(reflect.ValueOf(oldSchema).Elem().Interface())
		}
	default:

	}
}

func defaultSchema(oldSchema interface{}) {
	switch reflect.TypeOf(oldSchema).Kind() {
	case reflect.Map:
		v := reflect.ValueOf(oldSchema)
		iter := v.MapRange()
		for iter.Next() {
			switch iter.Key().String() {
			case _id, _createdAt, _creatorID, _creatorName, _updatedAt, _modifierID, _modifierName:
				defaultSchemaUpdatePermission2(iter.Value().Interface(), minPermissionNum)
			default:
				if defaultIsLayoutComponent(iter.Value().Interface()) {
					continue
				} else {
					if !defaultSchemaUpdatePermission2(iter.Value().Interface(), maxPermissionNum) {
						// TODO delete
						v.SetMapIndex(iter.Key(), reflect.Value{})
					}
					continue
				}
			}

		}
	default:

	}
}

func defaultIsLayoutComponent(oldSchema interface{}) bool {
	switch reflect.TypeOf(oldSchema).Kind() {
	case reflect.Map:
		v := reflect.ValueOf(oldSchema)
		if value := v.MapIndex(reflect.ValueOf("type")); value.IsValid() {
			if value.Elem().String() == arr {
				if itemValue := v.MapIndex(reflect.ValueOf("item")); itemValue.IsValid() {
					DefaultSchema(itemValue.Interface())
				}
			}
			if value.Elem().String() == obj {
				if propertiesValue := v.MapIndex(reflect.ValueOf("properties")); propertiesValue.IsValid() {
					defaultSchema(propertiesValue.Interface())
				}
			}
		}
		if value := v.MapIndex(reflect.ValueOf("x-internal")); value.IsValid() {
			if defaultIsLayoutComponent(value.Interface()) {
				if propertiesValue := v.MapIndex(reflect.ValueOf("properties")); propertiesValue.IsValid() {
					defaultSchema(propertiesValue.Interface())
				}
				return true
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

func defaultSchemaUpdatePermission2(oldSchema interface{}, permissionNum int) bool {
	switch reflect.TypeOf(oldSchema).Kind() {
	case reflect.Map:
		v := reflect.ValueOf(oldSchema)
		if value := v.MapIndex(reflect.ValueOf("x-internal")); value.IsValid() {
			return defaultSchemaUpdatePermission2(value.Interface(), permissionNum)
		}
		if value := v.MapIndex(reflect.ValueOf("permission")); value.IsValid() {
			var oldP = 0
			var setP = 0
			switch reflect.TypeOf(value.Interface()).Kind() {
			case reflect.Int:
				oldP = value.Interface().(int)
			case reflect.Int32:
				oldP = int(value.Interface().(int32))
			case reflect.Int64:
				oldP = int(value.Interface().(int64))
			case reflect.Float32:
				oldP = int(value.Interface().(float32))
			case reflect.Float64:
				oldP = int(value.Interface().(float64))
			default:
				return false
			}
			setP = getMaxPermission(permissionNum)
			basicsPermission := setBasicsPermission(oldP, maxPermissionNum, setP)
			v.SetMapIndex(reflect.ValueOf("permission"), reflect.ValueOf(basicsPermission))
			return true
		}

	default:
	}
	return true
}
