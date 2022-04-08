package mysql

import (
	"github.com/quanxiang-cloud/form/internal/models"
	"gorm.io/gorm"
)

type roleUserRepo struct{}

func (r *roleUserRepo) Get(db *gorm.DB, appID, userID string) (*models.UserRole, error) {
	userRole := new(models.UserRole)
	err := db.Table(r.TableName()).Where("user_id = ? and  app_id = ? ", userID, appID).Find(userRole).Error
	if err != nil {
		return nil, err
	}
	return userRole, nil
}

func (r *roleUserRepo) Delete(db *gorm.DB, query *models.UserRoleQuery) error {
	resp := make([]models.Role, 0)
	ql := db.Table(r.TableName())
	if query.AppID != "" {
		ql = ql.Where("app_id = ?", query.AppID)
	}
	if query.UserID != "" {
		ql = ql.Where("user_id = ?", query.UserID)
	}
	return ql.Delete(resp).Error
}

func (r *roleUserRepo) List(db *gorm.DB, query *models.UserRoleQuery, page, size int) ([]*models.UserRole, int64, error) {
	panic("implement me")
}

func (r *roleUserRepo) TableName() string {
	return "table_user"
}

func (r *roleUserRepo) BatchCreate(db *gorm.DB, userRole ...*models.UserRole) error {
	return db.Table(r.TableName()).CreateInBatches(userRole, len(userRole)).Error
}

func NewUserRoleRepo() models.UserRoleRepo {
	return &roleUserRepo{}
}
