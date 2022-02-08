package api

import (
	"github.com/gin-gonic/gin"
	"github.com/quanxiang-cloud/cabin/tailormade/header"
	"github.com/quanxiang-cloud/cabin/tailormade/resp"
	"github.com/quanxiang-cloud/form/internal/service"
	"github.com/quanxiang-cloud/form/pkg/misc/config"
	"net/http"
)

type DataSet struct {
	dataset service.DataSet
}

// NewDataSet 初始化
func NewDataSet(conf *config.Config) (*DataSet, error) {
	d, err := service.NewDataSet(conf)
	if err != nil {
		return nil, err
	}
	return &DataSet{
		dataset: d,
	}, nil
}

// CreateDataSet CreateDataSet
func (d *DataSet) CreateDataSet(c *gin.Context) {
	req := &service.CreateDataSetReq{}
	if err := c.ShouldBind(req); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	resp.Format(d.dataset.CreateDataSet(header.MutateContext(c), req)).Context(c)

}

// GetDataSet GetDataSet
func (d *DataSet) GetDataSet(c *gin.Context) {
	req := &service.GetDataSetReq{}
	if err := c.ShouldBind(req); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	resp.Format(d.dataset.GetDataSet(header.MutateContext(c), req)).Context(c)
}

// UpdateDataSet UpdateDataSet
func (d *DataSet) UpdateDataSet(c *gin.Context) {
	req := &service.UpdateDataSetReq{}
	if err := c.ShouldBind(req); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	resp.Format(d.dataset.UpdateDataSet(header.MutateContext(c), req)).Context(c)

}

// GetByConditionSet GetByConditionSet
func (d *DataSet) GetByConditionSet(c *gin.Context) {
	req := &service.GetByConditionSetReq{}
	if err := c.ShouldBind(req); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	resp.Format(d.dataset.GetByConditionSet(header.MutateContext(c), req)).Context(c)
}

// DeleteDataSet DeleteDataSet
func (d *DataSet) DeleteDataSet(c *gin.Context) {
	req := &service.DeleteDataSetReq{}
	if err := c.ShouldBind(req); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	resp.Format(d.dataset.DeleteDataSet(header.MutateContext(c), req)).Context(c)

}
