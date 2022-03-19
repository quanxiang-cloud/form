package mysql

import (
	"github.com/quanxiang-cloud/form/internal/models"
	"gorm.io/gorm"
)

type permitRepo struct{}

func (t *permitRepo) BatchCreate(db *gorm.DB, permits ...*models.Permit) error {
	return db.Table(t.TableName()).CreateInBatches(permits, len(permits)).Error
}

func (t *permitRepo) Get(db *gorm.DB, roleID, path string) (*models.Permit, error) {
	permits := new(models.Permit)
	err := db.Table(t.TableName()).Where("role_id = ? and  path = ? ", roleID, path).Find(permits).Error
	if err != nil {
		return nil, err
	}
	return permits, nil
}

func (t *permitRepo) Find(db *gorm.DB, query *models.PermitQuery) ([]*models.Permit, error) {
	return nil, nil
}

func (t *permitRepo) Delete(db *gorm.DB, query *models.PermitQuery) error {
	resp := make([]models.Permit, 0)
	ql := db.Table(t.TableName())
	if query.RoleID != "" {
		ql = ql.Where("role_id = ? ", query.RoleID)
	}
	if query.ID != "" {
		ql = ql.Where("id = ?", query.ID)
	}
	return ql.Delete(resp).Error
}

func (t *permitRepo) Update(db *gorm.DB, id string, permit *models.Permit) error {
	setMap := make(map[string]interface{})
	if permit.Params != nil {
		setMap["params"] = permit.Params
	}
	if permit.Response != nil {
		setMap["response"] = permit.Response
	}
	return db.Table(t.TableName()).Where("id = ? ", id).Updates(
		setMap).Error
	return nil
}

func (t *permitRepo) TableName() string {
	return "permit"
}

func NewPermitRepo() models.PermitRepo {
	return &permitRepo{}
}
