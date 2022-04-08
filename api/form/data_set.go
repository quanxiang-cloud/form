package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/quanxiang-cloud/cabin/logger"
	"github.com/quanxiang-cloud/cabin/tailormade/header"
	"github.com/quanxiang-cloud/cabin/tailormade/resp"
	"github.com/quanxiang-cloud/form/internal/service"
	"github.com/quanxiang-cloud/form/pkg/misc/config"
)

type DataSet struct {
	dataset service.DataSet
}

// Newdataset initialization.
func NewDataSet(conf *config.Config) (*DataSet, error) {
	d, err := service.NewDataSet(conf)
	if err != nil {
		return nil, err
	}
	return &DataSet{
		dataset: d,
	}, nil
}

// CreateDataSet CreateDataSet.
func (d *DataSet) CreateDataSet(c *gin.Context) {
	ctx := header.MutateContext(c)

	req := &service.CreateDataSetReq{}
	if err := c.ShouldBind(req); err != nil {
		logger.Logger.WithName("CreateDataSet").Errorw(err.Error(), header.GetRequestIDKV(ctx).Fuzzy()...)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	resp.Format(d.dataset.CreateDataSet(ctx, req)).Context(c)
}

// GetDataSet GetDataSet.
func (d *DataSet) GetDataSet(c *gin.Context) {
	ctx := header.MutateContext(c)

	req := &service.GetDataSetReq{}
	if err := c.ShouldBind(req); err != nil {
		logger.Logger.WithName("GetDataSet").Errorw(err.Error(), header.GetRequestIDKV(ctx).Fuzzy()...)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	resp.Format(d.dataset.GetDataSet(ctx, req)).Context(c)
}

// UpdateDataSet UpdateDataSet.
func (d *DataSet) UpdateDataSet(c *gin.Context) {
	ctx := header.MutateContext(c)

	req := &service.UpdateDataSetReq{}
	if err := c.ShouldBind(req); err != nil {
		logger.Logger.WithName("UpdateDataSet").Errorw(err.Error(), header.GetRequestIDKV(ctx).Fuzzy()...)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	resp.Format(d.dataset.UpdateDataSet(ctx, req)).Context(c)
}

// GetByConditionSet GetByConditionSet.
func (d *DataSet) GetByConditionSet(c *gin.Context) {
	ctx := header.MutateContext(c)

	req := &service.GetByConditionSetReq{}
	if err := c.ShouldBind(req); err != nil {
		logger.Logger.WithName("GetByConditionSet").Errorw(err.Error(), header.GetRequestIDKV(ctx).Fuzzy()...)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	resp.Format(d.dataset.GetByConditionSet(ctx, req)).Context(c)
}

// DeleteDataSet DeleteDataSet.
func (d *DataSet) DeleteDataSet(c *gin.Context) {
	ctx := header.MutateContext(c)

	req := &service.DeleteDataSetReq{}
	if err := c.ShouldBind(req); err != nil {
		logger.Logger.WithName("DeleteDataSet").Errorw(err.Error(), header.GetRequestIDKV(ctx).Fuzzy()...)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	resp.Format(d.dataset.DeleteDataSet(ctx, req)).Context(c)
}
