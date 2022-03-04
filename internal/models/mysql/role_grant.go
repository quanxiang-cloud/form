package mysql

import (
	"github.com/quanxiang-cloud/form/internal/models"
	"gorm.io/gorm"
)

type roleGrantRepo struct{}

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
