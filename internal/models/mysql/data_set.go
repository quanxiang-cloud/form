package mysql

import (
	"github.com/quanxiang-cloud/form/internal/models"
	"gorm.io/gorm"
)

type dataSetRepo struct{}

func (t *dataSetRepo) Insert(db *gorm.DB, dataset *models.DataSet) error {
	return nil
}

func (t *dataSetRepo) Update(db *gorm.DB, dataset *models.DataSet) error {
	panic("implement me")
}

func (t *dataSetRepo) Delete(db *gorm.DB, id string) error {
	panic("implement me")
}

func (t *dataSetRepo) GetByID(db *gorm.DB, id string) (*models.DataSet, error) {
	panic("implement me")
}

func (t *dataSetRepo) Find(db *gorm.DB, query *models.DataSetQuery) ([]*models.DataSet, error) {
	panic("implement me")
}

func (t *dataSetRepo) TableName() string {
	return "data_set"
}

func NewDataSetRepo() models.DataSetRepo {
	return &dataSetRepo{}
}
