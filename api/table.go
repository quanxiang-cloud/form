package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/quanxiang-cloud/cabin/logger"
	"github.com/quanxiang-cloud/cabin/tailormade/header"
	"github.com/quanxiang-cloud/cabin/tailormade/resp"
	"github.com/quanxiang-cloud/form/internal/service"
	config2 "github.com/quanxiang-cloud/form/pkg/misc/config"
)

//Table  table
type Table struct {
	table service.Table
}

// NewTable new table
func NewTable(conf *config2.Config) (*Table, error) {

	t, err := service.NewTable(conf)
	if err != nil {
		return nil, err
	}
	return &Table{
		table: t,
	}, nil
}

func (t *Table) CreateSchema(c *gin.Context) {
	req := &service.CreateTableReq{
		AppID:    c.Param(_appID),
		UserID:   c.GetHeader(_userID),
		UserName: c.GetHeader(_userName),
	}
	ctx := header.MutateContext(c)
	if err := c.ShouldBind(req); err != nil {
		logger.Logger.Errorw("should bind", header.GetRequestIDKV(ctx).Fuzzy(), err.Error())
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	resp.Format(t.table.CreateSchema(ctx, req)).Context(c)

}
