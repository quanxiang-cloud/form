package api

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/quanxiang-cloud/cabin/tailormade/header"
	"github.com/quanxiang-cloud/cabin/tailormade/resp"
	"github.com/quanxiang-cloud/form/internal/service/consensus"
	"github.com/quanxiang-cloud/form/internal/service/form"
	"github.com/quanxiang-cloud/form/internal/service/guidance"
	"github.com/quanxiang-cloud/form/internal/service/types"
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

func action(ctr guidance.Guidance) gin.HandlerFunc {
	return func(c *gin.Context) {
		bus := &consensus.Bus{}
		bus.UserID = c.GetHeader(_userID)
		bus.UserName = c.GetHeader(_userName)

		bus.Method = c.Param("action")

		var err error
		bus.AppID, bus.TableID, err = checkURL(c)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		if err = c.ShouldBind(bus); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		resp.Format(ctr.Do(header.MutateContext(c), bus)).Context(c)
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

func get(ctr guidance.Guidance) gin.HandlerFunc {
	return func(c *gin.Context) {
		bus := &consensus.Bus{}
		bus.UserID = c.GetHeader(_userID)
		bus.Method = "get"

		var err error
		bus.AppID, bus.TableID, err = checkURL(c)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		id := c.Param("id")
		if id == "" {
			c.AbortWithError(http.StatusBadRequest, fmt.Errorf("id is must"))
			return
		}

		bus.Query = types.Query{
			"term": map[string]interface{}{
				"_id": id,
			},
		}

		resp.Format(ctr.Do(header.MutateContext(c), bus)).Context(c)
	}
}

func create(ctr guidance.Guidance) gin.HandlerFunc {
	return func(c *gin.Context) {
		bus := &consensus.Bus{}
		bus.UserID = c.GetHeader(_userID)
		bus.Method = "create"

		var err error
		bus.AppID, bus.TableID, err = checkURL(c)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		if err := c.ShouldBind(bus); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		resp.Format(ctr.Do(header.MutateContext(c), bus)).Context(c)
	}
}

func update(ctr guidance.Guidance) gin.HandlerFunc {
	return func(c *gin.Context) {
		bus := &consensus.Bus{}
		bus.UserID = c.GetHeader(_userID)
		bus.Method = "update"

		var err error
		bus.AppID, bus.TableID, err = checkURL(c)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		id := c.Param("id")
		if id == "" {
			c.AbortWithError(http.StatusBadRequest, fmt.Errorf("id is must"))
			return
		}
		bus.Query = types.Query{
			"term": map[string]interface{}{
				"_id": id,
			},
		}

		if err := c.ShouldBind(bus); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		resp.Format(ctr.Do(header.MutateContext(c), bus)).Context(c)
	}
}

func delete(ctr guidance.Guidance) gin.HandlerFunc {
	return func(c *gin.Context) {
		bus := &consensus.Bus{}
		bus.UserID = c.GetHeader(_userID)
		bus.Method = "delete"

		var err error
		bus.AppID, bus.TableID, err = checkURL(c)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		id := c.Param("id")
		if id == "" {
			c.AbortWithError(http.StatusBadRequest, fmt.Errorf("id is must"))
			return
		}

		bus.Query = types.Query{
			"term": map[string]interface{}{
				"_id": id,
			},
		}

		resp.Format(ctr.Do(header.MutateContext(c), bus)).Context(c)
	}
}
