package service

import (
	"context"
	"github.com/quanxiang-cloud/cabin/tailormade/header"

	daprd "github.com/dapr/go-sdk/client"
	id2 "github.com/quanxiang-cloud/cabin/id"
	"github.com/quanxiang-cloud/cabin/logger"
	redis2 "github.com/quanxiang-cloud/cabin/tailormade/db/redis"
	time2 "github.com/quanxiang-cloud/cabin/time"
	"github.com/quanxiang-cloud/form/internal/component/event"
	"github.com/quanxiang-cloud/form/internal/models"
	"github.com/quanxiang-cloud/form/internal/models/mysql"
	"github.com/quanxiang-cloud/form/internal/models/redis"
	config2 "github.com/quanxiang-cloud/form/pkg/misc/config"
	"gorm.io/gorm"
)

const (
	form_permit = "form-permit"
)

type Permit interface {
	CreateRole(ctx context.Context, req *CreateRoleReq) (*CreateRoleResp, error)

	UpdateRole(ctx context.Context, req *UpdateRoleReq) (*UpdateRoleResp, error)

	DeleteRole(ctx context.Context, req *DeleteRoleReq) (*DeleteRoleResp, error) // 这个删除需要关心的东西比较多

	GetRole(ctx context.Context, req *GetRoleReq) (*GetRoleResp, error)

	FindRole(ctx context.Context, req *FindRoleReq) (*FindRoleResp, error)

	AssignRoleGrant(ctx context.Context, req *AssignRoleGrantReq) (*AssignRoleGrantResp, error)

	FindGrantRole(ctx context.Context, req *FindGrantRoleReq) (*FindGrantRoleResp, error)

	CreatePermit(ctx context.Context, req *CreatePerReq) (*CreatePerResp, error)

	UpdatePermit(ctx context.Context, req *UpdatePerReq) (*UpdatePerResp, error)

	DeletePermit(ctx context.Context, req *DeletePerReq) (*DeletePerResp, error)

	GetPermit(ctx context.Context, req *GetPermitReq) (*GetPermitResp, error)

	FindPermit(ctx context.Context, req *FindPermitReq) (*FindPermitResp, error)

	SaveUserPerMatch(ctx context.Context, req *SaveUserPerMatchReq) (*SaveUserPerMatchResp, error)

	ListPermit(ctx context.Context, req *ListPermitReq) (*ListPermitResp, error)

	ListAndSelect(ctx context.Context, req *ListAndSelectReq) (*ListAndSelectResp, error)
}

type permit struct {
	db            *gorm.DB
	roleRepo      models.RoleRepo
	roleGrantRepo models.RoleRantRepo
	permitRepo    models.PermitRepo
	limitRepo     models.LimitsRepo
	daprClient    daprd.Client
	conf          *config2.Config
}

type ListAndSelectReq struct {
	AppID  string `json:"appID"`
	UserID string `json:"userID"`
	DepID  string `json:"depID"`
}

type ListAndSelectResp struct {
	OptionPer []*Per `json:"optionPer"`
	SelectPer *Per   `json:"selectPer"`
}

type Per struct {
	RoleID   string `json:"roleID"`
	RoleName string `json:"roleName"`
}

func (p *permit) ListAndSelect(ctx context.Context, req *ListAndSelectReq) (*ListAndSelectResp, error) {
	//

	list, _, err := p.roleGrantRepo.List(p.db, &models.RoleGrantQuery{
		Owners: []string{req.DepID, req.UserID},
		AppID:  req.AppID,
	}, 1, 999)
	if err != nil {
		return nil, err
	}
	ids := make([]string, len(list))
	for index, value := range list {
		ids[index] = value.RoleID
	}
	roles, _, err := p.roleRepo.List(p.db, &models.RoleQuery{
		RoleIDS: ids,
	}, 1, 999)
	if err != nil {
		return nil, err
	}
	resp := &ListAndSelectResp{
		OptionPer: make([]*Per, len(roles)),
	}
	for index, value := range roles {
		resp.OptionPer[index] = &Per{
			RoleID:   value.ID,
			RoleName: value.Name,
		}
	}
	// TODO

	return resp, nil
}

