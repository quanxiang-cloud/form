package api

import (
	"net/http"

	"github.com/quanxiang-cloud/form/internal/service"
	"github.com/quanxiang-cloud/form/pkg/misc/config"

	"github.com/gin-gonic/gin"
	"github.com/quanxiang-cloud/cabin/logger"
	"github.com/quanxiang-cloud/cabin/tailormade/header"
	"github.com/quanxiang-cloud/cabin/tailormade/resp"
)

// Backup backup.
type Backup struct {
	backup service.Backup
}

// NewBackup new backup.
func NewBackup(conf *config.Config) (Backup, error) {
	backup, err := service.NewBackup(conf)
	if err != nil {
		return Backup{}, err
	}

	return Backup{
		backup: backup,
	}, nil
}

// ExportTable export table.
func (b *Backup) ExportTable(c *gin.Context) {
	req := &service.ExportTableReq{}

	ctx := header.MutateContext(c)
	if err := c.ShouldBind(req); err != nil {
		logger.Logger.WithName("ExportTable").Errorw(err.Error(), header.GetRequestIDKV(ctx).Fuzzy()...)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	resp.Format(b.backup.ExportTable(ctx, req)).Context(c)
}

// ExportPermit export permit.
func (b *Backup) ExportPermit(c *gin.Context) {
	req := &service.ExportPermitReq{}

	ctx := header.MutateContext(c)
	if err := c.ShouldBind(req); err != nil {
		logger.Logger.WithName("ExportPermit").Errorw(err.Error(), header.GetRequestIDKV(ctx).Fuzzy()...)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	resp.Format(b.backup.ExportPermit(ctx, req)).Context(c)
}

// ExportRole export role.
func (b *Backup) ExportRole(c *gin.Context) {
	req := &service.ExportRoleReq{}

	ctx := header.MutateContext(c)
	if err := c.ShouldBind(req); err != nil {
		logger.Logger.WithName("ExportRole").Errorw(err.Error(), header.GetRequestIDKV(ctx).Fuzzy()...)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	resp.Format(b.backup.ExportRole(ctx, req)).Context(c)
}

// ExportTableSchema export table schema.
func (b *Backup) ExportTableSchema(c *gin.Context) {
	req := &service.ExportTableSchemeReq{}

	ctx := header.MutateContext(c)
	if err := c.ShouldBind(req); err != nil {
		logger.Logger.WithName("ExportTableSchema").Errorw(err.Error(), header.GetRequestIDKV(ctx).Fuzzy()...)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	resp.Format(b.backup.ExportTableScheme(ctx, req)).Context(c)
}

// ExportTableRelation export table relation.
func (b *Backup) ExportTableRelation(c *gin.Context) {
	req := &service.ExportTableRelationReq{}

	ctx := header.MutateContext(c)
	if err := c.ShouldBind(req); err != nil {
		logger.Logger.WithName("ExportTableRelation").Errorw(err.Error(), header.GetRequestIDKV(ctx).Fuzzy()...)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	resp.Format(b.backup.ExportTableRelation(ctx, req)).Context(c)
}

// ImportTable import table.
func (b *Backup) ImportTable(c *gin.Context) {
	req := &service.ImportTableReq{}

	ctx := header.MutateContext(c)
	if err := c.ShouldBind(req); err != nil {
		logger.Logger.WithName("ImportTable").Errorw(err.Error(), header.GetRequestIDKV(ctx).Fuzzy()...)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	resp.Format(b.backup.ImportTable(ctx, req)).Context(c)
}

// ImportPermit import permit.
func (b *Backup) ImportPermit(c *gin.Context) {
	req := &service.ImportPermitReq{}

	ctx := header.MutateContext(c)
	if err := c.ShouldBind(req); err != nil {
		logger.Logger.WithName("ImportPermit").Errorw(err.Error(), header.GetRequestIDKV(ctx).Fuzzy()...)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	resp.Format(b.backup.ImportPermit(ctx, req)).Context(c)
}

// ImportRole import role.
func (b *Backup) ImportRole(c *gin.Context) {
	req := &service.ImportRoleReq{}

	ctx := header.MutateContext(c)
	if err := c.ShouldBind(req); err != nil {
		logger.Logger.WithName("ImportRole").Errorw(err.Error(), header.GetRequestIDKV(ctx).Fuzzy()...)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	resp.Format(b.backup.ImportRole(ctx, req)).Context(c)
}

// ImportTableSchema import table schema.
func (b *Backup) ImportTableSchema(c *gin.Context) {
	req := &service.ImportTableSchemeReq{}

	ctx := header.MutateContext(c)
	if err := c.ShouldBind(req); err != nil {
		logger.Logger.WithName("ImportTableSchema").Errorw(err.Error(), header.GetRequestIDKV(ctx).Fuzzy()...)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	resp.Format(b.backup.ImportTableScheme(ctx, req)).Context(c)
}

// ImportTableRelation import table relation.
func (b *Backup) ImportTableRelation(c *gin.Context) {
	req := &service.ImportTableRelationReq{}

	ctx := header.MutateContext(c)
	if err := c.ShouldBind(req); err != nil {
		logger.Logger.WithName("ImportTableRelation").Errorw(err.Error(), header.GetRequestIDKV(ctx).Fuzzy()...)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	resp.Format(b.backup.ImportTableRelation(ctx, req)).Context(c)
}
