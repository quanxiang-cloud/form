package service

import (
	"context"

	error2 "github.com/quanxiang-cloud/cabin/error"
	id2 "github.com/quanxiang-cloud/cabin/id"
	time2 "github.com/quanxiang-cloud/cabin/time"

	"github.com/quanxiang-cloud/form/internal/models"
	"github.com/quanxiang-cloud/form/pkg/misc/code"
	"github.com/quanxiang-cloud/form/pkg/misc/config"
	"go.mongodb.org/mongo-driver/mongo"
)

// DataSet DataSet
type DataSet interface {
	CreateDataSet(c context.Context, req *CreateDataSetReq) (*CreateDataSetResp, error)
	GetDataSet(c context.Context, req *GetDataSetReq) (*GetDataSetResp, error)
	UpdateDataSet(c context.Context, req *UpdateDataSetReq) (*UpdateDataSetResp, error)
	GetByConditionSet(c context.Context, req *GetByConditionSetReq) (*GetByConditionSetResp, error)
	DeleteDataSet(c context.Context, req *DeleteDataSetReq) (*DeleteDataSetResp, error)
}

type dataset struct {
	mongodb     *mongo.Database
	datasetRepo models.DataSetRepo
}

// NewDataSet NewDataSet
func NewDataSet(conf *config.Config) (DataSet, error) {
	u := &dataset{
		//datasetRepo: repo.NewDataSetRepo(),
	}

	return u, nil
}

// SetMongo SetMongo
func (per *dataset) SetMongo(client *mongo.Client, dbName string) {
	per.mongodb = client.Database(dbName)
}

// CreateDataSetReq CreateDataSetReq
type CreateDataSetReq struct {
	Name    string `json:"name" binding:"max=100"`
	Tag     string `json:"tag"  binding:"max=100"`
	Type    int64  `json:"type"`
	Content string `json:"content"`
}

// CreateDataSetResp CreateDataSetResp
type CreateDataSetResp struct {
	ID string `json:"id"`
}

// CreateDataSet CreateDataSet
func (per *dataset) CreateDataSet(c context.Context, req *CreateDataSetReq) (*CreateDataSetResp, error) {
	exist, err := per.datasetRepo.Find(c, per.mongodb, &models.DataSetQuery{Name: req.Name})
	if err != nil {
		return nil, err
	}
	if len(exist) > 0 {
		return nil, error2.New(code.ErrExistDataSetNameState)
	}
	dataset := &models.DataSet{
		ID:        id2.StringUUID(),
		Name:      req.Name,
		Tag:       req.Tag,
		Type:      req.Type,
		Content:   req.Content,
		CreatedAt: time2.NowUnix(),
	}
	err = per.datasetRepo.Insert(c, per.mongodb, dataset)
	if err != nil {
		return nil, err
	}
	return &CreateDataSetResp{
		ID: dataset.ID,
	}, nil
}

// GetDataSetReq GetDataSetReq
type GetDataSetReq struct {
	ID string `json:"id"`
}

// GetDataSetResp GetDataSetResp
type GetDataSetResp struct {
	ID      string `json:"id"`
	Name    string `json:"name" binding:"max=100"`
	Tag     string `json:"tag"  binding:"max=100"`
	Type    int64  `json:"type"`
	Content string `json:"content"`
}

// GetDataSet GetDataSet
func (per *dataset) GetDataSet(c context.Context, req *GetDataSetReq) (*GetDataSetResp, error) {
	dataset, err := per.datasetRepo.GetByID(c, per.mongodb, req.ID)
	if err != nil {
		return nil, err
	}
	if dataset == nil {
		//return nil, error2.NewError(code.ErrNODataSetNameState)
	}
	resp := &GetDataSetResp{
		ID:      dataset.ID,
		Name:    dataset.Name,
		Tag:     dataset.Tag,
		Type:    dataset.Type,
		Content: dataset.Content,
	}
	return resp, nil
}

// UpdateDataSetReq UpdateDataSetReq
type UpdateDataSetReq struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Tag     string `json:"tag"`
	Type    int64  `json:"type"`
	Content string `json:"content"`
}

// UpdateDataSetResp UpdateDataSetResp
type UpdateDataSetResp struct {
}

// UpdateDataSet UpdateDataSet
func (per *dataset) UpdateDataSet(c context.Context, req *UpdateDataSetReq) (*UpdateDataSetResp, error) {
	data, err := per.datasetRepo.Find(c, per.mongodb, &models.DataSetQuery{
		Name: req.Name,
	})
	if err != nil {
		return nil, err
	}
	for _, value := range data {
		if value.ID != req.ID {
			return nil, error2.New(code.ErrExistDataSetNameState)
		}
	}
	dataset := &models.DataSet{
		ID:      req.ID,
		Name:    req.Name,
		Tag:     req.Tag,
		Type:    req.Type,
		Content: req.Content,
	}
	err = per.datasetRepo.Update(c, per.mongodb, dataset)
	if err != nil {
		return nil, err
	}
	return &UpdateDataSetResp{}, nil
}

// GetByConditionSetReq GetByConditionSetReq
type GetByConditionSetReq struct {
	Name  string `json:"name"`
	Tag   string `json:"tag"`
	Types int64  `json:"type"`
}

// GetByConditionSetResp GetByConditionSetResp
type GetByConditionSetResp struct {
	List []*DataSetVo `json:"list"`
}

// DataSetVo DataSetVo
type DataSetVo struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Tag     string `json:"tag"`
	Type    int64  `json:"type"`
	Content string `json:"content"`
}

// GetByConditionSet GetByConditionSet
func (per *dataset) GetByConditionSet(c context.Context, req *GetByConditionSetReq) (*GetByConditionSetResp, error) {
	arr, err := per.datasetRepo.Find(c, per.mongodb, &models.DataSetQuery{
		Tag:   req.Tag,
		Name:  req.Name,
		Types: req.Types,
	})
	if err != nil {
		return nil, err
	}
	resp := &GetByConditionSetResp{
		List: make([]*DataSetVo, len(arr)),
	}
	for index, value := range arr {
		resp.List[index] = new(DataSetVo)
		cloneDataSet(value, resp.List[index])
	}
	return resp, nil

}
func cloneDataSet(src *models.DataSet, dst *DataSetVo) {
	dst.ID = src.ID
	dst.Name = src.Name
	dst.Tag = src.Tag
	dst.Type = src.Type
	dst.Content = src.Content

}

// DeleteDataSetReq DeleteDataSetReq
type DeleteDataSetReq struct {
	ID string `json:"id"`
}

// DeleteDataSetResp DeleteDataSetResp
type DeleteDataSetResp struct {
}

// DeleteDataSet DeleteDataSet
func (per *dataset) DeleteDataSet(c context.Context, req *DeleteDataSetReq) (*DeleteDataSetResp, error) {
	err := per.datasetRepo.Delete(c, per.mongodb, req.ID)
	if err != nil {
		return nil, err
	}
	return &DeleteDataSetResp{}, nil
}
