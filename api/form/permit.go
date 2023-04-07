package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/quanxiang-cloud/cabin/logger"
	"github.com/quanxiang-cloud/cabin/tailormade/header"
	"github.com/quanxiang-cloud/cabin/tailormade/resp"
	"github.com/quanxiang-cloud/form/internal/service"
	config2 "github.com/quanxiang-cloud/form/pkg/misc/config"
)

// Permit permit.
type Permit struct {
	permit service.Permit
}

// NewPermit new permit.
func NewPermit(conf *config2.Config) (*Permit, error) {
	permits, err := service.NewPermit(conf)
	if err != nil {
		return nil, err
	}

	return &Permit{
		permit: permits,
	}, nil
}

// CreateRole create role.
func (p *Permit) CreateRole(c *gin.Context) {
	req := &service.CreateRoleReq{
		AppID:    c.Param(_appID),
		UserID:   c.GetHeader(_userID),
		UserName: c.GetHeader(_userName),
	}
	ctx := header.MutateContext(c)
	if err := c.ShouldBind(req); err != nil {
		logger.Logger.WithName("CreateRole").Errorw(err.Error(), header.GetRequestIDKV(ctx).Fuzzy()...)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	resp.Format(p.permit.CreateRole(ctx, req)).Context(c)
}

// UpdateRole update role.
func (p *Permit) UpdateRole(c *gin.Context) {
	req := &service.UpdateRoleReq{}
	ctx := header.MutateContext(c)
	if err := c.ShouldBind(req); err != nil {
		logger.Logger.WithName("UpdateRole").Errorw(err.Error(), header.GetRequestIDKV(ctx).Fuzzy()...)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	resp.Format(p.permit.UpdateRole(ctx, req)).Context(c)
}

// DeleteRole delete role.
func (p *Permit) DeleteRole(c *gin.Context) {
	req := &service.DeleteRoleReq{
		RoleID: c.Param("id"),
	}
	ctx := header.MutateContext(c)
	if err := c.ShouldBind(req); err != nil {
		logger.Logger.WithName("DeleteRole").Errorw(err.Error(), header.GetRequestIDKV(ctx).Fuzzy()...)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	resp.Format(p.permit.DeleteRole(ctx, req)).Context(c)
}

// CratePermit create permit.
func (p *Permit) CratePermit(c *gin.Context) {
	pr := getProfile(c)
	req := &service.CreatePerReq{
		UserID:   pr.userID,
		UserName: pr.userName,
	}
	ctx := header.MutateContext(c)
	if err := c.ShouldBind(req); err != nil {
		logger.Logger.WithName("CratePermit").Errorw(err.Error(), header.GetRequestIDKV(ctx).Fuzzy()...)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	resp.Format(p.permit.CreatePermit(ctx, req)).Context(c)
}

// UpdatePermit update permit.
func (p *Permit) UpdatePermit(c *gin.Context) {
	req := &service.UpdatePerReq{
		ID: c.Param("id"),
	}
	ctx := header.MutateContext(c)
	if err := c.ShouldBind(req); err != nil {
		logger.Logger.WithName("UpdatePermit").Errorw(err.Error(), header.GetRequestIDKV(ctx).Fuzzy()...)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	resp.Format(p.permit.UpdatePermit(ctx, req)).Context(c)
}

// GetPermit get permit.
func (p *Permit) GetPermit(c *gin.Context) {
	req := &service.GetPermitReq{}
	ctx := header.MutateContext(c)
	if err := c.ShouldBind(req); err != nil {
		logger.Logger.WithName("GetPermit").Errorw(err.Error(), header.GetRequestIDKV(ctx).Fuzzy()...)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	resp.Format(p.permit.GetPermit(ctx, req)).Context(c)
}

// FindPermit find permit.
func (p *Permit) FindPermit(c *gin.Context) {
	req := &service.FindPermitReq{}
	ctx := header.MutateContext(c)
	if err := c.ShouldBind(req); err != nil {
		logger.Logger.WithName("FindPermit").Errorw(err.Error(), header.GetRequestIDKV(ctx).Fuzzy()...)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	resp.Format(p.permit.FindPermit(ctx, req)).Context(c)
}

// FindRole find role.
func (p *Permit) FindRole(c *gin.Context) {
	req := &service.FindRoleReq{
		AppID: c.Param(_appID),
	}
	ctx := header.MutateContext(c)
	//
	if err := c.ShouldBind(req); err != nil {
		logger.Logger.WithName("FindRole").Errorw(err.Error(), header.GetRequestIDKV(ctx).Fuzzy()...)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	resp.Format(p.permit.FindRole(ctx, req)).Context(c)
}

// GetRole get role.
func (p *Permit) GetRole(c *gin.Context) {
	req := &service.GetRoleReq{
		ID: c.Param("id"),
	}
	ctx := header.MutateContext(c)
	if err := c.ShouldBind(req); err != nil {
		logger.Logger.WithName("GetRole").Errorw(err.Error(), header.GetRequestIDKV(ctx).Fuzzy()...)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	resp.Format(p.permit.GetRole(ctx, req)).Context(c)
}

// FindGrantRole find grant role.
func (p *Permit) FindGrantRole(c *gin.Context) {
	req := &service.FindGrantRoleReq{
		AppID:  c.Param(_appID),
		RoleID: c.Param("roleID"),
	}
	ctx := header.MutateContext(c)
	if err := c.ShouldBind(req); err != nil {
		logger.Logger.WithName("FindGrantRole").Errorw(err.Error(), header.GetRequestIDKV(ctx).Fuzzy()...)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	resp.Format(p.permit.FindGrantRole(ctx, req)).Context(c)
}

// AssignRoleGrant assign role grant.
func (p *Permit) AssignRoleGrant(c *gin.Context) {
	req := &service.AssignRoleGrantReq{
		AppID:  c.Param(_appID),
		RoleID: c.Param("roleID"),
	}
	ctx := header.MutateContext(c)
	if err := c.ShouldBind(req); err != nil {
		logger.Logger.WithName("AssignRoleGrant").Errorw(err.Error(), header.GetRequestIDKV(ctx).Fuzzy()...)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	resp.Format(p.permit.AssignRoleGrant(ctx, req)).Context(c)
}

// DeletePermit delete permit.
func (p *Permit) DeletePermit(c *gin.Context) {
	req := &service.DeletePerReq{}
	ctx := header.MutateContext(c)
	if err := c.ShouldBind(req); err != nil {
		logger.Logger.WithName("DeletePermit").Errorw(err.Error(), header.GetRequestIDKV(ctx).Fuzzy()...)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	resp.Format(p.permit.DeletePermit(ctx, req)).Context(c)
}

// ListPermit list permit.
func (p *Permit) ListPermit(c *gin.Context) {
	req := &service.ListPermitReq{}
	ctx := header.MutateContext(c)
	if err := c.ShouldBind(req); err != nil {
		logger.Logger.WithName("ListPermit").Errorw(err.Error(), header.GetRequestIDKV(ctx).Fuzzy()...)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	resp.Format(p.permit.ListPermit(ctx, req)).Context(c)
}

// ListAndSelect list and select.
func (p *Permit) ListAndSelect(c *gin.Context) {
	pf := getProfile(c)
	req := &service.ListAndSelectReq{
		AppID:  c.Param("appID"),
		UserID: pf.userID,
		DepID:  pf.depID,
	}
	ctx := header.MutateContext(c)
	if err := c.ShouldBind(req); err != nil {
		logger.Logger.WithName("ListAndSelect").Errorw(err.Error(), header.GetRequestIDKV(ctx).Fuzzy()...)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	resp.Format(p.permit.ListAndSelect(ctx, req)).Context(c)
}

func (p *Permit) PathPermit(c *gin.Context) {
	pf := getProfile(c)
	req := &service.HomePerListReq{
		AppID:  c.Param(_appID),
		UserID: pf.userID,
		DepID:  pf.depID,
	}
	ctx := header.MutateContext(c)
	if err := c.ShouldBind(req); err != nil {
		logger.Logger.WithName("HomePerList").Errorw(err.Error(), header.GetRequestIDKV(ctx).Fuzzy()...)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	resp.Format(p.permit.HomePerList(ctx, req)).Context(c)
}

func (p *Permit) CreateUserRole(c *gin.Context) {
	pf := getProfile(c)
	req := &service.CreateUserRoleReq{
		UserID: pf.userID,
		AppID:  c.Param("appID"),
	}
	ctx := header.MutateContext(c)
	if err := c.ShouldBind(req); err != nil {
		logger.Logger.WithName("SaveUserPerMatch").Errorw(err.Error(), header.GetRequestIDKV(ctx).Fuzzy()...)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	resp.Format(p.permit.CreateUserRole(ctx, req)).Context(c)
}

func (p *Permit) GetUserRole(c *gin.Context) {
	pf := getProfile(c)
	req := &service.GetUserRoleReq{
		UserID: pf.userID,
		DepID:  pf.depID,
		AppID:  c.Param("appID"),
	}
	ctx := header.MutateContext(c)
	if err := c.ShouldBind(req); err != nil {
		logger.Logger.WithName("SaveUserPerMatch").Errorw(err.Error(), header.GetRequestIDKV(ctx).Fuzzy()...)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	resp.Format(p.permit.GetUserRole(ctx, req)).Context(c)
}

func (p *Permit) CopyRole(c *gin.Context) {
	pf := getProfile(c)
	req := &service.CopyRoleReq{
		UserID:   pf.userID,
		UserName: pf.userName,
		AppID:    c.Param("appID"),
	}
	ctx := header.MutateContext(c)
	if err := c.ShouldBind(req); err != nil {
		logger.Logger.WithName("SaveUserPerMatch").Errorw(err.Error(), header.GetRequestIDKV(ctx).Fuzzy()...)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	resp.Format(p.permit.CopyRole(ctx, req)).Context(c)
}

func (p *Permit) PerPoly(c *gin.Context) {
	req := &service.PerPolyReq{
		AppID: c.Param("appID"),
	}
	ctx := header.MutateContext(c)
	if err := c.ShouldBind(req); err != nil {
		logger.Logger.WithName("SaveUserPerMatch").Errorw(err.Error(), header.GetRequestIDKV(ctx).Fuzzy()...)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	resp.Format(p.permit.PerPoly(ctx, req)).Context(c)

}
