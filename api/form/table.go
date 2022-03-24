package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	id2 "github.com/quanxiang-cloud/cabin/id"
	"github.com/quanxiang-cloud/cabin/logger"
	"github.com/quanxiang-cloud/cabin/tailormade/header"
	"github.com/quanxiang-cloud/cabin/tailormade/resp"
	table2 "github.com/quanxiang-cloud/form/internal/service/tables"
	config2 "github.com/quanxiang-cloud/form/pkg/misc/config"
)

// Table  table
type Table struct {
	table    table2.Table
	guidance table2.Guidance
}

// NewTable new table
func NewTable(conf *config2.Config) (*Table, error) {
	t, err := table2.NewTable(conf)
	if err != nil {
		return nil, err
	}
	guidance, err := table2.NewWebTable(conf)
	if err != nil {
		return nil, err
	}
	return &Table{
		table:    t,
		guidance: guidance,
	}, nil
}

func (t *Table) CrateTable(c *gin.Context) {
	profiles := getProfile(c)
	req := &table2.Bus{
		AppID:    c.Param(_appID),
		UserID:   profiles.userID,
		UserName: profiles.userName,
	}
	ctx := header.MutateContext(c)
	if err := c.ShouldBind(req); err != nil {
		logger.Logger.WithName("api table").Errorw(err.Error(), header.GetRequestIDKV(ctx).Fuzzy()...)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	resp.Format(t.guidance.Do(ctx, req)).Context(c)
}

func (t *Table) GetTable(c *gin.Context) {
	req := &table2.GetTableReq{
		AppID: c.Param(_appID),
	}
	ctx := header.MutateContext(c)
	if err := c.ShouldBind(req); err != nil {
		logger.Logger.WithName("api table").Errorw(err.Error(), header.GetRequestIDKV(ctx).Fuzzy()...)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	resp.Format(t.table.GetTable(ctx, req)).Context(c)
}

func (t *Table) DeleteTable(c *gin.Context) {
	req := &table2.DeleteTableReq{
		AppID: c.Param("appID"),
	}

	ctx := header.MutateContext(c)
	if err := c.ShouldBind(req); err != nil {
		logger.Logger.WithName("api table").Errorw(err.Error(), header.GetRequestIDKV(ctx).Fuzzy()...)
		return
	}
	resp.Format(t.table.DeleteTable(ctx, req)).Context(c)
}

func (t *Table) CreateBlank(c *gin.Context) {
	resp.Format(struct {
		TableID string `json:"tableID"`
	}{
		TableID: id2.String(5),
	}, nil).Context(c)
}

func (t *Table) FindTable(c *gin.Context) {
	req := &table2.FindTableReq{
		AppID: c.Param("appID"),
	}
	ctx := header.MutateContext(c)
	if err := c.ShouldBind(req); err != nil {
		logger.Logger.WithName("api table").Errorw(err.Error(), header.GetRequestIDKV(ctx).Fuzzy()...)
		return
	}
	resp.Format(t.table.FindTable(ctx, req)).Context(c)
}

func (t *Table) UpdateConfig(c *gin.Context) {
	req := &table2.UpdateConfigReq{
		AppID: c.Param("appID"),
	}
	ctx := header.MutateContext(c)
	if err := c.ShouldBind(req); err != nil {
		logger.Logger.WithName("api table").Errorw(err.Error(), header.GetRequestIDKV(ctx).Fuzzy()...)
		return
	}
	resp.Format(t.table.UpdateConfig(ctx, req)).Context(c)
}
