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

func (p *Permit) DeleteRole(c *gin.Context) {
	req := &service.DeleteRoleReq{
		RoleID: c.Param("id"),
	}
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
	req := &service.UpdatePerReq{
		ID: c.Param("id"),
	}
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
		ID: match.List[0].RoleID,
	}
	resp.Format(p.permit.GetRole(ctx, reqRole)).Context(c)
}

func (p *Permit) FindPermit(c *gin.Context) {
	req := &service.FindPermitReq{}
	ctx := header.MutateContext(c)
	if err := c.ShouldBind(req); err != nil {
		logger.Logger.Errorw("should bind", header.GetRequestIDKV(ctx).Fuzzy(), err.Error())
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	resp.Format(p.permit.FindPermit(ctx, req)).Context(c)
}

func (p *Permit) FindRole(c *gin.Context) {
	req := &service.FindRoleReq{
		AppID: c.Param(_appID),
	}
	ctx := header.MutateContext(c)
	if err := c.ShouldBind(req); err != nil {
		logger.Logger.Errorw("should bind", header.GetRequestIDKV(ctx).Fuzzy(), err.Error())
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	resp.Format(p.permit.FindRole(ctx, req)).Context(c)
}

func (p *Permit) GetRole(c *gin.Context) {
	req := &service.GetRoleReq{
		ID: c.Param("id"),
	}
	ctx := header.MutateContext(c)
	if err := c.ShouldBind(req); err != nil {
		logger.Logger.Errorw("should bind", header.GetRequestIDKV(ctx).Fuzzy(), err.Error())
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	resp.Format(p.permit.GetRole(ctx, req)).Context(c)

}

func (p *Permit) FindGrantRole(c *gin.Context) {
	req := &service.FindGrantRoleReq{
		AppID:  c.Param(_appID),
		RoleID: c.Param("roleID"),
	}
	ctx := header.MutateContext(c)
	if err := c.ShouldBind(req); err != nil {
		logger.Logger.Errorw("should bind", header.GetRequestIDKV(ctx).Fuzzy(), err.Error())
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	resp.Format(p.permit.FindGrantRole(ctx, req)).Context(c)
}

func (p *Permit) AssignRoleGrant(c *gin.Context) {
	req := &service.AssignRoleGrantReq{
		AppID:  c.Param(_appID),
		RoleID: c.Param("roleID"),
	}
	ctx := header.MutateContext(c)
	if err := c.ShouldBind(req); err != nil {
		logger.Logger.Errorw("should bind", header.GetRequestIDKV(ctx).Fuzzy(), err.Error())
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	resp.Format(p.permit.AssignRoleGrant(ctx, req)).Context(c)
}

func (p *Permit) DeletePermit(c *gin.Context) {
	req := &service.DeletePerReq{}
	ctx := header.MutateContext(c)
	if err := c.ShouldBind(req); err != nil {
		logger.Logger.Errorw("should bind", header.GetRequestIDKV(ctx).Fuzzy(), err.Error())
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	resp.Format(p.permit.DeletePermit(ctx, req)).Context(c)
}

func (p *Permit) ListPermit(c *gin.Context) {
	ctx := header.MutateContext(c)
	var batch []service.GetPermitReq
	if err := c.ShouldBind(&batch); err != nil {
		logger.Logger.Errorw("should bind", header.GetRequestIDKV(ctx).Fuzzy(), err.Error())
		resp.Format(nil, err).Context(c, http.StatusBadRequest)
		return
	}

	list := make([]*service.GetPermitResp, len(batch))
	for index, get := range batch {
		r, err := p.permit.GetPermit(c, &get)
		if err != nil {
			logger.Logger.Errorw("get ", header.GetRequestIDKV(ctx).Fuzzy(), err.Error())
			list[index] = &service.GetPermitResp{}
			continue
		}
		list[index] = r
	}
	resp.Format(map[string]interface{}{
		"list": list,
	}, nil).Context(c)
}
