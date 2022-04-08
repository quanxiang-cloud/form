package api

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	id2 "github.com/quanxiang-cloud/cabin/id"
	"github.com/quanxiang-cloud/cabin/logger"
	"github.com/quanxiang-cloud/cabin/tailormade/header"
	"github.com/quanxiang-cloud/cabin/tailormade/resp"
	table2 "github.com/quanxiang-cloud/form/internal/service/tables"
	config2 "github.com/quanxiang-cloud/form/pkg/misc/config"
)

// Table  table.
type Table struct {
	table    table2.Table
	guidance table2.Guidance
}

// NewTable new table.
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

// CrateTable create table.
func (t *Table) CrateTable(c *gin.Context) {
	profiles := getProfile(c)
	req := &table2.Bus{
		AppID:    c.Param(_appID),
		UserID:   profiles.userID,
		UserName: profiles.userName,
	}
	ctx := header.MutateContext(c)
	if err := c.ShouldBind(req); err != nil {
		logger.Logger.WithName("CrateTable").Errorw(err.Error(), header.GetRequestIDKV(ctx).Fuzzy()...)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	resp.Format(t.guidance.Do(ctx, req)).Context(c)
}

// GetTable GetTable.
func (t *Table) GetTable(c *gin.Context) {
	req := &table2.GetTableReq{
		AppID:   c.Param(_appID),
		TableID: c.Param("tableName"),
	}
	ctx := header.MutateContext(c)
	if err := c.ShouldBind(req); err != nil {
		logger.Logger.WithName("GetTable").Errorw(err.Error(), header.GetRequestIDKV(ctx).Fuzzy()...)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	resp.Format(t.table.GetTable(ctx, req)).Context(c)
}

// DeleteTable DeleteTable.
func (t *Table) DeleteTable(c *gin.Context) {
	req := &table2.DeleteTableReq{
		AppID: c.Param("appID"),
	}

	ctx := header.MutateContext(c)
	if err := c.ShouldBind(req); err != nil {
		logger.Logger.WithName("DeleteTable").Errorw(err.Error(), header.GetRequestIDKV(ctx).Fuzzy()...)
		return
	}
	resp.Format(t.table.DeleteTable(ctx, req)).Context(c)
}

// CreateBlank create blank.
func (t *Table) CreateBlank(c *gin.Context) {
	resp.Format(struct {
		TableID string `json:"tableID"`
	}{
		TableID: id2.String(5),
	}, nil).Context(c)
}

// FindTable find table.
func (t *Table) FindTable(c *gin.Context) {
	req := &table2.FindTableReq{
		AppID: c.Param("appID"),
	}
	ctx := header.MutateContext(c)
	if err := c.ShouldBind(req); err != nil {
		logger.Logger.WithName("FindTable").Errorw(err.Error(), header.GetRequestIDKV(ctx).Fuzzy()...)
		return
	}
	resp.Format(t.table.FindTable(ctx, req)).Context(c)
}

// UpdateConfig update config.
func (t *Table) UpdateConfig(c *gin.Context) {
	req := &table2.UpdateConfigReq{
		AppID: c.Param("appID"),
	}
	ctx := header.MutateContext(c)
	if err := c.ShouldBind(req); err != nil {
		logger.Logger.WithName("UpdateConfig").Errorw(err.Error(), header.GetRequestIDKV(ctx).Fuzzy()...)
		return
	}
	resp.Format(t.table.UpdateConfig(ctx, req)).Context(c)
}

func (t *Table) GetTableInfo(c *gin.Context) {
	req := &table2.GetTableInfoReq{
		AppID: c.Param("appID"),
	}
	ctx := header.MutateContext(c)
	if err := c.ShouldBind(req); err != nil {
		logger.Logger.WithName("api table").Errorw(err.Error(), header.GetRequestIDKV(ctx).Fuzzy()...)
		return
	}
	resp.Format(t.table.GetTableInfo(ctx, req)).Context(c)

}

// GetXNameReq GetXNameReq
type GetXNameReq struct {
	TableID string `json:"TableID"`
	Action  string `json:"action"`
	AppID   string `json:"AppID"`
}

func (t *Table) GetXName(c *gin.Context) {
	req := &GetXNameReq{}
	req.AppID = c.Param("appID")
	if err := c.ShouldBind(req); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	name := GenXName(req.AppID, req.TableID, req.Action)
	resp.Format(map[string]string{
		"name": name,
	}, nil).Context(c)
}

func getProfile(c *gin.Context) *profile {
	depIDS := strings.Split(c.GetHeader(_departmentID), ",")
	return &profile{
		userID:   c.Param(_userID),
		userName: c.Param(_userName),
		depID:    depIDS[0],
	}
}
func GenXName(appID, tableID, tag string) string {
	return fmt.Sprintf("/system/app/%s/raw/inner/%s/%s/%s.r", appID, "form", tableID, GetInnerXName(tableID, tag))
}

// GetInnerXName GetInnerXName
func GetInnerXName(tableID, tag string) string {
	tableIDs := strings.Split(tableID, "_")
	return fmt.Sprintf("%s_%s", tableIDs[len(tableIDs)-1], tag)
}
func GenXName(appID, tableID, tag string) string {
	return fmt.Sprintf("/system/app/%s/raw/inner/%s/%s/%s.r", appID, "form", tableID, GetInnerXName(tableID, tag))
}

// GetInnerXName GetInnerXName
func GetInnerXName(tableID, tag string) string {
	tableIDs := strings.Split(tableID, "_")
	return fmt.Sprintf("%s_%s", tableIDs[len(tableIDs)-1], tag)
}
