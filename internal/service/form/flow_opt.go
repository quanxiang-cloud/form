package form

import (
	"context"
	"github.com/quanxiang-cloud/form/internal/service/types"
)

type optionReq struct {
	tableID string
	userId  string
	query   types.Query
}

type Options func(ctx context.Context, req *optionReq) error

// CreateOption CreateOption
func CreateOption() Options {
	return func(ctx context.Context, option *optionReq) error {
		//if f1, ok := f.(*form); ok {
		//	data := new(inform.FormData)
		//	data.TableID = option.tableID
		//	data.Entity = option.entity
		//	comet.DefaultFormFiled(ctx, data, Post)
		//	logger.Logger.Infof(" %s send kafk data:   %+v : ", logger.STDRequestID(ctx).String, data)
		//	f1.hook.Send <- data
		//}
		//return nil
		return nil
	}
}
