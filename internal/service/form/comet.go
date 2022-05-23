package form

import (
	"context"
	"fmt"
	"reflect"

	"github.com/quanxiang-cloud/form/internal/service/types"
	"github.com/quanxiang-cloud/form/pkg/misc/config"

	"github.com/quanxiang-cloud/form/internal/service/consensus"
	client2 "github.com/quanxiang-cloud/form/pkg/misc/client"
)

type comet struct {
	formClient *client2.FormAPI
}

func newForm(config *config.Config) (consensus.Guidance, error) {
	formApi, err := client2.NewFormAPI(config)
	if err != nil {
		return nil, err
	}
	return &comet{
		formClient: formApi,
	}, nil
}

func (c *comet) Do(ctx context.Context, bus *consensus.Bus) (*consensus.Response, error) {
	// TODO
	base := Base{
		AppID:   bus.AppID,
		TableID: bus.TableID,
		UserID:  bus.UserID,
	}
	switch bus.Foundation.Method {
	case "get":
		req := &GetReq{
			Base:  base,
			Query: bus.Query,
		}
		req.Base = base
		req.Query = bus.Query
		req.Aggs = bus.Aggs
		return c.callGet(ctx, req)

	case "find", "search":
		req := &SearchReq{
			Sort:  bus.List.Sort,
			Page:  bus.List.Page,
			Size:  bus.List.Size,
			Query: bus.Query,
			Base:  base,
			Aggs:  bus.Aggs,
		}
		return c.callSearch(ctx, req)
	case "create":
		req := &CreateReq{
			Entity: bus.CreatedOrUpdate.Entity,
			Base:   base,
		}
		return c.callCreate(ctx, req)
	case "update":
		req := &UpdateReq{
			Entity: bus.CreatedOrUpdate.Entity,
			Query:  bus.Query,
			Base:   base,
		}
		return c.callUpdate(ctx, req)
	case "delete":
		req := &DeleteReq{
			Query: bus.Query,
			Base:  base,
		}
		return c.callDelete(ctx, req)
	}
	return nil, nil
}

func (c *comet) callSearch(ctx context.Context, req *SearchReq) (*consensus.Response, error) {
	dsl := make(map[string]interface{})
	if req.Query != nil {
		dsl["query"] = req.Query
	}
	if len(dsl) == 0 {
		dsl = nil
	}
	if req.Aggs != nil {
		dsl["aggs"] = req.Aggs
	}
	formReq := &client2.FormReq{
		DslQuery: dsl,
	}
	formReq.Size = req.Size
	formReq.Page = req.Page
	formReq.Sort = req.Sort
	formReq.TableID = getTableID(req.AppID, req.TableID)

	searchResp, err := c.formClient.Search(ctx, formReq)
	if err != nil {
		return nil, err
	}
	resp := new(consensus.Response)
	resp.Total = searchResp.Total
	resp.Entities = searchResp.Entities
	return resp, nil
}

func (c *comet) callCreate(ctx context.Context, req *CreateReq) (*consensus.Response, error) {
	formReq := &client2.FormReq{
		Entity:  req.Entity,
		TableID: getTableID(req.AppID, req.TableID),
	}
	insert, err := c.formClient.Insert(ctx, formReq)
	if err != nil {
		return nil, err
	}
	resp := new(consensus.Response)
	resp.Total = insert.SuccessCount
	resp.Entity = req.Entity
	return resp, nil
}

func get(e consensus.Entity) types.Entity {
	if e == nil {
		return nil
	}
	value := reflect.ValueOf(e)
	switch _t := reflect.TypeOf(e); _t.Kind() {
	case reflect.Map:
		iter := value.MapRange()
		m := make(types.Entity)
		for iter.Next() {
			if !iter.Value().CanInterface() {
				continue
			}
			m[iter.Key().String()] = iter.Value()
		}
		return m
	default:
		return nil
	}
	return nil
}

func (c *comet) callUpdate(ctx context.Context, req *UpdateReq) (*consensus.Response, error) {
	dsl := make(map[string]interface{})
	if req.Query != nil {
		dsl["query"] = req.Query
	}
	if len(dsl) == 0 {
		dsl = nil
	}

	formReq := &client2.FormReq{
		Entity:   req.Entity,
		TableID:  getTableID(req.AppID, req.TableID),
		DslQuery: dsl,
	}
	updates, err := c.formClient.Update(ctx, formReq)
	if err != nil {
		return nil, err
	}
	resp := &consensus.Response{}
	resp.Total = updates.SuccessCount
	resp.Entity = req.Entity
	return resp, nil
}

func getTableID(appID, tableID string) string {
	if len(appID) == 36 {
		return fmt.Sprintf("%s%s%s", "A", appID, tableID)
	}
	return fmt.Sprintf("%s%s%s", "a", appID, tableID)
}

func (c *comet) callGet(ctx context.Context, req *GetReq) (*consensus.Response, error) {
	dsl := make(map[string]interface{})
	if req.Query != nil {
		dsl["query"] = req.Query
	}
	if req.Aggs != nil {
		dsl["aggs"] = req.Aggs
	}
	if len(dsl) == 0 {
		dsl = nil
	}

	formReq := &client2.FormReq{
		DslQuery: dsl,
		TableID:  getTableID(req.AppID, req.TableID),
	}
	gets, err := c.formClient.Get(ctx, formReq)
	if err != nil {
		return nil, err
	}
	resp := &consensus.Response{}
	resp.Entity = gets.Entity
	return resp, nil
}

func (c *comet) callDelete(ctx context.Context, req *DeleteReq) (*consensus.Response, error) {
	dsl := make(map[string]interface{})
	if req.Query != nil {
		dsl["query"] = req.Query
	}
	if len(dsl) == 0 {
		dsl = nil
	}
	formReq := &client2.FormReq{
		DslQuery: dsl,
		TableID:  getTableID(req.AppID, req.TableID),
	}
	deletes, err := c.formClient.Delete(ctx, formReq)
	if err != nil {
		return nil, err
	}
	resp := &consensus.Response{}
	resp.Total = deletes.SuccessCount
	return resp, nil
}
