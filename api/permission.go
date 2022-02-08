package api

import (
	"github.com/gin-gonic/gin"
	"github.com/quanxiang-cloud/cabin/tailormade/header"
	"github.com/quanxiang-cloud/cabin/tailormade/resp"
	"github.com/quanxiang-cloud/form/internal/service"
	"github.com/quanxiang-cloud/form/pkg/misc/config"
	"net/http"
	"strings"
)

type Permission struct {
	permission service.Permission
}

// NewPermission NewPermission
func NewPermission(config *config.Config) (*Permission, error) {
	permission, err := service.NewPermission(config)
	if err != nil {
		return nil, err
	}
	return &Permission{
		permission: permission,
	}, nil
}

// CreatePerGroup CreatePerGroup
func (per *Permission) CreatePerGroup(c *gin.Context) {

	req := &service.CreateGroup{
		AppID:       c.Param("appID"),
		UserID:      c.GetHeader(_userID),
		CreatorName: c.GetHeader(_userName),
	}
	if err := c.ShouldBind(req); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	resp.Format(per.permission.CreateGroup(header.MutateContext(c), req)).Context(c)
}

// UpdatePerGroup UpdatePerGroup
func (per *Permission) UpdatePerGroup(c *gin.Context) {
	req := &service.UpdateGroupReq{
		AppID: c.Param("appID"),
	}
	if err := c.ShouldBind(req); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	resp.Format(per.permission.UpdateGroup(header.MutateContext(c), req)).Context(c)
}

func (per *Permission) DelPerGroup(c *gin.Context) {
	req := &service.DelPerGroupReq{}
	if err := c.ShouldBind(req); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	req.AppID = c.Param(_appID)

	resp.Format(per.permission.DelPerGroup(header.MutateContext(c), req)).Context(c)
}

func (per *Permission) GetPerGroup(c *gin.Context) {
	req := &service.GetPerGroupReq{}
	if err := c.ShouldBind(req); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	resp.Format(per.permission.GetPerGroup(header.MutateContext(c), req)).Context(c)
}

func (per *Permission) FindPerGroup(c *gin.Context) {
	req := &service.FindPerGroupReq{
		AppID: c.Param(_appID),
	}
	if err := c.ShouldBind(req); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	resp.Format(per.permission.FindPerGroup(header.MutateContext(c), req)).Context(c)
}

func (per *Permission) DeleteForm(c *gin.Context) {

	req := &service.DeleteFormReq{}
	if err := c.ShouldBind(req); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	resp.Format(per.permission.DeleteForm(header.MutateContext(c), req)).Context(c)
}

func (per *Permission) SaveForm(c *gin.Context) {
	req := &service.SaveFormReq{}
	if err := c.ShouldBind(req); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	resp.Format(per.permission.SaveForm(header.MutateContext(c), req)).Context(c)
}

func (per *Permission) FindForm(c *gin.Context) {
	req := &service.FindFormReq{}
	if err := c.ShouldBind(req); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	resp.Format(per.permission.FindForm(header.MutateContext(c), req)).Context(c)
}

func (per *Permission) GetPerInfo(c *gin.Context) {
	req := &service.GetFormReq{}
	if err := c.ShouldBind(req); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	resp.Format(per.permission.GetForm(header.MutateContext(c), req)).Context(c)
}

func (per *Permission) GetPerGroupsByMenu(c *gin.Context) {
	req := &service.GetGroupsByMenuReq{}
	if err := c.ShouldBind(req); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	resp.Format(per.permission.GetGroupsByMenu(header.MutateContext(c), req)).Context(c)
}

func (per *Permission) GetOperate(c *gin.Context) {
	depIDS := strings.Split(c.GetHeader(_departmentID), ",")
	req := &service.GetPerInCacheReq{
		UserID: c.GetHeader(_userID),
		DepID:  depIDS[len(depIDS)-1],
		AppID:  c.Param(_appID),
	}
	if err := c.ShouldBind(req); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	cacheResp, err := per.permission.GetPerInCache(header.MutateContext(c), req)
	result := struct {
		ID        string `json:"id"`
		Authority int64  `json:"authority"`
	}{}
	if cacheResp != nil {
		result.ID = cacheResp.ID
		result.Authority = cacheResp.Authority
	}
	resp.Format(cacheResp, err).Context(c)
}

func (per *Permission) GetPerOption(c *gin.Context) {
	ctx := header.MutateContext(c)
	userID := c.GetHeader(_userID)
	depIDS := strings.Split(c.GetHeader(_departmentID), ",")
	appID := c.Param(_appID)
	selectPer, err := per.permission.GetPerInCache(ctx, &service.GetPerInCacheReq{
		AppID:  appID,
		UserID: userID,
		DepID:  depIDS[len(depIDS)-1],
	})
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	optionPer, err := per.permission.FindPerGroup(ctx, &service.FindPerGroupReq{
		AppID:  appID,
		DepID:  depIDS[len(depIDS)-1],
		UserID: userID,
	})
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	resp.Format(map[string]interface{}{
		"selectPer": struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		}{
			ID:   selectPer.ID,
			Name: selectPer.Name,
		},
		"optionPer": optionPer.ListVO,
	}, nil).Context(c)

}

func (per *Permission) SaveUserPerMatch(c *gin.Context) {
	req := &service.SaveUserPerMatchReq{
		UserID: c.GetHeader(_userID),
		AppID:  c.Param(_appID),
	}
	if err := c.ShouldBind(req); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	resp.Format(per.permission.SaveUserPerMatch(header.MutateContext(c), req)).Context(c)
}

// AddOwnerToGroup user and dep to group
func (per *Permission) AddOwnerToGroup(c *gin.Context) {
	req := &service.AddOwnerToGroupReq{
		AppID: c.Param(_appID),
	}
	if err := c.ShouldBind(req); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	resp.Format(per.permission.AddOwnerToGroup(header.MutateContext(c), req)).Context(c)
}

// AddOwnerToApp user and dep to app
func (per *Permission) AddOwnerToApp(c *gin.Context) {
	req := &service.AddOwnerToAppReq{
		AppID: c.Param(_appID),
	}
	if err := c.ShouldBind(req); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	resp.Format(per.permission.AddOwnerToApp(header.MutateContext(c), req)).Context(c)
}
