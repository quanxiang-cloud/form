package service

import (
	"context"

	"github.com/quanxiang-cloud/form/internal/models"
	"github.com/quanxiang-cloud/form/internal/models/mysql"
	config2 "github.com/quanxiang-cloud/form/pkg/misc/config"

	"github.com/quanxiang-cloud/cabin/logger"
	"github.com/quanxiang-cloud/cabin/tailormade/header"
	"gorm.io/gorm"
)

// Backup import and export data interface.
type Backup interface {
	ExportTable(context.Context, *ExportTableReq) (*ExportTableResp, error)
	ExportPermit(context.Context, *ExportPermitReq) (*ExportPermitResp, error)
	ExportTableRelation(context.Context, *ExportTableRelationReq) (*ExportTableRelationResp, error)
	ExportTableScheme(context.Context, *ExportTableSchemeReq) (*ExportTableSchemeResp, error)
	ExportRole(context.Context, *ExportRoleReq) (*ExportRoleResp, error)

	ImportTable(context.Context, *ImportTableReq) (*ImportTableResp, error)
	ImportPermit(context.Context, *ImportPermitReq) (*ImportPermitResp, error)
	ImportTableRelation(context.Context, *ImportTableRelationReq) (*ImportTableRelationResp, error)
	ImportTableScheme(context.Context, *ImportTableSchemeReq) (*ImportTableSchemeResp, error)
	ImportRole(context.Context, *ImportRoleReq) (*ImportRoleResp, error)
}

type backup struct {
	db                *gorm.DB
	roleRepo          models.RoleRepo
	permitRepo        models.PermitRepo
	tableRepo         models.TableRepo
	tableRelationRepo models.TableRelationRepo
	tableSchemeRepo   models.TableSchemeRepo
}

// NewBackup create a new backup service.
func NewBackup(conf *config2.Config) (Backup, error) {
	db, err := CreateMysqlConn(conf)
	if err != nil {
		return nil, err
	}

	return &backup{
		db:                db,
		roleRepo:          mysql.NewRoleRepo(),
		permitRepo:        mysql.NewPermitRepo(),
		tableRepo:         mysql.NewTableRepo(),
		tableRelationRepo: mysql.NewTableRelationRepo(),
		tableSchemeRepo:   mysql.NewTableSchema(),
	}, nil
}

type ExportTableReq struct {
	AppID string `uri:"appID"`
	Page  int    `json:"page"`
	Size  int    `json:"size"`
}

type ExportTableResp struct {
	Data  []*models.Table `json:"data"`
	Count int64           `json:"count"`
}

// ExportTable export table data.
func (b *backup) ExportTable(ctx context.Context, req *ExportTableReq) (*ExportTableResp, error) {
	tables, count, err := b.tableRepo.List(b.db, &models.TableQuery{
		AppID: req.AppID,
	},
		req.Page, req.Size)
	if err != nil {
		logger.Logger.WithName("ExportTable").Errorw(err.Error(), header.GetRequestIDKV(ctx).Fuzzy()...)
		return nil, err
	}

	return &ExportTableResp{
		Data:  tables,
		Count: count,
	}, nil
}

type ExportPermitReq struct {
	AppID string `uri:"appID"`
	Page  int    `json:"page"`
	Size  int    `json:"size"`
}

type ExportPermitResp struct {
	Data  []*models.Permit `json:"data"`
	Count int64            `json:"count"`
}

// ExportPermit export permit data.
func (b *backup) ExportPermit(ctx context.Context, req *ExportPermitReq) (*ExportPermitResp, error) {
	roles, count, err := b.roleRepo.List(b.db, &models.RoleQuery{
		AppID: req.AppID,
	},
		req.Page, req.Size)
	if err != nil {
		logger.Logger.WithName("ExportPermit").Errorw(err.Error(), header.GetRequestIDKV(ctx).Fuzzy()...)
		return nil, err
	}

	if len(roles) == 0 {
		return &ExportPermitResp{
			Data:  nil,
			Count: 0,
		}, nil
	}

	ids := make([]string, 0, count)
	for _, role := range roles {
		ids = append(ids, role.ID)
	}

	permits, count, err := b.permitRepo.List(b.db, &models.PermitQuery{
		RoleIDs: ids,
	},
		req.Page, req.Size)
	if err != nil {
		logger.Logger.WithName("ExportPermit").Errorw(err.Error(), header.GetRequestIDKV(ctx).Fuzzy()...)
		return nil, err
	}

	return &ExportPermitResp{
		Data:  permits,
		Count: count,
	}, nil
}

type ExportTableRelationReq struct {
	AppID string `uri:"appID"`
	Page  int    `json:"page"`
	Size  int    `json:"size"`
}

type ExportTableRelationResp struct {
	Data  []*models.TableRelation `json:"data"`
	Count int64                   `json:"count"`
}

