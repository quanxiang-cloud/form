package treasure

import (
	"os"
	"reflect"

	"github.com/quanxiang-cloud/form/internal/models"
)

var (
	intercept string
)

func init() {
	intercept = os.Getenv("INTERCEPT") //拦截
	if intercept == "" {
		intercept = "false"
	}
}

const (
	object = "object"
	array  = "array"
)

func Pre(entity interface{}, fieldPermit models.FiledPermit) bool {
	if intercept == "false" {
		return false
	}
	if entity == nil {
		return false
	}
	value := reflect.ValueOf(entity)
	switch reflect.TypeOf(entity).Kind() {
	case reflect.Ptr:
		if value.Elem().CanInterface() {
			return Pre(value.Elem().Interface(), fieldPermit)
		}
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
			if value.Index(i).CanInterface() {
				Pre(value.Index(i).Interface(), fieldPermit)
			}
		}
	}
	return false
}

func Post(response interface{}, fieldPermit models.FiledPermit) {
	if intercept == "false" {
		return
	}
	if response == nil || fieldPermit == nil {
		return
	}
	value := reflect.ValueOf(response)
	switch reflect.TypeOf(response).Kind() {
	case reflect.Ptr:
		if value.Elem().CanInterface() {
			Post(value.Elem().Interface(), fieldPermit)
		}
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
			if value.Index(i).CanInterface() {
				Post(value.Index(i).Interface(), fieldPermit)
			}
		}
	}
}
