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

// Project  pr.
type Project struct {
	project service.Project
}

// NewProject new project.
func NewProject(conf *config2.Config) (*Project, error) {
	pro, err := service.NewProject(conf)
	if err != nil {
		return nil, err
	}
	project := &Project{
		project: pro,
	}
	return project, nil
}

// CreateProject create table.
func (p *Project) CreateProject(c *gin.Context) {
	profiles := getProfile(c)
	req := &service.CreateProjectReq{
		CreatorID:   profiles.userID,
		CreatorName: profiles.userName,
	}

	ctx := header.MutateContext(c)
	if err := c.ShouldBind(req); err != nil {
		logger.Logger.WithName("CreateProject").Errorw(err.Error(), header.GetRequestIDKV(ctx).Fuzzy()...)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	resp.Format(p.project.CreateProject(ctx, req)).Context(c)
}

func (p *Project) DeleteProject(c *gin.Context) {
	req := &service.DeleteProjectReq{}
	ctx := header.MutateContext(c)
	if err := c.ShouldBind(req); err != nil {
		logger.Logger.WithName("CreateProject").Errorw(err.Error(), header.GetRequestIDKV(ctx).Fuzzy()...)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	resp.Format(p.project.DeleteProject(ctx, req)).Context(c)
}

// ListProject create table.
func (p *Project) ListProject(c *gin.Context) {
	req := &service.ListProjectReq{}
	ctx := header.MutateContext(c)
	if err := c.ShouldBind(req); err != nil {
		logger.Logger.WithName("CreateProject").Errorw(err.Error(), header.GetRequestIDKV(ctx).Fuzzy()...)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	resp.Format(p.project.ListProject(ctx, req)).Context(c)
}

// AssignProjectUser assign .
func (p *Project) AssignProjectUser(c *gin.Context) {
	req := &service.AssignProjectUserReq{}
	ctx := header.MutateContext(c)
	if err := c.ShouldBind(req); err != nil {
		logger.Logger.WithName("AssignProjectUser").Errorw(err.Error(), header.GetRequestIDKV(ctx).Fuzzy()...)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	resp.Format(p.project.AssignProjectUser(ctx, req)).Context(c)
}

// ListProjectUser  listProjectUser
func (p *Project) ListProjectUser(c *gin.Context) {
	req := &service.ListProjectUserReq{}
	ctx := header.MutateContext(c)
	if err := c.ShouldBind(req); err != nil {
		logger.Logger.WithName("ListProjectUser").Errorw(err.Error(), header.GetRequestIDKV(ctx).Fuzzy()...)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	resp.Format(p.project.ListProjectUser(ctx, req)).Context(c)
}
