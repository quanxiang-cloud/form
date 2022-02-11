package service

import (
	"context"
	id2 "github.com/quanxiang-cloud/cabin/id"
	time2 "github.com/quanxiang-cloud/cabin/time"
	"github.com/quanxiang-cloud/form/internal/filters"
	"github.com/quanxiang-cloud/form/internal/models/mysql"
	"github.com/quanxiang-cloud/form/pkg/client"
	"gorm.io/gorm"
	"time"

	error2 "github.com/quanxiang-cloud/cabin/error"
	redis2 "github.com/quanxiang-cloud/cabin/tailormade/db/redis"

	"github.com/quanxiang-cloud/form/internal/models"
	"github.com/quanxiang-cloud/form/internal/models/redis"
	"github.com/quanxiang-cloud/form/pkg/misc/code"
	"github.com/quanxiang-cloud/form/pkg/misc/config"
)

// Permission permission
type Permission interface {
	CreateGroup(ctx context.Context, req *CreateGroup) (*CreateGroupResp, error)
	UpdateGroup(ctx context.Context, req *UpdateGroupReq) (*UpdateGroupResp, error)
	AddOwnerToGroup(ctx context.Context, req *AddOwnerToGroupReq) (*AddOwnerToGroupResp, error)
	AddOwnerToApp(ctx context.Context, req *AddOwnerToAppReq) (*AddOwnerToAppResp, error)

	DelPerGroup(ctx context.Context, req *DelPerGroupReq) (*DelPerGroupResp, error)
	GetPerGroup(ctx context.Context, req *GetPerGroupReq) (*GetPerGroupResp, error)
	FindPerGroup(ctx context.Context, req *FindPerGroupReq) (*FindPerGroupResp, error)

	SaveForm(ctx context.Context, req *SaveFormReq) (*SaveFormResp, error)
	DeleteForm(ctx context.Context, req *DeleteFormReq) (*DeleteFormResp, error)

	FindForm(ctx context.Context, req *FindFormReq) (*FindFormResp, error)
	GetForm(ctx context.Context, req *GetFormReq) (*GetFormResp, error)

	GetGroupsByMenu(ctx context.Context, req *GetGroupsByMenuReq) (*GetGroupsByMenuResp, error)

	SaveUserPerMatch(ctx context.Context, req *SaveUserPerMatchReq) (*SaveUserPerMatchResp, error)

	GetPerInCache(ctx context.Context, req *GetPerInCacheReq) (*GetPerInCacheResp, error)
}

const (
	lockPermission = "lockPermission"
	lockPerMatch   = "lockPerMatch"
	lockTimeout    = time.Duration(30) * time.Second                   // 30秒
	perTime        = time.Hour * time.Duration(12) * time.Duration(30) // 30天
	timeSleep      = time.Millisecond * 500                            // 0.5 秒
	notAuthority   = -1
)

type permission struct {
	db              *gorm.DB
	perGroupRepo    models.PerGroupRepo
	perFormRepo     models.GroupFormRepo
	permitGrantRepo models.PerGrantRepo
	permissionRepo  models.PermissionRepo
	appClient       client.AppCenter
}

// NewPermission NewPermission
func NewPermission(conf *config.Config) (Permission, error) {
	redisClient, err := redis2.NewClient(conf.Redis)
	if err != nil {
		return nil, err
	}
	db, err := createMysqlConn(conf)
	if err != nil {
		return nil, err
	}

	u := &permission{
		db:              db,
		permissionRepo:  redis.NewPermissionRepo(redisClient),
		perGroupRepo:    mysql.NewPerGroupRepo(),
		perFormRepo:     mysql.NewGroupFormRepo(),
		permitGrantRepo: mysql.NewPerGrantRepo(),
	}
	return u, nil
}

