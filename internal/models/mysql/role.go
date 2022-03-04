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
	return nil, nil
}

func (t *roleRepo) Update(db *gorm.DB, id string, role *models.Role) error {
	return nil
}

func (t *roleRepo) Delete(db *gorm.DB, query *models.RoleQuery) error {
	return nil
}

func (t *roleRepo) TableName() string {
	return "role"
}

func NewRoleRepo() models.RoleRepo {
	return &roleRepo{}
}
