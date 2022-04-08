package service

import (
	"context"
	"encoding/json"

	"github.com/quanxiang-cloud/form/internal/models"
	"github.com/quanxiang-cloud/form/internal/models/mysql"
	config2 "github.com/quanxiang-cloud/form/pkg/misc/config"

	"github.com/quanxiang-cloud/cabin/logger"
	"github.com/quanxiang-cloud/cabin/tailormade/header"
	"gorm.io/gorm"
)

// Backup import and export data interface.
type Backup interface {
	ExportTable(context.Context, *ExportReq) (*ExportResp, error)
	ExportPermit(context.Context, *ExportReq) (*ExportResp, error)
	ExportTableRelation(context.Context, *ExportReq) (*ExportResp, error)
	ExportTableScheme(context.Context, *ExportReq) (*ExportResp, error)
	ExportRole(context.Context, *ExportReq) (*ExportResp, error)

	ImportTable(context.Context, *ImportReq) (*ImportResp, error)
	ImportPermit(context.Context, *ImportReq) (*ImportResp, error)
	ImportTableRelation(context.Context, *ImportReq) (*ImportResp, error)
	ImportTableScheme(context.Context, *ImportReq) (*ImportResp, error)
	ImportRole(context.Context, *ImportReq) (*ImportResp, error)
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
		tableRelationRepo: mysql.NewTableRelation(),
		tableSchemeRepo:   mysql.NewTableSchema(),
	}, nil
}

// ExportReq export request.
type ExportReq struct {
	AppID string `uri:"appID"`
	Page  int    `json:"page"`
	Size  int    `json:"size"`
}

// ExportResp export response.
type ExportResp struct {
	Data  Object `json:"data"`
	Count int64  `json:"count"`
}

// Object export data type.
type Object []interface{}

// ExportTable export table data.
func (b *backup) ExportTable(ctx context.Context, req *ExportReq) (*ExportResp, error) {
	tables, count, err := b.tableRepo.List(b.db, &models.TableQuery{
		AppID: req.AppID,
	},
		req.Page, req.Size)
	if err != nil {
		logger.Logger.WithName("ExportTable").Errorw(err.Error(), header.GetRequestIDKV(ctx).Fuzzy()...)
		return nil, err
	}

	data := make(Object, 0, len(tables))
	for _, table := range tables {
		data = append(data, table)
	}

	return &ExportResp{
		Data:  data,
		Count: count,
	}, nil
}

