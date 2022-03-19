package tables

import (
	"context"
	"github.com/quanxiang-cloud/form/internal/models"
	"github.com/quanxiang-cloud/form/internal/models/mysql"
	"github.com/quanxiang-cloud/form/internal/service"
	config2 "github.com/quanxiang-cloud/form/pkg/misc/config"
	"gorm.io/gorm"
)

type Table interface {
	GetTable(ctx context.Context, req *GetTableReq) (*GetTableResp, error)
	DeleteTable(ctx context.Context, req *DeleteTableReq) (*DeleteTableResp, error)
	FindTable(ctx context.Context, req *FindTableReq) (*FindTableResp, error)
	UpdateConfig(ctx context.Context, req *UpdateConfigReq) (*UpdateConfigResp, error)
}

type table struct {
	db              *gorm.DB
	tableRepo       models.TableRepo
	tableSchemaRepo models.TableSchemeRepo
}

func NewTable(conf *config2.Config) (Table, error) {
	db, err := service.CreateMysqlConn(conf)
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

type GetTableReq struct {
	AppID   string `json:"appID"`
	TableID string `json:"tableID"`
}

type GetTableResp struct {
	// id pk
	ID string `json:"id"`
	// table design json schema
	Schema models.WebSchema `json:"schema"`
	// table page config json schema
	Config models.Config `json:"config"`
}

func (t *table) GetTable(ctx context.Context, req *GetTableReq) (*GetTableResp, error) {
	tables, err := t.tableRepo.Get(t.db, req.AppID, req.TableID)
	if err != nil {
		return nil, err
	}
	return &GetTableResp{
		ID:     tables.ID,
		Schema: tables.Schema,
		Config: tables.Config,
	}, nil
}

type DeleteTableReq struct {
	AppID   string `json:"appID"`
	TableID string `json:"tableID"`
}

type DeleteTableResp struct {
}

// DeleteTable 不开启事务
func (t *table) DeleteTable(ctx context.Context, req *DeleteTableReq) (*DeleteTableResp, error) {
	err := t.tableRepo.Delete(t.db, &models.TableQuery{
		TableID: req.TableID,
		AppID:   req.AppID,
	})
	if err != nil {
		return nil, err
	}
	err = t.tableSchemaRepo.Delete(t.db, &models.TableSchemaQuery{
		TableID: req.TableID,
	})
	if err != nil {
		return nil, err
	}
	return &DeleteTableResp{}, nil

}

type FindTableReq struct {
	Title  string            `json:"title"`
	AppID  string            `json:"appID"`
	Page   int64             `json:"page"`
	Size   int64             `json:"size"`
	Source models.SourceType `json:"source"`
}

type FindTableResp struct {
	List  []*tableVo `json:"list"`
	Total int64      `json:"total"`
}

// tableVo tableVo
type tableVo struct {
	ID          string            `json:"id"`
	TableID     string            `json:"tableID"`
	Title       string            `json:"title"`
	Description string            `json:"description"`
	FieldLen    int64             `json:"fieldLen"`
	Source      models.SourceType `json:"source"`
	CreatedAt   int64             `json:"createdAt"`
	UpdatedAt   int64             `json:"updatedAt"`
	Editor      string            `json:"editor"`
	CreatorName string            `json:"creatorName"`
}

func (t *table) FindTable(ctx context.Context, req *FindTableReq) (*FindTableResp, error) {
	tables, total, err := t.tableSchemaRepo.Find(t.db, &models.TableSchemaQuery{}, req.Size, req.Page)
	if err != nil {
		return nil, err
	}
	resp := &FindTableResp{
		List: make([]*tableVo, len(tables)),
	}
	for index, v := range tables {
		vo := &tableVo{
			ID:          v.ID,
			TableID:     v.TableID,
			Title:       v.Title,
			Description: v.Description,
			Source:      v.Source,
			FieldLen:    v.FieldLen,
			CreatedAt:   v.CreatedAt,
			UpdatedAt:   v.UpdatedAt,
			Editor:      v.EditorName,
			CreatorName: v.CreatorName,
		}
		resp.List[index] = vo
	}
	resp.Total = total
	return resp, nil
}

type UpdateConfigReq struct {
	TableID string        `json:"tableID"`
	AppID   string        `json:"appID"`
	Config  models.Config `json:"config"`
}

type UpdateConfigResp struct {
}

func (t *table) UpdateConfig(ctx context.Context, req *UpdateConfigReq) (*UpdateConfigResp, error) {
	tables := &models.Table{
		TableID: req.TableID,
		Config:  req.Config,
		AppID:   req.AppID,
	}

	err := t.tableRepo.Update(t.db, req.AppID, req.TableID, tables)

	if err != nil {
		return nil, err
	}
	return &UpdateConfigResp{}, nil
}
