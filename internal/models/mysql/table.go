package mysql

import (
	"github.com/quanxiang-cloud/form/internal/models"
	"gorm.io/gorm"
)

type tableRepo struct{}

func (t *tableRepo) TableName() string {
	return "table"
}

func NewTableRepo() models.TableRepo {
	return &tableRepo{}
}

func (t *tableRepo) BatchCreate(db *gorm.DB, tables ...*models.Table) error {
	return db.Table(t.TableName()).CreateInBatches(tables, len(tables)).Error
}

func (t *tableRepo) Get(db *gorm.DB, appId, tableID string) (*models.Table, error) {
	table := new(models.Table)
	err := db.Table(t.TableName()).Where("app_id = ? and  table_id = ? ", appId, tableID).Find(table).Error
	if err != nil {
		return nil, err
	}
	return table, nil
}

func (t *tableRepo) Find(db *gorm.DB, query *models.TableQuery) ([]*models.Table, error) {
	return nil, nil
}

func (t *tableRepo) Delete(db *gorm.DB, query *models.TableQuery) error {
	return nil
}

func (t *tableRepo) Update(db *gorm.DB, appID, tableID string, table *models.Table) error {
	setMap := make(map[string]interface{})
	if table.Schema != nil {
		setMap["schema"] = table.Schema
	}
	if table.Config != nil {
		setMap["config"] = table.Config
	}
	return db.Table(t.TableName()).Where("app_id = ? and table_id = ? ", appID, tableID).Updates(
		setMap).Error
}
