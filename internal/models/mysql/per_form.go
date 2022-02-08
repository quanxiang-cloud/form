package mysql

import (
	"github.com/quanxiang-cloud/form/internal/models"
	"gorm.io/gorm"
)

type perFormRepo struct{}

func (p *perFormRepo) BatchCreate(db *gorm.DB, permitForm ...*models.PermitForm) error {
	return db.Table(p.TableName()).CreateInBatches(permitForm, len(permitForm)).Error
}

func (p *perFormRepo) Get(db *gorm.DB, permitID, formID string) (*models.PermitForm, error) {
	permitForm := new(models.PermitForm)
	err := db.Table(p.TableName()).Where("permit_id = ? and  form_id = ? ", permitID, formID).Find(permitForm).Error
	if err != nil {
		return nil, err
	}
	return permitForm, nil
}

func (p *perFormRepo) Find(db *gorm.DB, query *models.PerFormQuery) ([]*models.PermitForm, error) {
	ql := db.Table(p.TableName())
	if query.FormID != "" {
		ql = ql.Where("form_id = ?", query.FormID)
	}
	if query.PerGroupID != "" {
		ql = ql.Where("permit_id = ?", query.PerGroupID)
	}
	perFormList := make([]*models.PermitForm, 0)
	err := ql.Find(&perFormList).Error
	if err != nil {
		return nil, err
	}
	return perFormList, nil
}

func (p *perFormRepo) Delete(db *gorm.DB, query *models.PerFormQuery) error {
	panic("implement me")
}

func (p *perFormRepo) Update(db *gorm.DB, permitID, formID string, form *models.PermitForm) error {
	setMap := make(map[string]interface{})
	if form.Authority != 0 {
		setMap["authority"] = form.Authority
	}
	if form.FieldJSON != nil {
		setMap["field_json"] = form.FieldJSON
	}
	if &form.WebSchema != nil {
		setMap["web_schema"] = form.WebSchema
	}
	if form.Conditions != nil {
		setMap["conditions"] = form.Conditions
	}
	return db.Table(p.TableName()).Where("permit_id = ? and form_id = ? ", permitID, formID).Updates(
		setMap).Error
}

func (p *perFormRepo) TableName() string {
	return "permit_form"
}

func NewGroupFormRepo() models.GroupFormRepo {
	return &perFormRepo{}
}
