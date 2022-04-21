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

//idConditions := consensus.GetSimple(consensus.TermsKey, primitiveID, id)
//	keyCondition := consensus.GetSimple(consensus.TermKey, fieldName, c.key)
//	boolQuery := consensus.GetBool(consensus.Must, idConditions, keyCondition)
//
//	universal := consensus.Universal{
//		UserID:   c.userID,
//		UserName: c.userName,
//	}
//	foundation := consensus.Foundation{
//		AppID:   refData.AppID,
//		TableID: getRelationName(extraData.TableID, refData.TableID),
//		Method:  "search",
//	}
//	bus := new(consensus.Bus)
//	bus.Universal = universal
//	bus.Foundation = foundation
//	bus.Get.Query = boolQuery
//	list := consensus.List{
//		Size: 1000,
//		Page: 1,
//		Sort: []string{"created_at"},
//	}
//	bus.List = list
//
//	searchResp1, err := c.ref.Do(ctx, bus)
//	if err != nil {
//		return err
//	}
//	for _, value := range searchResp1.Entities {
//		_, ok := value[subIDs]
//		if !ok {
//			continue
//		}
//		data = append(data, value[subIDs])
//	}
//	if !isReplace {
//		setValue(c.primaryEntity, c.key, data)
//		return nil
//	}
//
//	idsQuery := consensus.GetSimple(consensus.TermsKey, "_id", data)
//	bus1 := new(consensus.Bus)
//	bus1.Universal = universal
//	bus1.Foundation = consensus.Foundation{
//		TableID: refData.TableID,
//		AppID:   refData.AppID,
//		Method:  "search",
//	}
//	bus1.Get.Query = idsQuery
//	bus.List = list
//
//	subResp, err := c.ref.Do(ctx, bus1)
//	if err != nil {
//		return err
//	}
//	err = c.findOnePost(&params{
//		ctx:       ctx,
//		subResp:   subResp,
//		refData:   refData,
//		extraData: extraData,
//		data:      data,
//	})
//	if err != nil {
//		return err
//	}
//	return nil
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
