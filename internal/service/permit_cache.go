package service

import (
	"context"

	"github.com/quanxiang-cloud/cabin/logger"
	"github.com/quanxiang-cloud/form/internal/models"
)

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

	err = p.limitRepo.DeletePermit(ctx, req.RoleID)
	if err != nil {

	}
	return &DeleteRoleResp{}, nil

}
