package mysql

import (
	"github.com/quanxiang-cloud/form/internal/models"
	"gorm.io/gorm"
)

type perGrantRepo struct{}

func (p *perGrantRepo) Get(db *gorm.DB, permit, owner string) (*models.PerGrant, error) {
	perGrant := new(models.PerGrant)
	err := db.Table(p.TableName()).Where("permit_id = ? and  owner = ? ", permit, owner).Find(perGrant).Error
	if err != nil {
		return nil, err
	}
	return perGrant, nil

}

func (p *perGrantRepo) Find(db *gorm.DB, query *models.PerGrantQuery) ([]*models.PerGrant, error) {
	ql := db.Table(p.TableName())
	if query.PerGroupID != "" {
		ql = ql.Where("permit_id = ?", query.PerGroupID)
	}
	if len(query.Owners) != 0 {
		ql = ql.Where("owner in  ?", query.Owners)
	}
	if len(query.Permits) != 0 {
		ql = ql.Where("permit_id in ?", query.Permits)
	}
	perFormList := make([]*models.PerGrant, 0)
	err := ql.Find(&perFormList).Error
	if err != nil {
		return nil, err
	}
	return perFormList, nil
}

func (p *perGrantRepo) BatchCreate(db *gorm.DB, perGrant ...*models.PerGrant) error {
	return db.Table(p.TableName()).CreateInBatches(perGrant, len(perGrant)).Error
}

func (p *perGrantRepo) Delete(db *gorm.DB, query *models.PerGrantQuery) error {
	return nil
}

func (p *perGrantRepo) TableName() string {
	return "permit_grant"
}

func NewPerGrantRepo() models.PerGrantRepo {
	return &perGrantRepo{}
}