// ExportPermit export permit data.
func (b *backup) ExportPermit(ctx context.Context, req *ExportReq) (*ExportResp, error) {
	roles, count, err := b.roleRepo.List(b.db, &models.RoleQuery{
		AppID: req.AppID,
	},
		req.Page, req.Size)
	if err != nil {
		logger.Logger.WithName("ExportPermit").Errorw(err.Error(), header.GetRequestIDKV(ctx).Fuzzy()...)
		return nil, err
	}

	if len(roles) == 0 {
		return &ExportResp{
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

	data := make(Object, 0, len(permits))
	for _, permit := range permits {
		data = append(data, permit)
	}

	return &ExportResp{
		Data:  data,
		Count: count,
	}, nil
}

// ExportTableRelation export table relation data.
func (b *backup) ExportTableRelation(ctx context.Context, req *ExportReq) (*ExportResp, error) {
	tableRelations, count, err := b.tableRelationRepo.List(b.db, &models.TableRelationQuery{
		AppID: req.AppID,
	},
		req.Page, req.Size)
	if err != nil {
		logger.Logger.WithName("ExportTableRelation").Errorw(err.Error(), header.GetRequestIDKV(ctx).Fuzzy()...)
		return nil, err
	}

	data := make(Object, 0, len(tableRelations))
	for _, tableRelation := range tableRelations {
		data = append(data, tableRelation)
	}

	return &ExportResp{
		Data:  data,
		Count: count,
	}, nil
}

// ExportTableScheme export table scheme data.
func (b *backup) ExportTableScheme(ctx context.Context, req *ExportReq) (*ExportResp, error) {
	tableSchemes, count, err := b.tableSchemeRepo.List(b.db, &models.TableSchemaQuery{
		AppID: req.AppID,
	},
		req.Page, req.Size)
	if err != nil {
		logger.Logger.WithName("ExportTableScheme").Errorw(err.Error(), header.GetRequestIDKV(ctx).Fuzzy()...)
		return nil, err
	}

	data := make(Object, 0, len(tableSchemes))
	for _, tableScheme := range tableSchemes {
		data = append(data, tableScheme)
	}

	return &ExportResp{
		Data:  data,
		Count: count,
	}, nil
}

// ExportRole export role data.
func (b *backup) ExportRole(ctx context.Context, req *ExportReq) (*ExportResp, error) {
	roles, count, err := b.roleRepo.List(b.db, &models.RoleQuery{
		AppID: req.AppID,
	},
		req.Page, req.Size)
	if err != nil {
		logger.Logger.WithName("ExportRole").Errorw(err.Error(), header.GetRequestIDKV(ctx).Fuzzy()...)
		return nil, err
	}

	data := make(Object, 0, len(roles))
	for _, role := range roles {
		data = append(data, role)
	}

	return &ExportResp{
		Data:  data,
		Count: count,
	}, nil
}

// ImportReq import request.
type ImportReq struct {
	Data Object `json:"data"`
}

// ImportResp import response.
type ImportResp struct{}

// ImportTable import table data.
func (b *backup) ImportTable(ctx context.Context, req *ImportReq) (*ImportResp, error) {
	tables := make([]*models.Table, 0, len(req.Data))
	for _, data := range req.Data {
		table := &models.Table{}
		if err := decode(data, table); err != nil {
			logger.Logger.WithName("ImportTable").Errorw(err.Error(), header.GetRequestIDKV(ctx).Fuzzy()...)
			return nil, err
		}
		tables = append(tables, table)
	}

	err := b.tableRepo.BatchCreate(b.db, tables...)
	if err != nil {
		logger.Logger.WithName("ImportTable").Errorw(err.Error(), header.GetRequestIDKV(ctx).Fuzzy()...)
		return nil, err
	}

	return &ImportResp{}, nil
}

// ImportPermit import permit data.
func (b *backup) ImportPermit(ctx context.Context, req *ImportReq) (*ImportResp, error) {
	permits := make([]*models.Permit, 0, len(req.Data))
	for _, data := range req.Data {
		permit := &models.Permit{}
		if err := decode(data, permit); err != nil {
			logger.Logger.WithName("ImportPermit").Errorw(err.Error(), header.GetRequestIDKV(ctx).Fuzzy()...)
			return nil, err
		}
		permits = append(permits, permit)
	}

	err := b.permitRepo.BatchCreate(b.db, permits...)
	if err != nil {
		logger.Logger.WithName("ImportPermit").Errorw(err.Error(), header.GetRequestIDKV(ctx).Fuzzy()...)
		return nil, err
	}

	return &ImportResp{}, nil
}

// ImportTableRelation import table relation data.
func (b *backup) ImportTableRelation(ctx context.Context, req *ImportReq) (*ImportResp, error) {
	tableRelations := make([]*models.TableRelation, 0, len(req.Data))
	for _, data := range req.Data {
		tableRelation := &models.TableRelation{}
		if err := decode(data, tableRelation); err != nil {
			logger.Logger.WithName("ImportTableRelation").Errorw(err.Error(), header.GetRequestIDKV(ctx).Fuzzy()...)
			return nil, err
		}
		tableRelations = append(tableRelations, tableRelation)
	}

	err := b.tableRelationRepo.BatchCreate(b.db, tableRelations...)
	if err != nil {
		logger.Logger.WithName("ImportTableRelation").Errorw(err.Error(), header.GetRequestIDKV(ctx).Fuzzy()...)
		return nil, err
	}

	return &ImportResp{}, nil
}

// ImportTableScheme import table scheme data.
func (b *backup) ImportTableScheme(ctx context.Context, req *ImportReq) (*ImportResp, error) {
	tableSchemes := make([]*models.TableSchema, 0, len(req.Data))
	for _, data := range req.Data {
		tableScheme := &models.TableSchema{}
		if err := decode(data, tableScheme); err != nil {
			logger.Logger.WithName("ImportTableScheme").Errorw(err.Error(), header.GetRequestIDKV(ctx).Fuzzy()...)
			return nil, err
		}
		tableSchemes = append(tableSchemes, tableScheme)
	}

	err := b.tableSchemeRepo.BatchCreate(b.db, tableSchemes...)
	if err != nil {
		logger.Logger.WithName("ImportTableScheme").Errorw(err.Error(), header.GetRequestIDKV(ctx).Fuzzy()...)
		return nil, err
	}

	return &ImportResp{}, nil
}

// ImportRole import role data.
func (b *backup) ImportRole(ctx context.Context, req *ImportReq) (*ImportResp, error) {
	roles := make([]*models.Role, 0, len(req.Data))
	for _, data := range req.Data {
		role := &models.Role{}
		if err := decode(data, role); err != nil {
			logger.Logger.WithName("ImportRole").Errorw(err.Error(), header.GetRequestIDKV(ctx).Fuzzy()...)
			return nil, err
		}
		roles = append(roles, role)
	}

	err := b.roleRepo.BatchCreate(b.db, roles...)
	if err != nil {
		logger.Logger.WithName("ImportRole").Errorw(err.Error(), header.GetRequestIDKV(ctx).Fuzzy()...)
		return nil, err
	}

	return &ImportResp{}, nil
}

func decode(data, obj interface{}) error {
	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return json.Unmarshal(bytes, obj)
}
