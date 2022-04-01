package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/quanxiang-cloud/cabin/tailormade/header"
	"github.com/quanxiang-cloud/cabin/tailormade/resp"
	"github.com/quanxiang-cloud/form/internal/service"
	"github.com/quanxiang-cloud/form/pkg/misc/config"
)

type Backup struct {
	backup service.Backup
}

func NewBackup(conf *config.Config) (Backup, error) {
	backup, err := service.NewBackup(conf)
	if err != nil {
		return Backup{}, err
	}

	return Backup{
		backup: backup,
	}, nil
}

func (b *Backup) ExportTable(c *gin.Context) {
	req := &service.ExportReq{}

	if err := c.ShouldBind(req); err != nil {
		c.AbortWithError(http.StatusBadGateway, err)
		return
	}

	resp.Format(b.backup.ExportTable(header.MutateContext(c), req)).Context(c)
}

func (b *Backup) ExportPermit(c *gin.Context) {
	req := &service.ExportReq{}

	if err := c.ShouldBind(req); err != nil {
		c.AbortWithError(http.StatusBadGateway, err)
		return
	}

	resp.Format(b.backup.ExportPermit(header.MutateContext(c), req)).Context(c)
}

func (b *Backup) ExportRole(c *gin.Context) {
	req := &service.ExportReq{}

	if err := c.ShouldBind(req); err != nil {
		c.AbortWithError(http.StatusBadGateway, err)
		return
	}

	resp.Format(b.backup.ExportRole(header.MutateContext(c), req)).Context(c)
}

func (b *Backup) ExportTableSchema(c *gin.Context) {
	req := &service.ExportReq{}

	if err := c.ShouldBind(req); err != nil {
		c.AbortWithError(http.StatusBadGateway, err)
		return
	}

	resp.Format(b.backup.ExportTableScheme(header.MutateContext(c), req)).Context(c)
}

func (b *Backup) ExportTableRelation(c *gin.Context) {
	req := &service.ExportReq{}

	if err := c.ShouldBind(req); err != nil {
		c.AbortWithError(http.StatusBadGateway, err)
		return
	}

	resp.Format(b.backup.ExportTableRelation(header.MutateContext(c), req)).Context(c)
}

func (b *Backup) ImportTable(c *gin.Context) {
	req := &service.ImportReq{}

	if err := c.ShouldBind(req); err != nil {
		c.AbortWithError(http.StatusBadGateway, err)
		return
	}

	resp.Format(b.backup.ImportTable(header.MutateContext(c), req)).Context(c)
}

func (b *Backup) ImportPermit(c *gin.Context) {
	req := &service.ImportReq{}

	if err := c.ShouldBind(req); err != nil {
		c.AbortWithError(http.StatusBadGateway, err)
		return
	}

	resp.Format(b.backup.ImportPermit(header.MutateContext(c), req)).Context(c)
}

func (b *Backup) ImportRole(c *gin.Context) {
	req := &service.ImportReq{}

	if err := c.ShouldBind(req); err != nil {
		c.AbortWithError(http.StatusBadGateway, err)
		return
	}

	resp.Format(b.backup.ImportRole(header.MutateContext(c), req)).Context(c)
}

func (b *Backup) ImportTableSchema(c *gin.Context) {
	req := &service.ImportReq{}

	if err := c.ShouldBind(req); err != nil {
		c.AbortWithError(http.StatusBadGateway, err)
		return
	}

	resp.Format(b.backup.ImportTableScheme(header.MutateContext(c), req)).Context(c)
}

func (b *Backup) ImportTableRelation(c *gin.Context) {
	req := &service.ImportReq{}

	if err := c.ShouldBind(req); err != nil {
		c.AbortWithError(http.StatusBadGateway, err)
		return
	}

	resp.Format(b.backup.ImportTableRelation(header.MutateContext(c), req)).Context(c)
}
