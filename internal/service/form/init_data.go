package form

import (
	"context"
	"reflect"

	error2 "github.com/quanxiang-cloud/cabin/error"
	"github.com/quanxiang-cloud/cabin/logger"
	"github.com/quanxiang-cloud/form/pkg/misc/code"
)

const (
	_type     = "type"
	createKey = "new"

	updateKey = "updated"
	// DeleteKey the key of "deleted" in the ref.
	deleteKey = "deleted"
	// AppIDKey the component's AppID in the ref.
	appIDKey = "appID"
	// TableIDKey the component's TableID in the ref.
	tableIDKey = "tableID"
	// EntityKey the component's entity in the ref.

	entityKey = "entity"

	ref = "ref"

	queryKey = "query"

	// PrimitiveID PrimitiveID.
	primitiveID = "primitiveID"
	// SubID subID.
	subIDs = "subID"
	// FieldName  FieldName.
	fieldName = "fieldName"
)

type OptionEntity struct {
	Entity map[string]interface{} `json:"entity"`
	Query  map[string]interface{} `json:"query"`
	Ref    map[string]interface{} `json:"ref"`
}

func initOptionEntity(ctx context.Context, optValue map[string]interface{}, optionEntity *OptionEntity) (err error) {
	if tableID, ok := optValue[entityKey]; ok {
		err = SetFieldValue(ctx, tableID, &optionEntity.Entity)
	}
	if appID, ok := optValue[ref]; ok {
		err = SetFieldValue(ctx, appID, &optionEntity.Ref)
	}
	if create, ok := optValue[queryKey]; ok {
		err = SetFieldValue(ctx, create, &optionEntity.Query)
	}

	return
}

// RefData  map[string] interface{}.
type RefData struct {
	AppID         string                 `json:"appID"`
	TableID       string                 `json:"tableID"`
	Type          string                 `json:"type"`
	New           []interface{}          `json:"new"`
	Deleted       []interface{}          `json:"deleted"`
	Updated       []interface{}          `json:"updated"`
	Query         map[string]interface{} `json:"query"`
	SourceFieldID string                 `json:"sourceFieldId"`
	FieldName     string                 `json:"fieldName"`
	AggType       string                 `json:"aggType"`
}

func initRefData(ctx context.Context, optValue map[string]interface{}, data *RefData) (err error) {
	if tableID, ok := optValue[tableIDKey]; ok {
		err = SetFieldValue(ctx, tableID, &data.TableID)
	}
	if appID, ok := optValue[appIDKey]; ok {
		err = SetFieldValue(ctx, appID, &data.AppID)
	}
	if create, ok := optValue[createKey]; ok {
		err = SetFieldValue(ctx, create, &data.New)
	}
	if update, ok := optValue[updateKey]; ok {
		err = SetFieldValue(ctx, update, &data.Updated)
	}
	if delete1, ok := optValue[deleteKey]; ok {
		err = SetFieldValue(ctx, delete1, &data.Deleted)
	}
	if types, ok := optValue[_type]; ok {
		err = SetFieldValue(ctx, types, &data.Type)
	}
	if query, ok := optValue[queryKey]; ok {
		err = SetFieldValue(ctx, query, &data.Query)
	}
	if sourceFieldID, ok := optValue["sourceFieldId"]; ok {
		err = SetFieldValue(ctx, sourceFieldID, &data.SourceFieldID)
	}
	if fieldName, ok := optValue["fieldName"]; ok {
		err = SetFieldValue(ctx, fieldName, &data.FieldName)
	}
	if aggType, ok := optValue["aggType"]; ok {
		err = SetFieldValue(ctx, aggType, &data.AggType)
	}
	return
}

// ExtraData ExtraData.
type ExtraData struct {
	AppID   string `json:"appID"`
	TableID string `json:"tableID"`
	ID      string `json:"id"`
}

func initExtraData(ctx context.Context, optValue map[string]interface{}, extraData *ExtraData) (err error) {
	if tableID, ok := optValue[tableIDKey]; ok {
		err = SetFieldValue(ctx, tableID, &extraData.TableID)
	}
	if appID, ok := optValue[appIDKey]; ok {
		err = SetFieldValue(ctx, appID, &extraData.AppID)
	}
	if id, ok := optValue["_id"]; ok {
		err = SetFieldValue(ctx, id, &extraData.ID)
	}
	return
}

// SetFieldValue SetFieldValue.
func SetFieldValue(ctx context.Context, data, ptr interface{}) error {
	if data == nil {
		return nil
	}
	dateKind := reflect.TypeOf(data).Kind()
	value := reflect.ValueOf(data)
	ptrKind := reflect.TypeOf(ptr).Elem().Kind()
	if dateKind == reflect.Ptr {
		return SetFieldValue(ctx, value.Elem(), ptr)
	}
	if dateKind != ptrKind {
		logger.Logger.Errorw("dateKind type is not ptrKind", "", "")
		return error2.New(code.ErrParameter)
	}
	reflect.ValueOf(ptr).Elem().Set(value)
	return nil
}
