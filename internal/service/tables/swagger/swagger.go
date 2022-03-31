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
