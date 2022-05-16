package service

import (
	"context"
	"fmt"
	daprd "github.com/dapr/go-sdk/client"
	error2 "github.com/quanxiang-cloud/cabin/error"
	id2 "github.com/quanxiang-cloud/cabin/id"
	"github.com/quanxiang-cloud/cabin/logger"
	redis2 "github.com/quanxiang-cloud/cabin/tailormade/db/redis"
	time2 "github.com/quanxiang-cloud/cabin/time"
	"github.com/quanxiang-cloud/form/internal/component/event"
	"github.com/quanxiang-cloud/form/internal/models"
	"github.com/quanxiang-cloud/form/internal/models/mysql"
	"github.com/quanxiang-cloud/form/internal/models/redis"
	"github.com/quanxiang-cloud/form/pkg/misc/code"
	config2 "github.com/quanxiang-cloud/form/pkg/misc/config"
	"gorm.io/gorm"
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

	CreateUserRole(ctx context.Context, req *CreateUserRoleReq, opts ...Option) (*CreateUserRoleResp, error)

	ListPermit(ctx context.Context, req *ListPermitReq) (*ListPermitResp, error)

	ListAndSelect(ctx context.Context, req *ListAndSelectReq) (*ListAndSelectResp, error)

	GetUserRole(ctx context.Context, req *GetUserRoleReq) (*GetUserRoleResp, error)

	CopyRole(ctx context.Context, req *CopyRoleReq) (*CopyRoleResp, error)
}

type permit struct {
	db            *gorm.DB
	roleRepo      models.RoleRepo
	roleGrantRepo models.RoleRantRepo
	permitRepo    models.PermitRepo
	limitRepo     models.LimitsRepo
	daprClient    daprd.Client
	conf          *config2.Config
	userRoleRepo  models.UserRoleRepo
}

type CopyRoleReq struct {
	RoleID      string `json:"roleID"`
	UserID      string `json:"userID"`
	AppID       string `json:"appID"`
	UserName    string `json:"userName"`
	Description string `json:"description"`
	Name        string `json:"name"`
}

type CopyRoleResp struct {
	RoleID string `json:"id"`
}

func (p *permit) CopyRole(ctx context.Context, req *CopyRoleReq) (*CopyRoleResp, error) {
	tx := p.db.Begin()
	roleID := id2.StringUUID()
	err := p.roleRepo.BatchCreate(tx, &models.Role{
		ID:          roleID,
		Description: req.Description,
		Name:        req.Name,
		AppID:       req.AppID,
		CreatorName: req.UserName,
		CreatorID:   req.UserID,
		Types:       models.CreateType,
		CreatedAt:   time2.NowUnix(),
	})
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	list, _, err := p.permitRepo.List(tx, &models.PermitQuery{
		RoleID: req.RoleID,
	}, 1, 999)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	lists := make([]*models.Permit, len(list))
	for index, value := range list {
		lists[index] = &models.Permit{
			ID:          id2.StringUUID(),
			RoleID:      roleID,
			Path:        value.Path,
			Params:      value.Params,
			Condition:   value.Condition,
			Method:      value.Method,
			ParamsAll:   value.ParamsAll,
			ResponseAll: value.ResponseAll,
			CreatedAt:   time2.NowUnix(),
			CreatorID:   req.UserID,
			CreatorName: req.UserName,
		}
	}
	err = p.permitRepo.BatchCreate(tx, lists...)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()
	return &CopyRoleResp{
		RoleID: roleID,
	}, nil
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
	ow := make([]string, 0)
	if req.UserID != "" {
		ow = append(ow, req.UserID)
	}
	if req.DepID != "" {
		ow = append(ow, req.DepID)
	}
	list, _, err := p.roleGrantRepo.List(p.db, &models.RoleGrantQuery{
		Owners: ow,
		AppID:  req.AppID,
	}, 1, 999)
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return &ListAndSelectResp{
			OptionPer: make([]*Per, 0),
			SelectPer: &Per{},
		}, nil
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
	//
	userRole, err := p.userRoleRepo.Get(p.db, req.AppID, req.UserID)
	if err != nil {
		return nil, err
	}
	get, err := p.roleRepo.Get(p.db, userRole.RoleID)
	if err != nil {
		return nil, err
	}
	resp.SelectPer = &Per{
		RoleID:   userRole.RoleID,
		RoleName: get.Name,
	}
	return resp, nil
}

type ListPermitReq struct {
	RoleID string     `json:"roleID"`
	List   []*ListRes `json:"paths" binding:"required"`
}

