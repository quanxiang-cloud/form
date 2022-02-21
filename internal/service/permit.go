package service

import "context"

type SaveUserPerMatchReq struct {
}

type SaveUserPerMatchResp struct {
}

type FindGrantRoleReq struct {
}

type FindGrantRoleResp struct {
}

type Permit interface {
	CreateRole(ctx context.Context, req *CreateRoleReq) (*CreateRoleResp, error)

	UpdateRole(ctx context.Context, req *UpdateRoleReq) (*UpdateRoleResp, error)

	DeleteRole(ctx context.Context, req *DeleteRoleReq) (*DeleteRoleResp, error) // 这个删除需要关心的东西比较多

	GetRole(ctx context.Context, req *GetRoleReq) (*GetRoleResp, error)

	FindRole(ctx context.Context, req *FindRoleReq) (*FindRoleResp, error)

	AddOwnerToRole(ctx context.Context, req *AddOwnerToRoleReq) (*AddOwnerToRoleResp, error)

	DeleteOwnerToRole(ctx context.Context, req *DeleteOwnerReq) (*DeleteOwnerResp, error)

	FindGrantRole(ctx context.Context, req *FindGrantRoleReq) (*FindGrantRoleResp, error)

	CreatePermit(ctx context.Context, req *CreatePerReq) (*CreatePerResp, error)

	UpdatePermit(ctx context.Context, req *UpdatePerReq) (*UpdatePerResp, error)

	DeletePermit(ctx context.Context, req *DeletePerReq) (*DeletePerResp, error)

	GetPerInCache(ctx context.Context, req *GetPerInCacheReq) (*GetPerInCacheResp, error)

	SaveUserPerMatch(ctx context.Context, req *SaveUserPerMatchReq) (*SaveUserPerMatchResp, error)
}

type permit struct {
}

func (p *permit) FindGrantRole(ctx context.Context, req *FindGrantRoleReq) (*FindGrantRoleResp, error) {
	panic("implement me")
}

func (p *permit) SaveUserPerMatch(ctx context.Context, req *SaveUserPerMatchReq) (*SaveUserPerMatchResp, error) {
	panic("implement me")
}

func NewPermit() Permit {
	return &permit{}
}

type CreateRoleReq struct {
}

type CreateRoleResp struct {
}

func (p *permit) CreateRole(ctx context.Context, req *CreateRoleReq) (*CreateRoleResp, error) {
	return nil, nil
}

type UpdateRoleReq struct {
}

type UpdateRoleResp struct {
}

func (p *permit) UpdateRole(ctx context.Context, req *UpdateRoleReq) (*UpdateRoleResp, error) {
	return nil, nil
}

type DeleteRoleReq struct {
}
type DeleteRoleResp struct {
}

func (p *permit) DeleteRole(ctx context.Context, req *DeleteRoleReq) (*DeleteRoleResp, error) {
	return nil, nil
}

type GetRoleReq struct {
}
type GetRoleResp struct {
}

func (p *permit) GetRole(ctx context.Context, req *GetRoleReq) (*GetRoleResp, error) {
	return nil, nil
}

type FindRoleReq struct {
}

type FindRoleResp struct {
}

func (p *permit) FindRole(ctx context.Context, req *FindRoleReq) (*FindRoleResp, error) {
	return nil, nil
}

type AddOwnerToRoleReq struct {
}

type AddOwnerToRoleResp struct {
}

func (p *permit) AddOwnerToRole(ctx context.Context, req *AddOwnerToRoleReq) (*AddOwnerToRoleResp, error) {
	return nil, nil
}

type DeleteOwnerReq struct {
}

type DeleteOwnerResp struct {
}

func (p *permit) DeleteOwnerToRole(ctx context.Context, req *DeleteOwnerReq) (*DeleteOwnerResp, error) {
	return nil, nil
}

type CreatePerReq struct {
}

type CreatePerResp struct {
}

func (p *permit) CreatePermit(ctx context.Context, req *CreatePerReq) (*CreatePerResp, error) {
	return nil, nil
}

type UpdatePerReq struct {
}

type UpdatePerResp struct {
}

func (p *permit) UpdatePermit(ctx context.Context, req *UpdatePerReq) (*UpdatePerResp, error) {
	return nil, nil
}

type DeletePerReq struct {
}

type DeletePerResp struct {
}

func (p *permit) DeletePermit(ctx context.Context, req *DeletePerReq) (*DeletePerResp, error) {
	return nil, nil
}

type GetPerInCacheReq struct {
}

type GetPerInCacheResp struct {
}

func (p *permit) GetPerInCache(ctx context.Context, req *GetPerInCacheReq) (*GetPerInCacheResp, error) {
	return nil, nil
}
