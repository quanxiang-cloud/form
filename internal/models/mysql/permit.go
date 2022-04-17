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

func (t *permitRepo) Delete(db *gorm.DB, query *models.PermitQuery) error {
	resp := make([]models.Permit, 0)
	ql := db.Table(t.TableName())
	if query.RoleID != "" {
		ql = ql.Where("role_id = ? ", query.RoleID)
	}
	if query.ID != "" {
		ql = ql.Where("id = ?", query.ID)
	}
	if query.Path != "" {
		ql = ql.Where("path = ?", query.Path)
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
	if permit.Condition != nil {
		setMap["condition"] = permit.Condition
	}
	return db.Table(t.TableName()).Where("id = ? ", id).Updates(
		setMap).Error
}

func (t *permitRepo) List(db *gorm.DB, query *models.PermitQuery, page, size int) ([]*models.Permit, int64, error) {
	page, size = pages(page, size)
	db = db.Table(t.TableName())
	if query.RoleID != "" {
		db = db.Where("role_id = ?", query.RoleID)
	}

	if query.Path != "" {
		db = db.Where("path =  ?", query.Path)
	}
	if len(query.Paths) != 0 {
		db = db.Where("path in ?", query.Paths)
	}

	if len(query.RoleIDs) != 0 {
		db = db.Where("role_id in ?", query.RoleIDs)
	}

	var (
		count   int64
		permits []*models.Permit
	)

	err := db.Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	err = db.Offset((page - 1) * size).Limit(size).Find(&permits).Error
	if err != nil {
		return nil, 0, err
	}

	return permits, count, nil
}

func (t *permitRepo) TableName() string {
	return "permit"
}

func NewPermitRepo() models.PermitRepo {
	return &permitRepo{}
}