type ListPermitReq struct {
	RoleID string   `json:"roleID,omitempty"`
	Paths  []string `json:"paths"`
	URIs   []string `json:"uris"`
}
type ListPermitResp map[string]*ListVo

type ListVo struct {
	Params    models.FiledPermit `json:"params"`
	Response  models.FiledPermit `json:"response"`
	Condition models.Condition   `json:"condition"`
}

func (p *permit) ListPermit(ctx context.Context, req *ListPermitReq) (*ListPermitResp, error) {
	if len(req.Paths) == 0 || len(req.URIs) == 0 {
		return &ListPermitResp{}, nil
	}
	if len(req.Paths) > 100 {
		req.Paths = req.Paths[0:100]
		req.URIs = req.URIs[0:100]
	}
	if len(req.URIs) != len(req.Paths) {
		return nil, nil
	}
	temp := make(map[string]string)

	for index, value := range req.URIs {
		temp[value] = req.Paths[index]
	}
	form := IsFormAPI(req.Paths[0])
	if form {
		req.Paths = req.URIs
	}
	permits, _, err := p.permitRepo.List(p.db, &models.PermitQuery{
		RoleID: req.RoleID,
		Paths:  req.Paths,
	}, 1, 100)
	if err != nil {
		return nil, err
	}
	resp := make(ListPermitResp)
	for _, value := range permits {
		key := value.Path
		if form {
			key = temp[value.Path]
		}
		resp[key] = &ListVo{
			Params:    value.Params,
			Response:  value.Response,
			Condition: value.Condition,
		}
	}
	return &resp, nil
}

type FindPermitReq struct {
	RoleID string `json:"roleID"`
	Page   int    `json:"page"`
	Size   int    `json:"size"`
}

type FindPermitResp struct {
	List  []*Permits `json:"list"`
	Total int64      `json:"total"`
}

type Permits struct {
	ID        string             `json:"id"`
	RoleID    string             `json:"roleID"`
	Path      string             `json:"path"`
	Params    models.FiledPermit `json:"params"`
	Response  models.FiledPermit `json:"response"`
	Condition models.Condition   `json:"condition"`
}

func (p *permit) FindPermit(ctx context.Context, req *FindPermitReq) (*FindPermitResp, error) {
	permits, total, err := p.permitRepo.List(p.db, &models.PermitQuery{
		RoleID: req.RoleID,
	}, req.Page, req.Size)
	if err != nil {
		return nil, err
	}
	resp := &FindPermitResp{
		List:  make([]*Permits, len(permits)),
		Total: total,
	}
	for index, value := range permits {
		resp.List[index] = &Permits{
			ID:        value.ID,
			RoleID:    value.RoleID,
			Path:      value.Path,
			Params:    value.Params,
			Condition: value.Condition,
		}
	}
	return resp, nil
}

type FindGrantRoleReq struct {
	Owners []string `json:"owners"`
	AppID  string   `json:"appID"`
	RoleID string   `json:"roleID"`
	Page   int      `json:"page"`
	Size   int      `json:"size"`
}

type FindGrantRoleResp struct {
	List  []*GrantRoles `json:"list"`
	Total int64         `json:"total"`
}

type GrantRoles struct {
	RoleID    string `json:"roleID"`
	Owner     string `json:"id"`
	OwnerName string `json:"name"`
	Types     int    `json:"type"`
}

func (p *permit) FindGrantRole(ctx context.Context, req *FindGrantRoleReq) (*FindGrantRoleResp, error) {
	grantRole, total, err := p.roleGrantRepo.List(p.db, &models.RoleGrantQuery{
		Owners: req.Owners,
		AppID:  req.AppID,
		RoleID: req.RoleID,
	}, req.Page, req.Size)
	if err != nil {
		return nil, err
	}
	resp := &FindGrantRoleResp{
		List:  make([]*GrantRoles, 0, len(grantRole)),
		Total: total,
	}
	for _, value := range grantRole {
		resp.List = append(resp.List, &GrantRoles{
			RoleID:    value.RoleID,
			Owner:     value.Owner,
			OwnerName: value.OwnerName,
			Types:     value.Types,
		})
	}
	return resp, nil
}

