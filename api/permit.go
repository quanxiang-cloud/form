package api

import (
	"github.com/gin-gonic/gin"
	"github.com/quanxiang-cloud/cabin/logger"
	"github.com/quanxiang-cloud/cabin/tailormade/header"
	"github.com/quanxiang-cloud/cabin/tailormade/resp"
	"github.com/quanxiang-cloud/form/internal/service"
	config2 "github.com/quanxiang-cloud/form/pkg/misc/config"
	"net/http"
)

type Permit struct {
	permit service.Permit
}

// NewPermit new permit
func NewPermit(conf *config2.Config) (*Permit, error) {

	permits, err := service.NewPermit(conf)
	if err != nil {
		return nil, err
	}
	return &Permit{
		permit: permits,
	}, nil
}

func (p *Permit) CreateRole(c *gin.Context) {
	req := &service.CreateRoleReq{
		AppID:    c.Param(_appID),
		UserID:   c.GetHeader(_userID),
		UserName: c.GetHeader(_userName),
	}
	ctx := header.MutateContext(c)
	if err := c.ShouldBind(req); err != nil {
		logger.Logger.Errorw("should bind", header.GetRequestIDKV(ctx).Fuzzy(), err.Error())
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	resp.Format(p.permit.CreateRole(ctx, req)).Context(c)
}

func (p *Permit) UpdateRole(c *gin.Context) {
	req := &service.UpdateRoleReq{}
	ctx := header.MutateContext(c)
	if err := c.ShouldBind(req); err != nil {
		logger.Logger.Errorw("should bind", header.GetRequestIDKV(ctx).Fuzzy(), err.Error())
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	resp.Format(p.permit.UpdateRole(ctx, req)).Context(c)
}

func (p *Permit) AddToRole(c *gin.Context) {
	req := &service.AddOwnerToRoleReq{}
	ctx := header.MutateContext(c)
	if err := c.ShouldBind(req); err != nil {
		logger.Logger.Errorw("should bind", header.GetRequestIDKV(ctx).Fuzzy(), err.Error())
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	resp.Format(p.permit.AddOwnerToRole(ctx, req)).Context(c)
}

func (p *Permit) DeleteOwner(c *gin.Context) {
	req := &service.DeleteOwnerReq{}
	ctx := header.MutateContext(c)
	if err := c.ShouldBind(req); err != nil {
		logger.Logger.Errorw("should bind", header.GetRequestIDKV(ctx).Fuzzy(), err.Error())
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	resp.Format(p.permit.DeleteOwnerToRole(ctx, req)).Context(c)
}

func (p *Permit) DeleteRole(c *gin.Context) {
	req := &service.DeleteRoleReq{}
	ctx := header.MutateContext(c)
	if err := c.ShouldBind(req); err != nil {
		logger.Logger.Errorw("should bind", header.GetRequestIDKV(ctx).Fuzzy(), err.Error())
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	resp.Format(p.permit.DeleteRole(ctx, req)).Context(c)
}

func (p *Permit) CratePermit(c *gin.Context) {
	req := &service.CreatePerReq{}
	ctx := header.MutateContext(c)
	if err := c.ShouldBind(req); err != nil {
		logger.Logger.Errorw("should bind", header.GetRequestIDKV(ctx).Fuzzy(), err.Error())
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	resp.Format(p.permit.CreatePermit(ctx, req)).Context(c)
}

func (p *Permit) UpdatePermit(c *gin.Context) {
	req := &service.UpdatePerReq{}
	ctx := header.MutateContext(c)
	if err := c.ShouldBind(req); err != nil {
		logger.Logger.Errorw("should bind", header.GetRequestIDKV(ctx).Fuzzy(), err.Error())
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	resp.Format(p.permit.UpdatePermit(ctx, req)).Context(c)
}

func (p *Permit) GetPermit(c *gin.Context) {
	req := &service.GetPermitReq{}
	ctx := header.MutateContext(c)
	if err := c.ShouldBind(req); err != nil {
		logger.Logger.Errorw("should bind", header.GetRequestIDKV(ctx).Fuzzy(), err.Error())
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	resp.Format(p.permit.GetPermit(ctx, req)).Context(c)
}

func (p *Permit) SaveUserPerMatch(c *gin.Context) {
	req := &service.SaveUserPerMatchReq{}
	ctx := header.MutateContext(c)
	if err := c.ShouldBind(req); err != nil {
		logger.Logger.Errorw("should bind", header.GetRequestIDKV(ctx).Fuzzy(), err.Error())
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	resp.Format(p.permit.SaveUserPerMatch(ctx, req)).Context(c)
}

func (p *Permit) UserRoleMatch(c *gin.Context) {

	req := &service.FindGrantRoleReq{}
	ctx := header.MutateContext(c)
	if err := c.ShouldBind(req); err != nil {
		logger.Logger.Errorw("should bind", header.GetRequestIDKV(ctx).Fuzzy(), err.Error())
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	match, err := p.permit.FindGrantRole(ctx, req)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}
	if len(match.List) == 0 {
		c.AbortWithError(http.StatusForbidden, err)
	}
	reqRole := &service.GetRoleReq{
		ID: match.List[0].ID,
	}
	resp.Format(p.permit.GetRole(ctx, reqRole)).Context(c)
}
