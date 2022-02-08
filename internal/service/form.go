package service

import (
	"context"
	"fmt"
	error2 "github.com/quanxiang-cloud/cabin/error"
	"github.com/quanxiang-cloud/form/internal/comet"
	"github.com/quanxiang-cloud/form/internal/service/engine"
	"github.com/quanxiang-cloud/form/internal/service/types"
	"github.com/quanxiang-cloud/form/pkg/client"
	"github.com/quanxiang-cloud/form/pkg/misc/code"
	config2 "github.com/quanxiang-cloud/form/pkg/misc/config"
)

const (
	put    = "put"
	post   = "post"
	delete = "delete"
	get    = "get"
)

type Form interface {
	Search(ctx context.Context, req *SearchReq) (*SearchResp, error)
	Create(ctx context.Context, req *CreateReq) (*CreateResp, error)
}

type form struct {
	auth    engine.Plugs
	notAuth engine.Plugs
	permit  Permission
}

type CreateReq struct {
	IsAuth   bool         `json:"-"`
	AppID    string       `json:"-"`
	TableID  string       `json:"-"`
	UserName string       `json:"_"`
	UserId   string       `json:"-"`
	Entity   comet.Entity `json:"entity"`
	Ref      types.Ref    `json:"ref"`
}

type CreateResp struct {
}

func (f *form) Create(ctx context.Context, req *CreateReq) (*CreateResp, error) {
	// 鉴权 ，不鉴权

	return nil, nil
}

func (f *form) callCreate(ctx context.Context, req *CreateReq) (*client.InsertResp, error) {
	tableID := getTableID(req.AppID, req.TableID)
	var _ string
	_ = tableID
	if req.IsAuth {
		return nil, nil
	}
	return nil, nil
}

func getTableID(appID, tableID string) string {
	if len(appID) == 36 {
		return fmt.Sprintf("%s%s%s", "A", appID, tableID)
	}
	return fmt.Sprintf("%s%s%s", "a", appID, tableID)
}

func (f *form) Search(ctx context.Context, req *SearchReq) (*SearchResp, error) {
	searchData, err := f.callSearch(ctx, req)
	if err != nil {
		return nil, err
	}
	return &SearchResp{
		Entities:     searchData.Entities,
		Total:        searchData.Total,
		Aggregations: searchData.Aggregations,
	}, nil
}

func (f *form) callSearch(ctx context.Context, req *SearchReq) (*client.SearchResp, error) {
	//searchReq := &engine.SearchReq{
	//	AppID:   req.AppID,
	//	TableID: getTableID(req.AppID, req.TableID),
	//	Page:    req.Page,
	//	Size:    req.Size,
	//	Sort:    req.Sort,
	//	Query:   req.Query,
	//	Aggs:    req.Aggs,
	//}
	//
	//authOption := f.notAuth
	//if req.IsAuth {
	//	authOption = f.auth
	//	permit, err := f.getPermit(ctx, &GetPerInCacheReq{
	//		AppID:  req.AppID,
	//		DepID:  req.DepID,
	//		UserID: req.UserID,
	//		FormID: req.TableID,
	//	})
	//	if err != nil {
	//		return nil, err
	//	}
	//	searchReq.Permit = permit
	//}

	//searchData, err := engine.SearchData(ctx, authOption, searchReq)
	//if err != nil {
	//	return nil, err
	//}
	//return searchData, nil
	return nil, nil
}

func (f *form) getPermit(ctx context.Context, req *GetPerInCacheReq) (*engine.Permit, error) {
	cache, err := f.permit.GetPerInCache(ctx, req)
	if err != nil {
		return nil, err
	}
	if cache == nil {
		return nil, error2.New(code.ErrNotPermit)
	}
	return &engine.Permit{
		Condition: cache.DataAccessPer,
		Authority: cache.Authority,
		Filter:    cache.Filter,
	}, nil
}

func NewForm(conf *config2.Config) (Form, error) {
	auth, err := engine.NewAuth(conf)
	if err != nil {
		return nil, err
	}
	noAuth, err := engine.NewNoAuth(conf)
	if err != nil {
		return nil, err
	}
	permit, err := NewPermission(conf)
	if err != nil {
		return nil, err
	}
	return &form{
		auth:    auth,
		notAuth: noAuth,
		permit:  permit,
	}, nil
}

type FindOptions struct {
	Page int64    `json:"page"`
	Size int64    `json:"size"`
	Sort []string `json:"sort"`
}

type SearchReq struct {
	IsAuth  bool   `json:"-"`
	AppID   string `json:"-"`
	TableID string `json:"-"`
	FindOptions
	Query  types.Query `json:"query"`
	Ref    types.Ref   `json:"ref"`
	Aggs   types.Any   `json:"aggs"`
	UserID string      `json:"userID"`
	DepID  string      `json:"depID"`
}

type SearchResp struct {
	Entities     types.Entities `json:"entities"`
	Total        int64          `json:"total"`
	Aggregations types.Any      `json:"aggregations"`
}
