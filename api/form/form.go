package api

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/quanxiang-cloud/cabin/logger"
	"github.com/quanxiang-cloud/cabin/tailormade/header"
	"github.com/quanxiang-cloud/cabin/tailormade/resp"
	"github.com/quanxiang-cloud/form/internal/service/consensus"
	"github.com/quanxiang-cloud/form/internal/service/types"
)

type profile struct {
	userID   string
	depID    string
	userName string
}

func action(ctr consensus.Guidance) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := header.MutateContext(c)

		bus := &consensus.Bus{}
		err := initBus(c, bus, c.Param("action"))
		if err != nil {
			logger.Logger.WithName("action").Errorw(err.Error(), header.GetRequestIDKV(ctx).Fuzzy()...)
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		if err = c.ShouldBind(bus); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		do, err := ctr.Do(ctx, bus)

		resp.Format(do, err).Context(c)
	}
}

// checkURL CheckURL.
func checkURL(c *gin.Context) (appID, tableName string, err error) {
	appID, ok := c.Params.Get("appID")
	tableName, okt := c.Params.Get("tableName")
	if !ok || !okt {
		err = errors.New("invalid URI")
		return
	}
	return
}

func get(ctr consensus.Guidance) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := header.MutateContext(c)

		bus := &consensus.Bus{}
		err := initBus(c, bus, "get")
		if err != nil {
			logger.Logger.WithName("get").Errorw(err.Error(), header.GetRequestIDKV(ctx).Fuzzy()...)
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		bus.Get.Query = types.Query{
			"term": types.M{
				"_id": c.Param("id"),
			},
		}
		if err = c.ShouldBind(bus); err != nil {
			logger.Logger.WithName("get").Errorw(err.Error(), header.GetRequestIDKV(ctx).Fuzzy()...)
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		do, err := ctr.Do(ctx, bus)

		resp.Format(do, err).Context(c)
	}
}

func search(ctr consensus.Guidance) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := header.MutateContext(c)

		bus := &consensus.Bus{}
		err := initBus(c, bus, "search")
		if err != nil {
			logger.Logger.WithName("search").Errorw(err.Error(), header.GetRequestIDKV(ctx).Fuzzy()...)
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		if err = c.ShouldBind(bus); err != nil {
			logger.Logger.WithName("search").Errorw(err.Error(), header.GetRequestIDKV(ctx).Fuzzy()...)
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		do, err := ctr.Do(header.MutateContext(c), bus)
		resp.Format(do, err).Context(c)
	}
}

func delete(ctr consensus.Guidance) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := header.MutateContext(c)
		bus := &consensus.Bus{}
		err := initBus(c, bus, "delete")
		if err != nil {
			logger.Logger.WithName("delete").Errorw(err.Error(), header.GetRequestIDKV(ctx).Fuzzy()...)
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		bus.Get.Query = types.Query{
			"term": types.M{
				"_id": c.Param("id"),
			},
		}
		if err = c.ShouldBind(bus); err != nil {
			logger.Logger.WithName("delete").Errorw(err.Error(), header.GetRequestIDKV(ctx).Fuzzy()...)
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		do, err := ctr.Do(header.MutateContext(c), bus)

		resp.Format(do, err).Context(c)
	}
}

func update(ctr consensus.Guidance) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := header.MutateContext(c)
		bus := &consensus.Bus{}
		err := initBus(c, bus, "update")
		if err != nil {
			logger.Logger.WithName("update").Errorw(err.Error(), header.GetRequestIDKV(ctx).Fuzzy()...)
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		bus.Get.Query = types.Query{
			"term": types.M{
				"_id": c.Param("id"),
			},
		}
		if err = c.ShouldBind(bus); err != nil {
			logger.Logger.WithName("update").Errorw(err.Error(), header.GetRequestIDKV(ctx).Fuzzy()...)
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		do, err := ctr.Do(header.MutateContext(c), bus)

		resp.Format(do, err).Context(c)
	}
}

func create(ctr consensus.Guidance) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := header.MutateContext(c)
		bus := &consensus.Bus{}
		err := initBus(c, bus, "create")
		if err != nil {
			logger.Logger.WithName("create").Errorw(err.Error(), header.GetRequestIDKV(ctx).Fuzzy()...)
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		if err = c.ShouldBind(bus); err != nil {
			logger.Logger.WithName("create").Errorw(err.Error(), header.GetRequestIDKV(ctx).Fuzzy()...)
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		do, err := ctr.Do(header.MutateContext(c), bus)

		resp.Format(do, err).Context(c)
	}
}

func initBus(c *gin.Context, bus *consensus.Bus, method string) error {
	var err error
	bus.AppID, bus.TableID, err = checkURL(c)
	if err != nil {
		return errors.New("bad path")
	}
	bus.Method = method
	bus.UserID = c.GetHeader(_userID)
	bus.UserName = c.GetHeader(_userName)
	depIDS := strings.Split(c.GetHeader(_departmentID), ",")
	bus.DepID = depIDS[0]
	bus.Path = c.Request.RequestURI
	return nil
}

type relationReq struct {
	TableID    string      `json:"tableID"`
	SubTableID string      `json:"subTableID"`
	FieldKey   string      `json:"fieldKey"`
	Page       int         `json:"page"`
	Size       int         `json:"size"`
	Query      types.Query `json:"query"`
}

func relation(ctr consensus.Guidance) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := header.MutateContext(c)
		bus := &consensus.Bus{}
		err := initBus(c, bus, "create")
		if err != nil {
			logger.Logger.WithName("create").Errorw(err.Error(), header.GetRequestIDKV(ctx).Fuzzy()...)
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		if err = c.ShouldBind(bus); err != nil {
			logger.Logger.WithName("create").Errorw(err.Error(), header.GetRequestIDKV(ctx).Fuzzy()...)
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		do, err := ctr.Do(header.MutateContext(c), bus)

		resp.Format(do, err).Context(c)
	}
}
