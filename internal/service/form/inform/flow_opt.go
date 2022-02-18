package inform

import (
	"context"
	"git.internal.yunify.com/qxp/misc/logger"
	"github.com/quanxiang-cloud/form/internal/service/form/base"
	"github.com/quanxiang-cloud/form/internal/service/types"
	"reflect"
)

const (
	termKey = "term"
	idKey   = "_id"
)

type OptionReq struct {
	CommonReq
	FlowReq
}

type CommonReq struct {
	tableID string
	userId  string
}

type FlowReq struct {
	query  types.Query
	entity interface{}
	send   chan *FormData
}

type Options func(ctx context.Context, req *OptionReq) error

// CreateOption CreateOption
func CreateOption() Options {
	return func(ctx context.Context, option *OptionReq) error {
		data := new(FormData)
		data.TableID = option.tableID
		data.Entity = option.entity
		DefaultFormFiled(ctx, data, "post")
		logger.Logger.Infof(" %s send kafk data:   %+v : ", logger.STDRequestID(ctx).String, data)
		option.send <- data
		return nil
	}
}

func UpdateOption() Options {
	return func(ctx context.Context, req *OptionReq) error {
		term := term(req.query)
		if term == nil {
			return nil
		}
		val := reflect.ValueOf(term)
		if !val.CanInterface() {
			//return error2.NewError(code.ErrParameter)
		}
		if v, ok := val.Interface().([]interface{}); ok {
			for _, value := range v {
				id1 := value.(string)
				req.entity = base.DefaultField(req.entity,
					base.WithUpdateID(id1),
				)

				data := &FormData{
					TableID: req.tableID,
					Entity:  req.entity,
				}
				DefaultFormFiled(ctx, data, "put")
				logger.Logger.Infof(" %s send kafk data:   %+v : ", logger.STDRequestID(ctx).String, data)
				req.send <- data
				return nil
			}

		}
		if reflect.TypeOf(term).Kind() == reflect.String {
			id1 := term.(string)
			req.entity = base.DefaultField(req.entity,
				base.WithUpdateID(id1),
			)
			//option.handleNum = 1
			data := &FormData{
				TableID: req.tableID,
				Entity:  req.entity,
			}
			DefaultFormFiled(ctx, data, "put")
			logger.Logger.Infof(" %s send kafk data:   %+v : ", logger.STDRequestID(ctx).String, data)
			req.send <- data
			return nil
		}
		return nil
	}
}

//DeleteOption DeleteOption
func DeleteOption() Options {
	return func(ctx context.Context, req *OptionReq) error {
		term := term(req.query)
		if val := reflect.ValueOf(term); val.CanInterface() {
			if v, ok := val.Interface().([]interface{}); ok {
				//option.handleNum = len(v)
				data := &FormData{
					//	TableID: option.tableID,
					Entity: map[string]interface{}{
						"data":      v,
						"delete_id": req.userId},
				}
				DefaultFormFiled(ctx, data, "delete")
				logger.Logger.Infof(" %s send kafk data:   %+v : ", logger.STDRequestID(ctx).String, data)
				req.send <- data
			}
		}
		return nil
	}
}

// Term Term
func term(data interface{}) interface{} {
	switch reflect.TypeOf(data).Kind() {
	case reflect.Map:
		// 看有没有terms
		v := reflect.ValueOf(data)
		if value := v.MapIndex(reflect.ValueOf(termKey)); value.IsValid() {
			return term(value.Elem().Interface())
		}
		if value := v.MapIndex(reflect.ValueOf(termKey)); value.IsValid() {
			return term(value.Elem().Interface())
		}
		if value := v.MapIndex(reflect.ValueOf(idKey)); value.IsValid() {
			return value.Interface()
		}

	}
	return nil
}
