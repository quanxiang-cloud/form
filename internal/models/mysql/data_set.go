package mysql

import (
	"github.com/quanxiang-cloud/form/internal/models"
	"gorm.io/gorm"
)

type dataSetRepo struct{}

func (t *dataSetRepo) Insert(db *gorm.DB, dataset *models.DataSet) error {
	return db.Table(t.TableName()).Create(dataset).Error
}

func (t *dataSetRepo) Update(db *gorm.DB, id string, dataset *models.DataSet) error {
	setMap := make(map[string]interface{})
	if dataset.Name != "" {
		setMap["name"] = dataset.Name
	}
	if dataset.Name != "" {
		setMap["tag"] = dataset.Name
	}
	if dataset.Type != 0 {
		setMap["type"] = dataset.Type
	}
	if dataset.Content != "" {
		setMap["content"] = dataset.Content
	}
	return db.Table(t.TableName()).Where("id = ? ", id).Updates(
		setMap).Error
}

func (t *dataSetRepo) Delete(db *gorm.DB, id string) error {
	resp := make([]models.Permit, 0)
	ql := db.Table(t.TableName()).Where("id = ?", id)
	return ql.Delete(resp).Error
}

func (t *dataSetRepo) GetByID(db *gorm.DB, id string) (*models.DataSet, error) {
	dataSet := new(models.DataSet)
	err := db.Table(t.TableName()).Where("id = ?", id).Find(dataSet).Error
	if err != nil {
		return nil, err
	}
	return dataSet, nil
}

func (t *dataSetRepo) Find(db *gorm.DB, query *models.DataSetQuery, page, size int) ([]*models.DataSet, int64, error) {
	panic("implement me")
}

func (t *dataSetRepo) TableName() string {
	return "data_set"
}

func NewDataSetRepo() models.DataSetRepo {
	return &dataSetRepo{}
}