type CreateGroup struct {
	AppID       string         `json:"appID"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Types       models.PerType `json:"types"`
	UserID      string         `json:"userID"`
	CreatorName string         `json:"creatorName"`
}

type CreateGroupResp struct {
	ID string `json:"id"`
}

func (per *permission) CreateGroup(ctx context.Context, req *CreateGroup) (*CreateGroupResp, error) {
	perGroup, err := per.perGroupRepo.Find(per.db, &models.PerGroupQuery{
		Name:  req.Name,
		AppID: req.AppID,
	})
	if err != nil {
		return nil, err
	}
	if len(perGroup) > 0 {
		return nil, error2.New(code.ErrExistGroupNameState)
	}
	perGroups := &models.PerGroup{
		ID:          id2.HexUUID(true),
		Name:        req.Name,
		Description: req.Description,
		AppID:       req.AppID,
		CreatedAt:   time2.NowUnix(),
		CreatorID:   req.UserID,
		CreatorName: req.CreatorName,
		Types:       req.Types,
	}
	if req.Types == 0 {
		perGroups.Types = models.CreateType
	}
	err = per.perGroupRepo.BatchCreate(per.db, perGroups)
	if err != nil {
		return nil, err
	}
	return &CreateGroupResp{ID: perGroups.ID}, nil
}

type UpdateGroupReq struct {
	AppID       string `json:"_"`
	Name        string `json:"name"`
	ID          string `json:"id"`
	Description string `json:"description"`
}

type UpdateGroupResp struct {
}

// UpdateGroup UpdateGroup 修改
func (per *permission) UpdateGroup(ctx context.Context, req *UpdateGroupReq) (*UpdateGroupResp, error) {
	if req.Name != "" { // 如果 是空， 判断重复性
		perGroup, err := per.perGroupRepo.Find(per.db, &models.PerGroupQuery{
			Name:  req.Name,
			AppID: req.AppID,
		})
		if err != nil {
			return nil, err
		}
		if len(perGroup) > 0 {
			return nil, error2.New(code.ErrExistGroupNameState)
		}
	}
	permit := &models.PerGroup{
		Name:        req.Name,
		Description: req.Description,
	}
	err := per.perGroupRepo.Update(per.db, req.ID, permit)
	if err != nil {
		return nil, err
	}
	return &UpdateGroupResp{}, nil
}

type GetPerGroupReq struct {
	ID string `json:"id"`
}

type GetPerGroupResp struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	AppID       string     `json:"appID"`
	Grants      []*grantVO `json:"scopes"`
	Description string     `json:"description"`
}

type grantVO struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Types int    `json:"types"`
}

func (per *permission) GetPerGroup(ctx context.Context, req *GetPerGroupReq) (*GetPerGroupResp, error) {
	perGroup, err := per.perGroupRepo.Get(per.db, req.ID)
	if err != nil {
		return nil, err
	}
	resp := &GetPerGroupResp{
		ID:          perGroup.ID,
		Name:        perGroup.Name,
		Description: perGroup.Description,
	}
	return resp, nil
}

type FindPerGroupReq struct {
	AppID  string `json:"appID"`
	UserID string `json:"userID"`
	DepID  string `json:"depID"`
}

type FindPerGroupResp struct {
	ListVO []*listVO `json:"list"`
}

type listVO struct {
	ID          string         `json:"id"`
	Name        string         `json:"name"`
	CreatedBy   string         `json:"createdBy"`
	Scopes      []*grantVO     `json:"scopes"`
	Description string         `json:"description"`
	AppID       string         `json:"appID"`
	Add         bool           `json:"add"`
	Types       models.PerType `json:"types"`
}

func (per *permission) FindPerGroup(ctx context.Context, req *FindPerGroupReq) (*FindPerGroupResp, error) {
	list, err := per.perGroupRepo.Find(per.db, &models.PerGroupQuery{
		AppID: req.AppID,
	})
	if err != nil {
		return nil, err
	}
	resp := &FindPerGroupResp{
		ListVO: make([]*listVO, len(list)),
	}
	for index, value := range list {
		permit := new(listVO)
		clone(permit, value)
		// 查询人
		resp.ListVO[index] = permit
	}
	return resp, nil
}

func clone(dst *listVO, src *models.PerGroup) {
	dst.ID = src.ID
	dst.Name = src.Name
	dst.Description = src.Description
	dst.AppID = src.AppID
	dst.Types = src.Types
}

type DeleteFormReq struct {
	PerGroupID string `json:"perGroupID"`
	FormID     string `json:"formID"`
}

type DeleteFormResp struct {
}

func (per *permission) DeleteForm(ctx context.Context, req *DeleteFormReq) (*DeleteFormResp, error) {
	err := per.perFormRepo.Delete(per.db, &models.PerFormQuery{
		PerGroupID: req.PerGroupID,
		FormID:     req.FormID,
	})
	if err != nil {
		return nil, err
	}
	return &DeleteFormResp{}, nil
}

type SaveFormReq struct {
	FormID     string                  `json:"formID"`
	FormName   string                  `json:"formName"`
	PerGroupID string                  `json:"perGroupID"`
	Authority  int64                   `json:"authority"`
	Conditions map[string]models.Query `json:"conditions"`
	Schema     models.Schema           `json:"schema"`
}

type SaveFormResp struct {
}

func (per *permission) SaveForm(ctx context.Context, req *SaveFormReq) (*SaveFormResp, error) {
	perForm, err := per.perFormRepo.Get(per.db, req.PerGroupID, req.FormID)
	if err != nil {
		return nil, err
	}
	perForms := &models.PermitForm{
		PerGroupID: req.PerGroupID,
		FormID:     req.FormID,
		FormType:   "111",
		Authority:  req.Authority,
		Conditions: req.Conditions,
		WebSchema:  req.Schema,
		FieldJSON:  filters.DealSchemaToFilterType(req.Schema),
	}
	if perForm.FormID == "" {
		err = per.perFormRepo.BatchCreate(per.db, perForms)
		if err != nil {
			return nil, err
		}
	}
	err = per.perFormRepo.Update(per.db, req.PerGroupID, req.FormID, perForms)
	if err != nil {
		return nil, err
	}
	return &SaveFormResp{}, nil
}

// FindFormReq GetFormReq
type FindFormReq struct {
	PerGroupID string `json:"perGroupID"`
}

// FindFormResp GetFormResp
type FindFormResp struct {
	FormArr []*FormVo `json:"formArr"`
}

// FormVo FormVo
type FormVo struct {
	PerGroupID string
	FormID     string
	Authority  int64
}

//FindForm FindForm
func (per *permission) FindForm(ctx context.Context, req *FindFormReq) (*FindFormResp, error) {
	list, err := per.perFormRepo.Find(per.db, &models.PerFormQuery{PerGroupID: req.PerGroupID})
	if err != nil {
		return nil, err
	}
	resp := &FindFormResp{
		FormArr: make([]*FormVo, len(list)),
	}
	for index, value := range list {
		formVo := &FormVo{
			Authority: value.Authority,
		}
		resp.FormArr[index] = formVo
	}
	return resp, err
}

type GetFormReq struct {
	FormID     string `json:"formID"`
	PerGroupID string `json:"perGroupID"`
}

type GetFormResp struct {
	FormID     string            `json:"formID"`
	PerGroupID string            `json:"perGroupID"`
	Conditions models.Conditions `json:"dataAccess"`
	Opt        int64             `json:"opt"`
	Filter     models.Schema     `json:"filter"`
}

func (per *permission) GetForm(ctx context.Context, req *GetFormReq) (*GetFormResp, error) {
	permitForm, err := per.perFormRepo.Get(per.db, req.PerGroupID, req.FormID)
	if err != nil {
		return nil, err
	}
	resp := &GetFormResp{
		FormID:     permitForm.FormID,
		PerGroupID: permitForm.PerGroupID,
		Conditions: permitForm.Conditions,
		Opt:        permitForm.Authority,
		Filter:     permitForm.WebSchema,
	}
	return resp, nil
}

type GetGroupsByMenuReq struct {
}

type GetGroupsByMenuResp struct {
}

func (per *permission) GetGroupsByMenu(ctx context.Context, req *GetGroupsByMenuReq) (*GetGroupsByMenuResp, error) {
	return nil, nil
}

type GetOperateReq struct {
	FormID string `json:"formID"`
	UserID string `json:"userID"`
	DepID  string `json:"depID"`
	AppID  string `json:"appID"`
}

type GetOperateResp struct {
	Authority int64 `json:"authority"`
}

func (per *permission) GetOperate(ctx context.Context, req *GetOperateReq) (*GetOperateResp, error) {
	match, err := per.getPerMatch(ctx, req.UserID, req.DepID, req.AppID)
	if err != nil {
		return nil, err
	}
	if match == nil {
		return &GetOperateResp{}, nil
	}
	info, err := per.getPerInfo(ctx, match.PerGroupID, req.FormID)
	if err != nil {
		return nil, err
	}
	if info == nil {
		return &GetOperateResp{}, nil
	}
	return &GetOperateResp{
		Authority: info.Authority,
	}, nil
}

type GetPerSelectReq struct {
	UserID string `json:"userID"`
	DepID  string `json:"depID"`
	AppID  string `json:"appID"`
}

type GetPerSelectResp struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (per *permission) GetPerSelect(ctx context.Context, req *GetPerSelectReq) (*GetPerSelectResp, error) {
	match, err := per.getPerMatch(ctx, req.UserID, req.DepID, req.AppID)
	if err != nil {
		return nil, err
	}
	if match == nil {
		return nil, nil
	}
	perGroup, err := per.perGroupRepo.Get(per.db, match.PerGroupID)
	if err != nil {
		return nil, err
	}
	if perGroup != nil {
		return &GetPerSelectResp{
			ID:   match.PerGroupID,
			Name: perGroup.Name,
		}, nil
	}
	return nil, nil
}
