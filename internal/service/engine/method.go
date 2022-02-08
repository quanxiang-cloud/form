package engine

import (
	"context"
	error2 "github.com/quanxiang-cloud/cabin/error"
	"github.com/quanxiang-cloud/cabin/logger"
	comet2 "github.com/quanxiang-cloud/form/internal/comet"
	"github.com/quanxiang-cloud/form/internal/filters"
	"github.com/quanxiang-cloud/form/internal/models"
	"github.com/quanxiang-cloud/form/internal/service/types"
	"github.com/quanxiang-cloud/form/pkg/client"
	"github.com/quanxiang-cloud/form/pkg/misc/code"
)

type comet struct {
	formClient *client.FormAPI
}

const notAuthority = -1

type optionEntity interface{}

type bus struct {
	userID    string
	depID     string
	tableName string
	AppID     string
	entity    comet2.Entity
	permit    *Permit
}

// SearchReq SearchReq
type SearchReq struct {
	AppID   string
	TableID string
	Page    int64
	Size    int64
	Sort    []string
	Query   types.Query
	Aggs    types.Any
	UserID  string
	DepID   string
	Permit  *Permit
}

type Permit struct {
	Filter      map[string]interface{}
	Condition   map[string]models.Query
	Authority   int64
	permitTypes models.PerType
}

func SearchData(ctx context.Context, req *SearchReq, opts ...optionEntity) (*client.SearchResp, error) {
	var (
		searchResp *client.SearchResp
		err        error
		b          *bus
	)
	b = &bus{
		AppID:     req.AppID,
		userID:    req.UserID,
		depID:     req.DepID,
		tableName: req.TableID,
		permit:    req.Permit,
	}

	for _, option := range opts {
		err = preStep(ctx, b, "find", option, CheckOperate())
		if err != nil {
			return nil, err // 返回错误
		}
		if handle, ok := option.(Search); ok {
			searchResp, err = handle.Search(ctx, req)
			if err != nil {
				return nil, err
			}

		}

		if _, ok := option.(Pre); ok {
			post(ctx, searchResp.Entities, b, option, JSONFilter())
		}
	}
	return searchResp, nil
}

func (c1 *comet) Search(ctx context.Context, req *SearchReq) (*client.SearchResp, error) {
	dsl := make(map[string]interface{})
	if req.Aggs != nil {
		dsl["aggs"] = req.Aggs
	}
	if req.Query != nil {
		dsl["query"] = req.Query
	}

	if len(dsl) == 0 {
		dsl = nil
	}
	searchResp, err := c1.formClient.Search(ctx, client.FindOptions{
		Sort: req.Sort,
		Size: req.Size,
		Page: req.Page,
	}, dsl, req.TableID)
	if err != nil {
		return nil, err
	}
	return searchResp, nil
}

func post(ctx context.Context, data interface{}, b *bus, option optionEntity, opts ...FilterOption) {
	if post, ok := option.(Post); ok {
		err := post.Postfix(ctx, data, b, opts...)
		if err != nil {
			logger.Logger.Errorw("msg", err.Error())
		}

	}
}

// Postfix Postfix
func (c1 *comet) Postfix(ctx context.Context, data interface{}, bus *bus, opt ...FilterOption) error {
	if bus.permit.permitTypes == models.InitType {
		return nil
	}
	for _, opt := range opt {
		opt(data, bus.permit.Filter)
	}
	return nil
}

func preStep(ctx context.Context, b *bus, method string, option optionEntity, opts ...PreOption) error {
	if pre, ok := option.(Pre); ok {
		err := pre.Pre(ctx, b, method, opts...)
		if err != nil {
			return err
		}
	}
	return nil
}

// JSONFilter JSONFilter
func JSONFilter() FilterOption {
	return func(data interface{}, filter map[string]interface{}) error {
		if data == nil {
			return nil
		}

		if filter == nil {
			return error2.New(code.ErrNotPermit)
		}
		filters.JSONFilter2(data, filter)
		return nil
	}
}
