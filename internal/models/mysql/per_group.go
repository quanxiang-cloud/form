package mysql

import (
	"github.com/quanxiang-cloud/form/internal/models"
	"gorm.io/gorm"
)

type perGroupRepo struct{}

func (p *perGroupRepo) BatchCreate(db *gorm.DB, permissionGroup ...*models.PerGroup) error {
	return db.Table(p.TableName()).CreateInBatches(permissionGroup, len(permissionGroup)).Error
}

func (p *perGroupRepo) Get(db *gorm.DB, id string) (*models.PerGroup, error) {
	return nil, nil
}

func (p *perGroupRepo) Find(db *gorm.DB, query *models.PerGroupQuery) ([]*models.PerGroup, error) {
	ql := db.Table(p.TableName())
	if query.AppID != "" {
		ql = ql.Where("app_id = ?", query.AppID)
	}
	perList := make([]*models.PerGroup, 0)
	err := ql.Find(&perList).Error
	if err != nil {
		return nil, err
	}
	return perList, nil
}

func (p *perGroupRepo) Update(db *gorm.DB, id string, perGroup *models.PerGroup) error {
	setMap := make(map[string]interface{})
	if perGroup.Description != "" {
		setMap["description"] = perGroup.Description
	}
	if perGroup.Name != "" {
		setMap["name"] = perGroup.Name
	}
	return db.Table(p.TableName()).Where("id = ?", id).Updates(
		setMap).Error
}

func (p *perGroupRepo) Delete(db *gorm.DB, query *models.PerGroupQuery) error {
	ql := db.Table(p.TableName())
	if query.AppID != "" {
		ql = ql.Where("app_id = ?", query.AppID)
	}
	if len(query.PerGroupIDs) != 0 {
		ql = ql.Where("id in ?", query.PerGroupIDs)
	}
	if query.ID != "" {
		ql = ql.Where("id  =  ?", query.ID)
	}
	return ql.Delete(&models.PerGroup{}).Error
}

func (p *perGroupRepo) TableName() string {
	return "permit_group"
}

func NewPerGroupRepo() models.PerGroupRepo {
	return &perGroupRepo{}
}