type ListRes struct {
	URI        string `json:"uri"`
	AccessPath string `json:"accessPath"`
	Method     string `json:"method"`
}

type ListPermitResp map[string]bool

func (p *permit) ListPermit(ctx context.Context, req *ListPermitReq) (*ListPermitResp, error) {
	resp := make(ListPermitResp)
	if req.RoleID == "" {
		return &resp, nil
	}
	for _, values := range req.List {
		url := values.AccessPath
		if IsFormAPI(values.AccessPath) {
			url = values.URI
		}
		per, err := p.permitRepo.Get(p.db, req.RoleID, url, values.Method)
		if err != nil {
			continue
		}
		if per.ID != "" {
			key := fmt.Sprintf("%s-%s", values.AccessPath, values.Method)
			resp[key] = true
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
	Methods   string             `json:"methods"`
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
	Types  int      `json:"type"`
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
		Types:  req.Types,
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

type CreateUserRoleReq struct {
	RoleID string `json:"roleID"`
	UserID string `json:"userID"`
	AppID  string `json:"appID"`
}

type CreateUserRoleResp struct{}

func (p *permit) CreateUserRole(ctx context.Context, req *CreateUserRoleReq, opts ...Option) (resp *CreateUserRoleResp, err error) {
	// TO 删除

	defer func() {
		userSpec := &OptionReq{
			data: event.Data{
				UserSpec: &event.UserSpec{
					RoleID: req.RoleID,
					UserID: req.UserID,
					AppID:  req.AppID,
					Action: "create",
				},
			},
		}
		if err == nil {
			for _, opt := range opts {
				opt(ctx, userSpec)
			}
		}
	}()

	resp = &CreateUserRoleResp{}
	err = p.userRoleRepo.Delete(p.db, &models.UserRoleQuery{
		UserID: req.UserID,
		AppID:  req.AppID,
	})
	if err != nil {
		return
	}
	err = p.userRoleRepo.BatchCreate(p.db, &models.UserRole{
		UserID: req.UserID,
		RoleID: req.RoleID,
		AppID:  req.AppID,
		ID:     id2.StringUUID(),
	})
	if err != nil {
		return
	}

	return
}

type OptionReq struct {
	data event.Data
}

//Option Option
type Option func(ctx context.Context, req *OptionReq)

func RoleUserOption(permit2 Permit) Option {
	return func(ctx context.Context, req *OptionReq) {
		k2, ok := permit2.(*permit)
		if !ok {
			return
		}
		err := k2.publish(ctx, "form-user-match", req.data)
		if err != nil {
			logger.Logger.Errorw("", "xxxxx")
			return
		}

	}
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
	return &permit{
		db:   db,
		conf: conf,

		roleRepo:      mysql.NewRoleRepo(),
		roleGrantRepo: mysql.NewRoleGrantRepo(),
		permitRepo:    mysql.NewPermitRepo(),
		userRoleRepo:  mysql.NewUserRoleRepo(),
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
	_, total, err := p.roleRepo.List(p.db, &models.RoleQuery{
		Name:  req.Name,
		AppID: req.AppID,
	}, 1, 999)
	if err != nil {
		return nil, err
	}
	if total > 0 {
		return nil, error2.New(code.ErrExistRoleNameState)
	}
	roles := &models.Role{
		ID:          id2.StringUUID(),
		AppID:       req.AppID,
		Name:        req.Name,
		Description: req.Description,
		CreatedAt:   time2.NowUnix(),
		CreatorName: req.UserName,
		CreatorID:   req.UserID,
	}
	roles.Types = req.Types
	if req.Types == 0 {
		roles.Types = models.CreateType
	}
	err = p.roleRepo.BatchCreate(p.db, roles)
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
	roles, err := p.roleRepo.Get(p.db, req.ID)
	if err != nil {
		return nil, err
	}
	return &GetRoleResp{
		ID:          roles.ID,
		Types:       roles.Types,
		Name:        roles.Name,
		Description: roles.Description,
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
	tx := p.db.Begin()
	roleGrants := make([]*models.RoleGrant, len(req.Add))
	for index, value := range req.Add {
		roleGrants[index] = &models.RoleGrant{
			ID:        id2.StringUUID(),
			RoleID:    req.RoleID,
			Owner:     value.Owner,
			OwnerName: value.OwnerName,
			Types:     value.Types,
			AppID:     req.AppID,
			CreatedAt: time2.NowUnix() + int64(index),
		}
	}
	err := p.roleGrantRepo.BatchCreate(tx, roleGrants...)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	if len(req.Removes) == 0 {
		tx.Commit()
		return &AssignRoleGrantResp{}, nil
	}
	err = p.roleGrantRepo.Delete(p.db, &models.RoleGrantQuery{
		RoleID: req.RoleID,
		Owners: req.Removes,
	})
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	// 删除关联关系
	err = p.userRoleRepo.Delete(p.db, &models.UserRoleQuery{
		UserIDS: req.Removes,
		AppID:   req.AppID,
		RoleID:  req.RoleID,
	})
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()
	return &AssignRoleGrantResp{}, nil
}

type CreatePerReq struct {
	AccessPath string             `json:"path"`
	URI        string             `json:"uri"`
	Params     models.FiledPermit `json:"params"`
	Response   models.FiledPermit `json:"response"`
	Condition  models.Condition   `json:"condition"`
	RoleID     string             `json:"roleID"`
	UserID     string             `json:"userID"`
	UserName   string             `json:"userName"`
	Method     string             `json:"method"`
}

type CreatePerResp struct{}

// CreatePermit CreatePermit 如果是表单， (uri , post)    (accessPath , post ,get)
func (p *permit) CreatePermit(ctx context.Context, req *CreatePerReq) (*CreatePerResp, error) {

	exist, err := p.checkExist(ctx, req)
	if err != nil {
		return nil, err
	}
	if exist {
		return nil, error2.New(code.ErrExistPermitState)
	}
	permitArr := make([]*models.Permit, 0)
	if IsFormAPI(req.AccessPath) { // is form api
		permitArr = append(permitArr, &models.Permit{
			ID:          id2.StringUUID(),
			Path:        req.URI,
			Params:      req.Params,
			Response:    req.Response,
			RoleID:      req.RoleID,
			CreatorID:   req.UserID,
			CreatorName: req.UserName,
			CreatedAt:   time2.NowUnix(),
			Condition:   req.Condition,
			Method:      req.Method,
			ParamsAll:   true,
			ResponseAll: true,
		})
	}
	permits := &models.Permit{
		ID:          id2.StringUUID(),
		Path:        req.AccessPath,
		Params:      req.Params,
		Response:    req.Response,
		RoleID:      req.RoleID,
		CreatorID:   req.UserID,
		CreatorName: req.UserName,
		CreatedAt:   time2.NowUnix(),
		Condition:   req.Condition,
		ParamsAll:   true,
		ResponseAll: true,
		Method:      req.Method,
	}
	permitArr = append(permitArr, permits)
	err = p.permitRepo.BatchCreate(p.db, permitArr...)
	if err != nil {
		return nil, err
	}
	return &CreatePerResp{}, nil
}

type check struct {
	method string
	path   string
}

func (p *permit) checkExist(ctx context.Context, req *CreatePerReq) (bool, error) {
	checks := make([]*check, 0)
	if IsFormAPI(req.AccessPath) {
		checks = append(checks, &check{
			method: req.Method,
			path:   req.URI,
		})
	}
	checks = append(checks, &check{
		method: req.Method,
		path:   req.AccessPath,
	})
	for _, value := range checks {
		exit, err := p.permitRepo.Get(p.db, req.RoleID, value.path, value.method)
		if err != nil {
			return false, err
		}
		if exit.ID != "" {
			return true, nil
		}
		continue
	}
	return false, nil
}

type UpdatePerReq struct {
	ID          string             `json:"id"`
	Params      models.FiledPermit `json:"params"`
	Response    models.FiledPermit `json:"response"`
	Condition   models.Condition   `json:"condition"`
	ParamsAll   bool               `json:"paramsAll"`
	ResponseAll bool               `json:"responseAll"`
	Path        string             `json:"accessPath"`
	URI         string             `json:"uri"`
	Method      string             `json:"method"`
}

type UpdatePerResp struct{}

func (p *permit) UpdatePermit(ctx context.Context, req *UpdatePerReq) (*UpdatePerResp, error) {
	tx := p.db.Begin()
	if IsFormAPI(req.Path) {
		err := p.permitRepo.Update(p.db, &models.PermitQuery{
			Path:   req.URI,
			Method: req.Method,
		}, &models.Permit{
			Params:      req.Params,
			Response:    req.Response,
			Condition:   req.Condition,
			ParamsAll:   req.ParamsAll,
			ResponseAll: req.ResponseAll,
		})
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	err := p.permitRepo.Update(p.db, &models.PermitQuery{
		Path:   req.Path,
		Method: req.Method,
	}, &models.Permit{
		Params:      req.Params,
		Response:    req.Response,
		Condition:   req.Condition,
		ParamsAll:   req.ParamsAll,
		ResponseAll: req.ResponseAll,
	})
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()
	// add redis cache
	return &UpdatePerResp{}, nil
}

type DeletePerReq struct {
	RoleID string `json:"roleID"`
	Path   string `json:"path"`
	URI    string `json:"uri"`
	Method string `json:"method"`
}

type DeletePerResp struct{}

func (p *permit) DeletePermit(ctx context.Context, req *DeletePerReq) (*DeletePerResp, error) {
	if IsFormAPI(req.Path) {
		err := p.permitRepo.Delete(p.db, &models.PermitQuery{
			RoleID: req.RoleID,
			Path:   req.URI,
			Method: req.Method,
		})
		if err != nil {
			return nil, err
		}
	}
	err := p.permitRepo.Delete(p.db, &models.PermitQuery{
		RoleID: req.RoleID,
		Path:   req.Path,
		Method: req.Method,
	})
	if err != nil {
		return nil, err
	}
	return &DeletePerResp{}, nil
}

type GetPermitReq struct {
	RoleID string `json:"roleID"`
	Path   string `json:"path"`
	URI    string `json:"uri"`
	Method string `json:"method" binding:"required"`
}

type GetPermitResp struct {
	ID          string             `json:"id"`
	RoleID      string             `json:"roleID"`
	Path        string             `json:"path,omitempty"`
	Params      models.FiledPermit `json:"params,omitempty"`
	Response    models.FiledPermit `json:"response,omitempty"`
	Condition   models.Condition   `json:"condition,omitempty"`
	ResponseAll bool               `json:"responseAll"`
	ParamsAll   bool               `json:"paramsAll"`
}

func (p *permit) GetPermit(ctx context.Context, req *GetPermitReq) (*GetPermitResp, error) {
	if IsFormAPI(req.Path) && req.URI != "" {
		req.Path = req.URI
	}

	permits, err := p.permitRepo.Get(p.db, req.RoleID, req.Path, req.Method)
	if err != nil {
		return nil, err
	}
	return &GetPermitResp{
		ID:          permits.ID,
		RoleID:      permits.RoleID,
		Path:        permits.Path,
		Params:      permits.Params,
		Response:    permits.Response,
		Condition:   permits.Condition,
		ResponseAll: permits.ResponseAll,
		ParamsAll:   permits.ParamsAll,
	}, nil
}

type DeleteRoleReq struct {
	RoleID string `json:"-"`
	AppID  string `json:"-"`
}
type DeleteRoleResp struct{}

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
	err = p.userRoleRepo.Delete(p.db, &models.UserRoleQuery{
		AppID:  req.AppID,
		RoleID: req.RoleID,
	})
	if err != nil {
		return nil, err
	}
	return &DeleteRoleResp{}, nil
}

type GetUserRoleReq struct {
	UserID string `json:"userID"`
	DepID  string `json:"depID"`
	AppID  string `json:"appID"`
}

type GetUserRoleResp struct {
	RoleID string          `json:"id"`
	Types  models.RoleType `json:"type"`
}

func (p *permit) GetUserRole(ctx context.Context, req *GetUserRoleReq) (*GetUserRoleResp, error) {
	userRole, err := p.userRoleRepo.Get(p.db, req.AppID, req.UserID)
	if err != nil {
		return nil, err
	}
	resp := &GetUserRoleResp{
		RoleID: userRole.RoleID,
	}
	if userRole.RoleID != "" {
		resp.RoleID = userRole.RoleID
		roles, err := p.roleRepo.Get(p.db, userRole.RoleID)
		if err != nil {
			return nil, err
		}
		resp.Types = roles.Types
		return resp, nil
	}
	// 根据
	ow := make([]string, 0)
	if req.UserID != "" {
		ow = append(ow, req.UserID)
	}
	if req.DepID != "" {
		ow = append(ow, req.DepID)
	}
	grant, total, err := p.roleGrantRepo.List(p.db, &models.RoleGrantQuery{
		Owners: ow,
		AppID:  req.AppID,
	}, 1, 999)
	if total == 0 || len(grant) == 0 {
		return resp, nil
	}
	// get role
	role, err := p.roleRepo.Get(p.db, grant[0].RoleID)
	if err != nil {
		return nil, err
	}

	// save user role
	err = p.userRoleRepo.BatchCreate(p.db, &models.UserRole{
		UserID: req.UserID,
		RoleID: role.ID,
		AppID:  req.AppID,
		ID:     id2.StringUUID(),
	})
	if err != nil {
		return nil, err
	}
	resp.Types = role.Types
	resp.RoleID = role.ID
	return resp, nil
}

func (p *permit) publish(ctx context.Context, topic string, data interface{}) error {
	return nil
}
