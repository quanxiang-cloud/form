package api

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/quanxiang-cloud/cabin/tailormade/header"
	"github.com/quanxiang-cloud/cabin/tailormade/resp"
	"github.com/quanxiang-cloud/form/internal/service/form"
)

// CheckURL CheckURL
func checkURL(c *gin.Context) (appID, tableName string, err error) {
	appID, ok := c.Params.Get("appID")
	tableName, okt := c.Params.Get("tableName")
	if !ok || !okt {
		err = errors.New("invalid URI")
		return
	}
	return

}

func Search(f form.Form) gin.HandlerFunc {
	return func(c *gin.Context) {
		req := &form.SearchReq{}
		p := getProfile(c)
		req.UserID = p.userID
		req.DepID = p.depID
		var err error
		req.AppID, req.TableID, err = checkURL(c)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		if err = c.ShouldBind(req); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		resp.Format(f.Search(header.MutateContext(c), req)).Context(c)
	}
}

func Create(f form.Form) gin.HandlerFunc {
	return func(c *gin.Context) {
		req := &form.CreateReq{}
		p := getProfile(c)
		req.UserID = p.userID
		req.DepID = p.depID
		var err error
		req.AppID, req.TableID, err = checkURL(c)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		if err = c.ShouldBind(req); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		resp.Format(f.Create(header.MutateContext(c), req)).Context(c)
	}
}

func Get(f form.Form) gin.HandlerFunc {
	return func(c *gin.Context) {
		req := &form.GetReq{}
		p := getProfile(c)
		req.UserID = p.userID
		req.DepID = p.depID
		var err error
		req.AppID, req.TableID, err = checkURL(c)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		if err = c.ShouldBind(req); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		resp.Format(f.Get(header.MutateContext(c), req)).Context(c)
	}
}

func Update(f form.Form) gin.HandlerFunc {
	return func(c *gin.Context) {
		req := &form.UpdateReq{}
		p := getProfile(c)
		req.UserID = p.userID
		req.DepID = p.depID
		var err error
		req.AppID, req.TableID, err = checkURL(c)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		if err = c.ShouldBind(req); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		resp.Format(f.Update(header.MutateContext(c), req)).Context(c)
	}
}

func Delete(f form.Form) gin.HandlerFunc {
	return func(c *gin.Context) {
		req := &form.DeleteReq{}
		p := getProfile(c)
		req.UserID = p.userID
		req.DepID = p.depID
		var err error
		req.AppID, req.TableID, err = checkURL(c)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		if err = c.ShouldBind(req); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		resp.Format(f.Delete(header.MutateContext(c), req)).Context(c)
	}
}

type profile struct {
	userID   string
	depID    string
	userName string
}

func Action(p *form.Poly) gin.HandlerFunc {
	return func(c *gin.Context) {
		req := &form.ProxyReq{}
		p2 := getProfile(c)
		req.DepID = p2.depID
		req.UserID = p2.userID
		var err error
		req.AppID, req.TableID, err = checkURL(c)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		if err = c.ShouldBind(req); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		resp.Format(p.Proxy(header.MutateContext(c), req)).Context(c)
	}
}

func getProfile(c *gin.Context) *profile {
	userID := c.GetHeader(_userID)
	userName := c.GetHeader(_userName)
	depIDS := strings.Split(c.GetHeader(_departmentID), ",")
	return &profile{
		userID:   userID,
		userName: userName,
		depID:    depIDS[len(depIDS)-1],
	}
}