// ExportTableRelation export table relation data.
func (b *backup) ExportTableRelation(ctx context.Context, req *ExportTableRelationReq) (*ExportTableRelationResp, error) {
	tableRelations, count, err := b.tableRelationRepo.List(b.db, &models.TableRelationQuery{
		AppID: req.AppID,
	},
		req.Page, req.Size)
	if err != nil {
		logger.Logger.WithName("ExportTableRelation").Errorw(err.Error(), header.GetRequestIDKV(ctx).Fuzzy()...)
		return nil, err
	}

	return &ExportTableRelationResp{
		Data:  tableRelations,
		Count: count,
	}, nil
}

type ExportTableSchemeReq struct {
	AppID string `uri:"appID"`
	Page  int    `json:"page"`
	Size  int    `json:"size"`
}

type ExportTableSchemeResp struct {
	Data  []*models.TableSchema `json:"data"`
	Count int64                 `json:"count"`
}

// ExportTableScheme export table scheme data.
func (b *backup) ExportTableScheme(ctx context.Context, req *ExportTableSchemeReq) (*ExportTableSchemeResp, error) {
	tableSchemes, count, err := b.tableSchemeRepo.List(b.db, &models.TableSchemaQuery{
		AppID: req.AppID,
	},
		req.Page, req.Size)
	if err != nil {
		logger.Logger.WithName("ExportTableScheme").Errorw(err.Error(), header.GetRequestIDKV(ctx).Fuzzy()...)
		return nil, err
	}

	return &ExportTableSchemeResp{
		Data:  tableSchemes,
		Count: count,
	}, nil
}

type ExportRoleReq struct {
	AppID string `uri:"appID"`
	Page  int    `json:"page"`
	Size  int    `json:"size"`
}

type ExportRoleResp struct {
	Data  []*models.Role `json:"data"`
	Count int64          `json:"count"`
}

// ExportRole export role data.
func (b *backup) ExportRole(ctx context.Context, req *ExportRoleReq) (*ExportRoleResp, error) {
	roles, count, err := b.roleRepo.List(b.db, &models.RoleQuery{
		AppID: req.AppID,
	},
		req.Page, req.Size)
	if err != nil {
		logger.Logger.WithName("ExportRole").Errorw(err.Error(), header.GetRequestIDKV(ctx).Fuzzy()...)
		return nil, err
	}

	return &ExportRoleResp{
		Data:  roles,
		Count: count,
	}, nil
}

type ImportTableReq struct {
	Data []*models.Table `json:"data"`
}

type ImportTableResp struct{}

// ImportTable import table data.
func (b *backup) ImportTable(ctx context.Context, req *ImportTableReq) (*ImportTableResp, error) {
	err := b.tableRepo.BatchCreate(b.db, req.Data...)
	if err != nil {
		logger.Logger.WithName("ImportTable").Errorw(err.Error(), header.GetRequestIDKV(ctx).Fuzzy()...)
		return nil, err
	}

	return &ImportTableResp{}, nil
}

type ImportPermitReq struct {
	Data []*models.Permit `json:"data"`
}

type ImportPermitResp struct{}

// ImportPermit import permit data.
func (b *backup) ImportPermit(ctx context.Context, req *ImportPermitReq) (*ImportPermitResp, error) {
	err := b.permitRepo.BatchCreate(b.db, req.Data...)
	if err != nil {
		logger.Logger.WithName("ImportPermit").Errorw(err.Error(), header.GetRequestIDKV(ctx).Fuzzy()...)
		return nil, err
	}

	return &ImportPermitResp{}, nil
}

type ImportTableRelationReq struct {
	Data []*models.TableRelation `json:"data"`
}

type ImportTableRelationResp struct{}

// ImportTableRelation import table relation data.
func (b *backup) ImportTableRelation(ctx context.Context, req *ImportTableRelationReq) (*ImportTableRelationResp, error) {
	err := b.tableRelationRepo.BatchCreate(b.db, req.Data...)
	if err != nil {
		logger.Logger.WithName("ImportTableRelation").Errorw(err.Error(), header.GetRequestIDKV(ctx).Fuzzy()...)
		return nil, err
	}

	return &ImportTableRelationResp{}, nil
}

type ImportTableSchemeReq struct {
	Data []*models.TableSchema `json:"data"`
}

type ImportTableSchemeResp struct{}

// ImportTableScheme import table scheme data.
func (b *backup) ImportTableScheme(ctx context.Context, req *ImportTableSchemeReq) (*ImportTableSchemeResp, error) {
	err := b.tableSchemeRepo.BatchCreate(b.db, req.Data...)
	if err != nil {
		logger.Logger.WithName("ImportTableScheme").Errorw(err.Error(), header.GetRequestIDKV(ctx).Fuzzy()...)
		return nil, err
	}

	return &ImportTableSchemeResp{}, nil
}

type ImportRoleReq struct {
	Data []*models.Role `json:"data"`
}

type ImportRoleResp struct{}

// ImportRole import role data.
func (b *backup) ImportRole(ctx context.Context, req *ImportRoleReq) (*ImportRoleResp, error) {
	err := b.roleRepo.BatchCreate(b.db, req.Data...)
	if err != nil {
		logger.Logger.WithName("ImportRole").Errorw(err.Error(), header.GetRequestIDKV(ctx).Fuzzy()...)
		return nil, err
	}

	return &ImportRoleResp{}, nil
}
