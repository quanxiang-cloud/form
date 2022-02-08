package comet

import (
	"reflect"

	id2 "github.com/quanxiang-cloud/cabin/id"
	time2 "github.com/quanxiang-cloud/cabin/time"
)

// Entity Entity
type Entity interface{}

// EntityOpt entity options
type EntityOpt func(e defaultFieldMap)

type defaultFieldMap map[string]interface{}

const (
	_id           = "_id"
	_createdAt    = "created_at"
	_creatorID    = "creator_id"
	_creatorName  = "creator_name"
	_updatedAt    = "updated_at"
	_modifierID   = "modifier_id"
	_modifierName = "modifier_name"
)

// WithID default field with id
func WithID() EntityOpt {
	return func(d defaultFieldMap) {
		d[_id] = id2.ShortID(8)
	}
}

// WithUpdateID WithUpdateID
func WithUpdateID(updateID string) EntityOpt {
	return func(d defaultFieldMap) {
		d[_id] = updateID
	}
}

// WithCreated default field with created_at、creator_id and creator_name
func WithCreated(userID, userName string) EntityOpt {
	return func(d defaultFieldMap) {
		d[_createdAt] = time2.Now()
		d[_creatorID] = userID
		d[_creatorName] = userName
	}
}

// WithUpdated default field with updated_at、modifier_id and modifier_name
func WithUpdated(userID, userName string) EntityOpt {
	return func(d defaultFieldMap) {
		d[_updatedAt] = time2.Now()
		d[_modifierID] = userID
		d[_modifierName] = userName
	}
}

func defaultFieldWithDep(e Entity, dep int, opts ...EntityOpt) Entity {
	if e == nil {
		return e
	}
	value := reflect.ValueOf(e)
	switch _t := reflect.TypeOf(e); _t.Kind() {
	case reflect.Ptr:
		return defaultFieldWithDep(value.Elem(), dep, opts...)
	case reflect.Array, reflect.Slice:
		for i := 0; i < value.Len(); i++ {
			if !value.Index(i).CanInterface() {
				continue
			}
			val := defaultFieldWithDep(value.Index(i).Interface(), dep, opts...)
			value.Index(i).Set(reflect.ValueOf(val))
		}
	case reflect.Map:
		if dep > 0 {
			dep--
			iter := value.MapRange()
			for iter.Next() {
				if !iter.Value().CanInterface() {
					continue
				}
				val := defaultFieldWithDep(iter.Value().Interface(), dep, opts...)
				value.SetMapIndex(reflect.ValueOf(iter.Key().String()), reflect.ValueOf(val))
			}
			return e
		}
		defaultFieldMap := make(map[string]interface{})
		for _, opt := range opts {
			opt(defaultFieldMap)
		}
		for key, val := range defaultFieldMap {
			value.SetMapIndex(reflect.ValueOf(key), reflect.ValueOf(val))
		}
	default:
		return e
	}

	return e
}

// DefaultField DefaultField
func DefaultField(e Entity, opts ...EntityOpt) Entity {
	return defaultFieldWithDep(e, 0, opts...)
}
