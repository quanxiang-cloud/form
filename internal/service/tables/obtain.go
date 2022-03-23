package tables

import (
	"git.internal.yunify.com/qxp/misc/error2"
	"github.com/quanxiang-cloud/form/pkg/misc/code"
)

func getMapToMap(schema map[string]interface{}, key string) (map[string]interface{}, error) {
	value, ok := schema[key]
	if !ok {
		return nil, error2.NewError(code.ErrParameter)
	}
	if v, ok := value.(map[string]interface{}); ok {
		return v, nil
	}
	return nil, error2.NewError(code.ErrParameter)
}

// getMapToString getMapToString
func getMapToString(schema map[string]interface{}, key string) string {
	value, ok := schema[key]
	if !ok {
		return ""
	}
	if v, ok := value.(string); ok {
		return v
	}
	return ""
}

// GetAsMap GetAsMap
func getAsMap(v interface{}) (map[string]interface{}, error) {
	if m, ok := v.(map[string]interface{}); ok {
		return m, nil
	}
	return nil, error2.NewError(code.ErrParameter)
}

func getMapToBool(schema map[string]interface{}, key string) (bool, error) {
	value, ok := schema[key]
	if !ok {
		return false, error2.NewError(code.ErrParameter)
	}
	if v, ok := value.(bool); ok {
		return v, nil
	}
	return false, error2.NewError(code.ErrParameter)
}
