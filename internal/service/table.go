package service

import (
	"context"
	id2 "github.com/quanxiang-cloud/cabin/id"
	"github.com/quanxiang-cloud/cabin/logger"
	"github.com/quanxiang-cloud/cabin/tailormade/header"
	"github.com/quanxiang-cloud/form/internal/models"
	"github.com/quanxiang-cloud/form/internal/models/mysql"
	"github.com/quanxiang-cloud/form/internal/service/swagger"
	config2 "github.com/quanxiang-cloud/form/pkg/misc/config"
	"gorm.io/gorm"
)

type Table interface {
	CreateSchema(ctx context.Context, req *CreateTableReq, options ...CreateSchemaOption) (*CreateTableResp, error)
}

type table struct {
	db              *gorm.DB
	tableRepo       models.TableRepo
	tableSchemaRepo models.TableSchemeRepo
}

func NewTable(conf *config2.Config) (Table, error) {
	db, err := createMysqlConn(conf)
	if err != nil {
		return nil, err
	}
	return &table{
		db:              db,
		tableRepo:       mysql.NewTableRepo(),
		tableSchemaRepo: mysql.NewTableSchema(),
	}, nil
}

type CreateTableReq struct {
	AppID    string            `json:"app_id"`
	TableID  string            `json:"tableID"`
	Schema   models.WebSchema  `json:"schema"`
	UserID   string            `json:"user_id"`
	UserName string            `json:"user_name"`
	Source   models.SourceType `json:"source"` // source 1 是表单驱动，2是模型驱动
}

type CreateTableResp struct {
}

type Base struct {
	UserName string
	UserID   string
	AppID    string
	TableID  string
	t        *table
}
type TableSchema struct {
	Title       string
	FieldLen    int64
	Description string
	Source      models.SourceType
	Schema      models.TableSchemas
}
type Poly struct {
	PolySchema map[string]interface{}
}

type DDL struct {
	OldSchema models.TableSchemas

	NewSchema models.TableSchemas
}
type CreateOptReq struct {
	Base
	TableSchema
	Poly
}

//CreateSchemaOption CreateSchemaOption
type CreateSchemaOption func(ctx context.Context, req *CreateOptReq) error

func ConvertSchema() CreateSchemaOption {
	return func(ctx context.Context, req *CreateOptReq) error {
		one, err := req.t.tableSchemaRepo.Get(req.t.db, req.TableID, req.AppID)
		if err != nil {
			return err
		}
		tables := &models.TableSchema{
			Title:       req.Title,
			Schema:      req.Schema,
			FieldLen:    req.FieldLen,
			Description: req.Description,
		}
		if one == nil {
			tables.ID = id2.HexUUID(true)
			tables.Source = req.Source
			tables.AppID = req.AppID
			tables.TableID = req.TableID
			//table.CreatedAt = time.Now()
			tables.CreatorName = req.UserName
			tables.CreatorID = req.UserID
			err = req.t.tableSchemaRepo.BatchCreate(req.t.db, tables)
			if err != nil {
				return err
			}
			return nil
		}
		//tables.UpdatedAt = time2.NowUnix()
		tables.EditorID = req.UserID
		tables.EditorName = req.UserName
		err = req.t.tableSchemaRepo.Update(req.t.db, req.AppID, req.TableID, tables)
		if err != nil {
			return err
		}
		return nil
	}
}

func RegisterSwagger() CreateSchemaOption {
	return func(ctx context.Context, req *CreateOptReq) error {

		return nil
	}
}

// ComponentHandle ComponentHandle 创建表结构， 创建索引
func ComponentHandle() CreateSchemaOption {
	return func(ctx context.Context, req *CreateOptReq) error {
		return nil
	}
}

func (t *table) CreateSchema(ctx context.Context, req *CreateTableReq, options ...CreateSchemaOption) (resp *CreateTableResp, err error) {
	// 1、保存schema
	defer func() {
		if err == nil {
			t.createSchemaPost(ctx, req, options...)
		}
	}()
	one, err := t.tableRepo.Get(t.db, req.AppID, req.TableID)
	if err != nil {
		return nil, err
	}
	tables := &models.Table{
		Schema: req.Schema,
	}
	if one == nil {
		tables.ID = id2.HexUUID(true)
		tables.TableID = req.TableID
		tables.AppID = req.AppID
		err = t.tableRepo.BatchCreate(t.db, tables)
		if err != nil {
			return nil, err
		}
		return &CreateTableResp{}, nil
	}
	err = t.tableRepo.Update(t.db, req.AppID, req.TableID, tables)
	if err != nil {
		return nil, err
	}
	// 准备数据

	return &CreateTableResp{}, nil
}

func (t *table) createSchemaPost(ctx context.Context, req *CreateTableReq, options ...CreateSchemaOption) {
	properties, err := getMapToMao(req.Schema, "properties")
	if err != nil {
		return
	}

	optReq := &CreateOptReq{}
	base := Base{
		UserName: req.UserName,
		UserID:   req.UserID,
		AppID:    req.AppID,
		TableID:  req.TableID,
		t:        t,
	}
	optReq.Base = base

	convert, total, err := swagger.Convert1(properties)
	description := getMapToString(req.Schema, "description")
	title := getMapToString(req.Schema, "title")

	baseSchema := TableSchema{
		Title:       title,
		Description: description,
		Source:      req.Source,
		Schema:      convert,
		FieldLen:    total,
	}
	optReq.TableSchema = baseSchema

	optReq.Poly = Poly{
		PolySchema: convert,
	}
	for _, options := range options {
		err := options(ctx, optReq)
		if err != nil {
			// 打印错误日志
			logger.Logger.Errorw("create ms err ", err.Error(), header.GetRequestIDKV(ctx).Fuzzy())
		}
	}
}

func getMapToMao(schema map[string]interface{}, key string) (map[string]interface{}, error) {
	value, ok := schema[key]
	if !ok {
		return nil, nil
	}
	if v, ok := value.(map[string]interface{}); ok {
		return v, nil
	}
	return nil, nil
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
