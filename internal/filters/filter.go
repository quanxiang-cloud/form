package filters

import (
	"github.com/quanxiang-cloud/form/internal/models"
	"reflect"
)

const (
	object = "object"
	array  = "array"
)

func pre(entity interface{}, fieldPermit models.FiledPermit) bool {
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
				if !pre(iter.Value(), permit.Properties) {
					return false
				}
			}
		}
		return true
	case reflect.Slice, reflect.Array:
		for i := 0; i < value.Len(); i++ {
			pre(value.Index(i), fieldPermit)
		}
	}
	return false
}

func post(response interface{}, fieldPermit models.FiledPermit) {
	if fieldPermit == nil {
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
				post(iter.Value(), permit.Properties)
			}
		}
	case reflect.Slice, reflect.Array:
		for i := 0; i < value.Len(); i++ {
			post(value.Index(i), fieldPermit)
		}
	}
}
