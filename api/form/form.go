package api

import (
	"errors"
	"fmt"
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

type batchCreateResp struct {
	Entity []consensus.Entity `json:"entity"`
	Total  int                `json:"total"`
}

func batchCreate(ctr consensus.Guidance) gin.HandlerFunc {
	return func(c *gin.Context) {
		var err error
		ctx := header.MutateContext(c)
		var batch []*consensus.Bus
		if err = c.ShouldBind(&batch); err != nil {
			logger.Logger.WithName("action").Errorw(err.Error(), header.GetRequestIDKV(ctx).Fuzzy()...)
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		total := 0
		entitys := make([]consensus.Entity, 0)
		for _, bus := range batch {
			err = initBus(c, bus, c.Param("action"))
			if err != nil {
				continue
			}
			do, errs := ctr.Do(ctx, bus)
			if errs != nil {
				continue
			}
			total++
			entitys = append(entitys, do.Entity)
		}
		resp1 := &batchCreateResp{
			Entity: entitys,
			Total:  total,
		}
		resp.Format(resp1, nil).Context(c)
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
		if bus.Sub.PID == "" { // is  normal
			do, err := ctr.Do(header.MutateContext(c), bus)
			resp.Format(do, err).Context(c)
			return
		}
		ids := consensus.GetSimple(consensus.TermKey, "primitiveID", bus.Sub.PID)
		keys := consensus.GetSimple(consensus.TermKey, "fieldName", bus.Sub.FieldKey)
		boolQuery := consensus.GetBool(consensus.Must, ids, keys)
		bus1 := getBus(bus.AppID, getRelationName(bus.PTableID, bus.TableID), boolQuery, 1, 300)
		searchResp1, err := ctr.Do(ctx, bus1)
		if err != nil {
			logger.Logger.WithName("search err").Errorw(err.Error(), header.GetRequestIDKV(ctx).Fuzzy()...)
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		data := make([]interface{}, 0)
		for _, value := range searchResp1.Entities {
			_, ok := value["subID"]
			if !ok {
				continue
			}
			data = append(data, value["subID"])
		}
		subQuery := consensus.GetSimple(consensus.TermsKey, "_id", data)
		if len(bus.Get.Query) != 0 {
			subQuery = consensus.GetBool(consensus.Must, subQuery, bus.Get.Query)
		}
		bus2 := getBus(bus.AppID, bus.TableID, subQuery, bus.Page, bus.Size)
		resp.Format(ctr.Do(ctx, bus2)).Context(c)
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
	ID         string      `json:"id" form:"id"`
	TableID    string      `json:"tableID"`
	SubTableID string      `json:"subTableID" form:"subTableID"`
	FieldKey   string      `json:"fieldKey" form:"fieldKey"`
	Page       int64       `json:"page" form:"page"`
	Size       int64       `json:"size" form:"size"`
	Query      types.Query `json:"query" form:"query"`
	AppID      string      `json:"appID"`
}

func relation(ctr consensus.Guidance) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := header.MutateContext(c)
		req := &relationReq{}
		req.AppID = c.Param("appID")
		req.TableID = c.Param("tableName")
		if err := c.ShouldBind(req); err != nil {
			logger.Logger.WithName("relation").Errorw(err.Error(), header.GetRequestIDKV(ctx).Fuzzy()...)
			return
		}
		ids := consensus.GetSimple(consensus.TermKey, "primitiveID", req.ID)
		keys := consensus.GetSimple(consensus.TermKey, "fieldName", req.FieldKey)
		boolQuery := consensus.GetBool(consensus.Must, ids, keys)
		bus := getBus(req.AppID, getRelationName(req.TableID, req.SubTableID), boolQuery, 1, 300)
		searchResp1, err := ctr.Do(ctx, bus)
		if err != nil {
			logger.Logger.WithName("search err").Errorw(err.Error(), header.GetRequestIDKV(ctx).Fuzzy()...)
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		data := make([]interface{}, 0)
		for _, value := range searchResp1.Entities {
			_, ok := value["subID"]
			if !ok {
				continue
			}
			data = append(data, value["subID"])
		}
		subQuery := consensus.GetSimple(consensus.TermsKey, "_id", data)
		if len(req.Query) != 0 {
			subQuery = consensus.GetBool(consensus.Must, subQuery, req.Query)
		}
		bus1 := getBus(req.AppID, req.SubTableID, subQuery, req.Page, req.Size)
		resp.Format(ctr.Do(ctx, bus1)).Context(c)
	}
}

func getBus(appID, tableID string, query types.Query, page, size int64) *consensus.Bus {
	bus := new(consensus.Bus)
	bus.Foundation = consensus.Foundation{
		AppID:   appID,
		TableID: tableID,
		Method:  "search",
	}
	bus.Get.Query = query
	bus.List = consensus.List{
		Size: size,
		Page: page,
		Sort: []string{"created_at"},
	}
	return bus
}

func getRelationName(primary, sub string) string {
	return fmt.Sprintf("%s_%s", primary, sub)
}
