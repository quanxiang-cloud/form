package service

import (
	"context"
	"fmt"
	"github.com/quanxiang-cloud/cabin/logger"
	"github.com/quanxiang-cloud/cabin/tailormade/header"
	config2 "github.com/quanxiang-cloud/form/pkg/misc/config"
)

type CreateTableReq struct {
	AppID    string
	TableID  string                 `json:"tableID"`
	Schema   map[string]interface{} `json:"schema"`
	UserID   string                 `json:"user_id"`
	UserName string                 `json:"user_name"`
	//Source   models.SourceType      `json:"source"` // source 1 是表单驱动，2是模型驱动
}

type CreateTableResp struct {
}

//CreateSchemaOption CreateSchemaOption
type CreateSchemaOption func(ctx context.Context) error

func ConvertSchema(t Table, req *CreateTableReq) CreateSchemaOption {
	return func(ctx context.Context) error {
		t1, ok := t.(*table)
		if !ok {
			return nil
		}
		fmt.Println(t1)
		return nil
	}
}

func RegisterSwagger(t Table, req *CreateTableReq) CreateSchemaOption {
	return func(ctx context.Context) error {
		t1, ok := t.(*table)
		if !ok {
			return nil
		}
		fmt.Println(t1)
		return nil
	}
}

func ComponentHandle(t Table, req *CreateTableReq) CreateSchemaOption {
	return func(ctx context.Context) error {
		t1, ok := t.(*table)
		if !ok {
			return nil
		}
		fmt.Println(t1)
		return nil
	}
}

type Table interface {
	CreateSchema(ctx context.Context, req *CreateTableReq, options ...CreateSchemaOption) (*CreateTableResp, error)
}

type table struct {
}

func NewTable(conf *config2.Config) (Table, error) {
	return &table{}, nil
}

func (t *table) CreateSchema(ctx context.Context, req *CreateTableReq, options ...CreateSchemaOption) (resp *CreateTableResp, err error) {
	// 1、保存schema
	defer func() {
		if err == nil {
			t.createSchemaPost(ctx, options...)
		}
	}()

	//

	return &CreateTableResp{}, nil
}

func (t *table) createSchemaPost(ctx context.Context, options ...CreateSchemaOption) {
	for _, options := range options {
		err := options(ctx)
		if err != nil {
			// 打印错误日志
			logger.Logger.Errorw("create ms err ", err.Error(), header.GetRequestIDKV(ctx).Fuzzy())
		}
	}
}
