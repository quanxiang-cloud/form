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
	intercept = os.Getenv("FORM_INTERCEPT") //拦截
	if intercept == "" {
		intercept = "true"
	}
}

const (
	object = "object"
	array  = "array"
)

func Filter(entity interface{}, fieldPermit models.FiledPermit) {
	if intercept == "true" {
		return
	}
	if entity == nil || fieldPermit == nil {
		return
	}
	value := reflect.ValueOf(entity)
	switch reflect.TypeOf(entity).Kind() {
	case reflect.Ptr:
		if value.Elem().CanInterface() {
			Filter(value.Elem().Interface(), fieldPermit)
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
				Filter(iter.Value().Interface(), permit.Properties)
			}
		}
	case reflect.Slice, reflect.Array:
		for i := 0; i < value.Len(); i++ {
			if value.Index(i).CanInterface() {
				Filter(value.Index(i).Interface(), fieldPermit)
			}
		}
	}
}