type SaveUserPerMatchReq struct {
	RoleID   string `json:"roleID"`
	RoleName string `json:"roleName"`
	UserID   string `json:"userID"`
	AppID    string `json:"appID"`
}

type SaveUserPerMatchResp struct{}

func (p *permit) SaveUserPerMatch(ctx context.Context, req *SaveUserPerMatchReq) (*SaveUserPerMatchResp, error) {
	userSpec := event.Data{
		UserSpec: &event.UserSpec{
			RoleID: req.RoleID,
			UserID: req.UserID,
			AppID:  req.AppID,
			Action: "create",
		},
	}
	err := p.publish(ctx, "form-user-match", userSpec)
	if err != nil {
		return nil, err
	}
	return &SaveUserPerMatchResp{}, nil
}

func NewPermit(conf *config2.Config) (Permit, error) {
	db, err := CreateMysqlConn(conf)
	if err != nil {
		return nil, err
	}
	redisClient, err := redis2.NewClient(conf.Redis)
	if err != nil {
		return nil, err
	}
	client, err := daprd.NewClient()
	if err != nil {
		return nil, err
	}
	return &permit{
		db:            db,
		conf:          conf,
		daprClient:    client,
		roleRepo:      mysql.NewRoleRepo(),
		roleGrantRepo: mysql.NewRoleGrantRepo(),
		permitRepo:    mysql.NewPermitRepo(),
		limitRepo:     redis.NewLimitRepo(redisClient),
	}, nil
}

type CreateRoleReq struct {
	UserID      string          `json:"user_id"`
	UserName    string          `json:"user_name"`
	AppID       string          `json:"app_id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Types       models.RoleType `json:"types"`
}

type CreateRoleResp struct {
	ID string `json:"id"`
}

func (p *permit) CreateRole(ctx context.Context, req *CreateRoleReq) (*CreateRoleResp, error) {
	roles := &models.Role{
		ID:          id2.HexUUID(true),
		AppID:       req.AppID,
		Name:        req.Name,
		Description: req.Description,
		CreatedAt:   time2.NowUnix(),
		CreatorName: req.UserName,
		CreatorID:   req.UserID,
	}
	if req.Types == 0 {
		roles.Types = models.CreateType
	}
	roles.Types = req.Types
	err := p.roleRepo.BatchCreate(p.db, roles)
	if err != nil {
		return nil, err
	}
	return &CreateRoleResp{
		ID: roles.ID,
	}, nil
}

type UpdateRoleReq struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type UpdateRoleResp struct{}

func (p *permit) UpdateRole(ctx context.Context, req *UpdateRoleReq) (*UpdateRoleResp, error) {
	err := p.roleRepo.Update(p.db, req.ID, &models.Role{
		Name:        req.Name,
		Description: req.Description,
	})
	if err != nil {
		return nil, err
	}
	return &UpdateRoleResp{}, nil
}

type GetRoleReq struct {
	ID string `json:"id"`
}
type GetRoleResp struct {
	Types       models.RoleType `json:"type"`
	ID          string          `json:"id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
}

func (p *permit) GetRole(ctx context.Context, req *GetRoleReq) (*GetRoleResp, error) {
	permits, err := p.roleRepo.Get(p.db, req.ID)
	if err != nil {
		return nil, err
	}
	return &GetRoleResp{
		ID:          permits.ID,
		Types:       permits.Types,
		Name:        permits.Name,
		Description: permits.Description,
	}, nil
}

type FindRoleReq struct {
	AppID string `json:"appID"`
	Page  int    `json:"page"`
	Size  int    `json:"size"`
}

type FindRoleResp struct {
	List  []*roleVo `json:"list"`
	Total int64     `json:"total"`
}

