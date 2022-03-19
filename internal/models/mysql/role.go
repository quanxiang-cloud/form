package mysql

import (
	"github.com/quanxiang-cloud/form/internal/models"
	"gorm.io/gorm"
)

type roleRepo struct{}

func (t *roleRepo) BatchCreate(db *gorm.DB, role ...*models.Role) error {
	return db.Table(t.TableName()).CreateInBatches(role, len(role)).Error
}

func (t *roleRepo) Get(db *gorm.DB, id string) (*models.Role, error) {
	role := new(models.Role)
	err := db.Table(t.TableName()).Where("id = ? ", id).Find(role).Error
	if err != nil {
		return nil, err
	}
	return role, nil

}

func (t *roleRepo) Find(db *gorm.DB, query *models.RoleQuery) ([]*models.Role, error) {
	ql := db.Table(t.TableName())

	if query.AppID != "" {
		ql = ql.Where("app_id = ?", query.AppID)
	}
	ql = ql.Order("created_at DESC")
	role := make([]*models.Role, 0)
	err := ql.Find(&role).Error
	return role, err
}

func (t *roleRepo) Update(db *gorm.DB, id string, role *models.Role) error {
	setMap := make(map[string]interface{})
	if role.Name != "" {
		setMap["name"] = role.Name
	}
	if role.Description != "" {
		setMap["description"] = role.Description
	}
	return db.Table(t.TableName()).Where("id = ? ", id).Updates(
		setMap).Error
}

func (t *roleRepo) Delete(db *gorm.DB, query *models.RoleQuery) error {
	resp := make([]models.Role, 0)
	ql := db.Table(t.TableName())
	if query.ID != "" {
		ql = ql.Where("id = ?", query.ID)
	}
	return ql.Delete(resp).Error
}

func (t *roleRepo) TableName() string {
	return "role"
}

func NewRoleRepo() models.RoleRepo {
	return &roleRepo{}
}
