package mysql

import (
	"github.com/quanxiang-cloud/form/internal/models"
	"gorm.io/gorm"
)

type roleGrantRepo struct{}

func (t *roleGrantRepo) List(db *gorm.DB, query *models.RoleGrantQuery, page, size int) ([]*models.RoleGrant, int64, error) {
	page, size = pages(page, size)
	var (
		count          int64
		tableRelations []*models.RoleGrant
	)
	ql := db.Table(t.TableName())
	if len(query.Owners) != 0 {
		ql = ql.Where("owner in ? ", query.Owners)
	}
	if query.AppID != "" {
		ql = ql.Where("app_id = ?", query.AppID)
	}
	if query.RoleID != "" {
		ql = ql.Where("role_id = ?", query.RoleID)
	}
	if len(query.RoleIDs) != 0 {
		ql = ql.Where("role_id in ?", query.RoleIDs)
	}
	ql = ql.Order("created_at DESC")

	err := db.Count(&count).Error
	if err != nil {
		return nil, 0, err
	}
	err = db.Order("created_at desc").Offset((page - 1) * size).Limit(size).Find(&tableRelations).Error
	if err != nil {
		return nil, 0, err
	}

	return tableRelations, count, nil
}

func (t *roleGrantRepo) BatchCreate(db *gorm.DB, roleGrant ...*models.RoleGrant) error {
	return db.Table(t.TableName()).CreateInBatches(roleGrant, len(roleGrant)).Error
}
func (t *roleGrantRepo) Get(db *gorm.DB, id string) (*models.RoleGrant, error) {
	roleGrant := new(models.RoleGrant)
	err := db.Table(t.TableName()).Where("id = ? ", id).Find(roleGrant).Error
	if err != nil {
		return nil, err
	}
	return roleGrant, nil
}

func (t *roleGrantRepo) Find(db *gorm.DB, query *models.RoleGrantQuery) ([]*models.RoleGrant, error) {
	ql := db.Table(t.TableName())
	if len(query.Owners) != 0 {
		ql = ql.Where("owner in ? ", query.Owners)
	}
	if query.AppID != "" {
		ql = ql.Where("app_id = ?", query.AppID)
	}
	if query.RoleID != "" {
		ql = ql.Where("role_id = ?", query.RoleID)
	}
	if len(query.RoleIDs) != 0 {
		ql = ql.Where("role_id in ?", query.RoleIDs)
	}
	ql = ql.Order("created_at DESC")
	roleGrant := make([]*models.RoleGrant, 0)
	err := ql.Find(&roleGrant).Error
	return roleGrant, err

}

func (t *roleGrantRepo) Update(db *gorm.DB, id string, roleGrant *models.RoleGrant) error {
	return nil
}

func (t *roleGrantRepo) Delete(db *gorm.DB, query *models.RoleGrantQuery) error {
	resp := make([]models.RoleGrant, 0)
	ql := db.Table(t.TableName())
	if query.RoleID != "" {
		ql = ql.Where("role_id = ? ", query.RoleID)
	}
	if len(query.Owners) != 0 {
		ql = ql.Where("owner in ? ", query.Owners)
	}
	return ql.Delete(resp).Error
}

func (t *roleGrantRepo) TableName() string {
	return "role_grant"
}

func NewRoleGrantRepo() models.RoleRantRepo {
	return &roleGrantRepo{}
}
