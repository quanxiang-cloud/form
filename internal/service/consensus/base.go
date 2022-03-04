package consensus

import (
	"reflect"
)

const (
	// TermsKey TermsKey
	TermsKey = "terms"
	TermKey  = "term"
	IDKey    = "_id"
)

// KeyValue KeyValue
type KeyValue map[string]interface{}

// GetSimple GetSimple
func GetSimple(tag, fieldName string, value interface{}) map[string]interface{} {
	return map[string]interface{}{
		tag: KeyValue{
			fieldName: value,
		},
	}
}

// GetBool GetBool
func GetBool(tag string, value ...interface{}) map[string]interface{} {
	return map[string]interface{}{
		"bool": KeyValue{
			tag: value,
		},
	}
}

func GetIDByQuery(query map[string]interface{}) []string {
	term := Term(query)
	if term == nil {
		return nil
	}

	val := reflect.ValueOf(term)
	if !val.CanInterface() {
		return nil
	}
	ids := make([]string, 0)
	if v, ok := val.Interface().([]interface{}); ok {
		for _, value := range v {
			id1, ok := value.(string)
			if ok {
				ids = append(ids, id1)
			}
		}

	} else if reflect.TypeOf(term).Kind() == reflect.String {
		id1, ok := term.(string)
		if ok {
			ids = append(ids, id1)
		}

	}
	return ids
}

// Term Term
func Term(data interface{}) interface{} {
	switch reflect.TypeOf(data).Kind() {
	case reflect.Map:
		// 看有没有terms
		v := reflect.ValueOf(data)
		if value := v.MapIndex(reflect.ValueOf(TermsKey)); value.IsValid() {
			return Term(value.Elem().Interface())
		}
		if value := v.MapIndex(reflect.ValueOf(TermKey)); value.IsValid() {
			return Term(value.Elem().Interface())
		}
		if value := v.MapIndex(reflect.ValueOf(IDKey)); value.IsValid() {
			return value.Interface()
		}

	}
	return nil
}
