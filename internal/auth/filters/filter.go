package filters

import (
	"reflect"

	"github.com/quanxiang-cloud/form/internal/models"
)

const (
	object = "object"
	array  = "array"
)

func Pre(entity interface{}, fieldPermit models.FiledPermit) bool {
	value := reflect.ValueOf(entity)
	switch reflect.TypeOf(entity).Kind() {
	case reflect.Map:
		iter := value.MapRange()
		for iter.Next() {
			if !iter.Value().CanInterface() {
				continue
			}
			key := iter.Key().String()
			permit, ok := fieldPermit[key]
			if !ok {
				return false
			}
			// 如果
			if permit.Type == object {
				if !Pre(iter.Value(), permit.Properties) {
					return false
				}
			}
		}
		return true
	case reflect.Slice, reflect.Array:
		for i := 0; i < value.Len(); i++ {
			Pre(value.Index(i), fieldPermit)
		}
	}
	return false
}

func Post(response interface{}, fieldPermit models.FiledPermit) {

	if response == nil || fieldPermit == nil {
		return
	}
	value := reflect.ValueOf(response)
	switch reflect.TypeOf(response).Kind() {
	case reflect.Map:
		iter := value.MapRange()
		for iter.Next() {
			if !iter.Value().CanInterface() {
				continue
			}
			key := iter.Key().String()
			permit, ok := fieldPermit[key]
			if !ok {
				value.SetMapIndex(iter.Key(), reflect.Value{})
			}
			// 如果
			if permit.Type == object || permit.Type == array {
				Post(iter.Value(), permit.Properties)
			}
		}
	case reflect.Slice, reflect.Array:
		for i := 0; i < value.Len(); i++ {
			Post(value.Index(i), fieldPermit)
		}
	}
}