type roleVo struct {
	Types       models.RoleType `json:"type"`
	ID          string          `json:"id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
}

func (p *permit) FindRole(ctx context.Context, req *FindRoleReq) (*FindRoleResp, error) {
	list, total, err := p.roleRepo.List(p.db, &models.RoleQuery{
		AppID: req.AppID,
	}, req.Page, req.Size)
	if err != nil {
		return nil, err
	}
	resp := &FindRoleResp{
		Total: total,
		List:  make([]*roleVo, len(list)),
	}
	for index, value := range list {
		resp.List[index] = &roleVo{
			ID:          value.ID,
			Name:        value.Name,
			Types:       value.Types,
			Description: value.Description,
		}
	}
	return resp, nil
}

type AssignRoleGrantReq struct {
	Add     []*Owners `json:"add"`
	RoleID  string    `json:"roleID"`
	AppID   string    `json:"appID"`
	Removes []string  `json:"removes"`
}
type Owners struct {
	Owner     string `json:"id"`
	OwnerName string `json:"name"`
	Types     int    `json:"type"`
}

type AssignRoleGrantResp struct{}

func (p *permit) AssignRoleGrant(ctx context.Context, req *AssignRoleGrantReq) (*AssignRoleGrantResp, error) {
	roleGrants := make([]*models.RoleGrant, len(req.Add))
	for index, value := range req.Add {
		roleGrants[index] = &models.RoleGrant{
			ID:        id2.HexUUID(true),
			RoleID:    req.RoleID,
			Owner:     value.Owner,
			OwnerName: value.OwnerName,
			Types:     value.Types,
			AppID:     req.AppID,
			CreatedAt: time2.NowUnix() + int64(index),
		}
	}
	err := p.roleGrantRepo.BatchCreate(p.db, roleGrants...)
	if err != nil {
		return nil, err
	}

	if len(req.Removes) == 0 {
		return &AssignRoleGrantResp{}, nil
	}
	err = p.roleGrantRepo.Delete(p.db, &models.RoleGrantQuery{
		RoleID: req.RoleID,
		Owners: req.Removes,
	})
	if err != nil {
		return nil, err
	}
	return &AssignRoleGrantResp{}, nil
}

type CreatePerReq struct {
	AccessPath string             `json:"path"`
	URI        string             `json:"uri"`
	Params     models.FiledPermit `json:"params"`
	Response   models.FiledPermit `json:"response"`
	RoleID     string             `json:"roleID"`
	UserID     string             `json:"userID"`
	UserName   string             `json:"userName"`
	Condition  models.Condition   `json:"condition"`
}

type CreatePerResp struct{}

func (p *permit) CreatePermit(ctx context.Context, req *CreatePerReq) (*CreatePerResp, error) {
	if IsFormAPI(req.AccessPath) {
		req.AccessPath = req.URI
	}
	permits := &models.Permit{
		ID:          id2.HexUUID(true),
		Path:        req.AccessPath,
		Params:      req.Params,
		Response:    req.Response,
		RoleID:      req.RoleID,
		CreatorID:   req.UserID,
		CreatorName: req.UserName,
		CreatedAt:   time2.NowUnix(),
		Condition:   req.Condition,
	}
	err := p.permitRepo.BatchCreate(p.db, permits)
	if err != nil {
		return nil, err
	}
	spec := &event.PermitSpec{
		RoleID:    req.RoleID,
		Path:      req.AccessPath,
		Condition: req.Condition,
		Response:  req.Response,
		Params:    req.Params,
		Action:    "create",
	}
	err = p.publish(ctx, form_permit, &event.Data{
		PermitSpec: spec,
	})
	if err != nil {
		logger.Logger.WithName("publish permit create ").Errorw("publish", "topic", form_permit, header.GetRequestIDKV(ctx).Fuzzy(), err.Error())
	}
	return &CreatePerResp{}, nil
}

type UpdatePerReq struct {
	ID        string             `json:"id"`
	Params    models.FiledPermit `json:"params"`
	Response  models.FiledPermit `json:"response"`
	Condition models.Condition   `json:"condition"`
}

type UpdatePerResp struct{}

func (p *permit) UpdatePermit(ctx context.Context, req *UpdatePerReq) (*UpdatePerResp, error) {
	err := p.permitRepo.Update(p.db, req.ID, &models.Permit{
		Params:   req.Params,
		Response: req.Response,
	})
	if err != nil {
		return nil, err
	}
	// add redis cache
	return &UpdatePerResp{}, nil
}

type DeletePerReq struct {
	RoleID string `json:"roleID"`
	Path   string `json:"path"`
	URI    string `json:"uri"`
}

type DeletePerResp struct{}

func (p *permit) DeletePermit(ctx context.Context, req *DeletePerReq) (*DeletePerResp, error) {
	if IsFormAPI(req.Path) {
		req.Path = req.URI
	}
	err := p.permitRepo.Delete(p.db, &models.PermitQuery{
		RoleID: req.RoleID,
		Path:   req.Path,
	})

	if err != nil {
		return nil, err
	}
	spec := &event.PermitSpec{
		RoleID: req.RoleID,
		Path:   req.Path,
		Action: "delete",
	}

	// TODO dapr
	err = p.publish(ctx, form_permit, &event.Data{
		PermitSpec: spec,
	})
	if err != nil {
		logger.Logger.WithName("publish permit create ").Errorw("publish", "topic", form_permit, "err is", err.Error())
	}
	return &DeletePerResp{}, nil
}

type GetPermitReq struct {
	RoleID string `json:"roleID"`
	Path   string `json:"path"`
	URI    string `json:"uri"`
}

type GetPermitResp struct {
	ID        string             `json:"id"`
	RoleID    string             `json:"roleID"`
	Path      string             `json:"path,omitempty"`
	Params    models.FiledPermit `json:"params,omitempty"`
	Response  models.FiledPermit `json:"response,omitempty"`
	Condition models.Condition   `json:"condition,omitempty"`
}

func (p *permit) GetPermit(ctx context.Context, req *GetPermitReq) (*GetPermitResp, error) {
	if IsFormAPI(req.Path) {
		req.Path = req.URI
	}
	permits, err := p.permitRepo.Get(p.db, req.RoleID, req.Path)
	if err != nil {
		return nil, err
	}
	return &GetPermitResp{
		ID:        permits.ID,
		RoleID:    permits.RoleID,
		Path:      permits.Path,
		Params:    permits.Params,
		Response:  permits.Response,
		Condition: permits.Condition,
	}, nil
}

type DeleteRoleReq struct {
	RoleID string `json:"-"`
	AppID  string `json:"-"`
}
type DeleteRoleResp struct {
}

func (p *permit) DeleteRole(ctx context.Context, req *DeleteRoleReq) (*DeleteRoleResp, error) {
	err := p.roleRepo.Delete(p.db, &models.RoleQuery{
		ID: req.RoleID,
	})
	if err != nil {
		return nil, err
	}
	// 删除对应 角色的人
	err = p.roleGrantRepo.Delete(p.db, &models.RoleGrantQuery{
		RoleID: req.RoleID,
	})
	if err != nil {
		return nil, err
	}
	// 删除，role 对应的permit
	err = p.permitRepo.Delete(p.db, &models.PermitQuery{RoleID: req.RoleID})
	if err != nil {
		return nil, err
	}
	// 删除缓存
	err = p.limitRepo.DeletePerMatch(ctx, req.AppID)
	if err != nil {
		//
		logger.Logger.Errorw("delete per match", req.RoleID, err.Error())
	}
	//
	err = p.publish(ctx, "form-user-match", &event.Data{
		UserSpec: &event.UserSpec{
			RoleID: req.RoleID,
			AppID:  req.AppID,
			Action: "delete",
		},
	})
	logger.Logger.Errorw("")

	//err = p.limitRepo.DeletePermit(ctx, req.RoleID)
	//if err != nil {
	//
	//}
	return &DeleteRoleResp{}, nil

}

func (p *permit) publish(ctx context.Context, topic string, data interface{}) error {
	if err := p.daprClient.PublishEvent(ctx, p.conf.PubSubName, topic, data); err != nil {
		logger.Logger.WithName("public").Errorw("publish error", "topic", topic, "publicName", p.conf.PubSubName)
		return err
	}
	return nil
}
