package lowcode

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	error2 "github.com/quanxiang-cloud/cabin/error"
	"github.com/quanxiang-cloud/form/internal/filters"
	"github.com/quanxiang-cloud/form/internal/models"
	"github.com/quanxiang-cloud/form/internal/service"
	"github.com/quanxiang-cloud/form/internal/service/consensus"
	"github.com/quanxiang-cloud/form/pkg/misc/code"
	"github.com/quanxiang-cloud/form/pkg/misc/config"
)

type FormAuth struct {
	permit     service.Permit
	req        *FormReq
	consPermit *consensus.Permit
}

func NewFormAuth(conf *config.Config, req *FormReq) *FormAuth {
	permit, _ := service.NewPermit(conf)

	return &FormAuth{
		permit: permit,
		req:    req,
	}
}

type FormReq struct {
	UserID   string      `json:"userID,omitempty"`
	DepID    string      `json:"depID,omitempty"`
	UserName string      `json:"userName,omitempty"`
	Method   string      `json:"method,omitempty"`
	Path     string      `json:"path,omitempty"`
	AppID    string      `json:"appID,omitempty"`
	TableID  string      `json:"tableID,omitempty"`
	Entity   interface{} `json:"entity,omitempty"`
}

func (f *FormAuth) Auth(ctx *gin.Context) bool {
	permit, err := f.getPermit(ctx, f.req)
	if err != nil {
		return false
	}

	f.consPermit = permit

	// judge whether there is permission
	if permit.Types == models.InitType {
		return true
	}

	if !filters.Pre(f.req.Entity, permit.Params) {
		return false
	}

	return true
}

func (f *FormAuth) Filter(resp *http.Response) error {
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	res := consensus.Response{}
	err = json.Unmarshal(data, &res)
	if err != nil {
		return err
	}

	var entity interface{}
	switch f.req.Method {
	case "get":
		entity = res.GetResp.Entity
	case "search":
		entity = res.ListResp.Entities
	}
	filters.Post(entity, f.consPermit.Response)

	resp.Body = io.NopCloser(bytes.NewBuffer(data))
	return nil
}

func (f *FormAuth) post(ctx context.Context, bus *consensus.Bus, resp *consensus.Response) {
	var entity interface{}
	switch bus.Method {
	case "get":
		entity = resp.GetResp.Entity
	case "search":
		entity = resp.ListResp.Entities
	}
	filters.Post(entity, bus.Permit.Response)
}

func (f *FormAuth) getPermit(ctx context.Context, req *FormReq) (*consensus.Permit, error) {
	cache, err := f.permit.GetPerInCache(ctx, &service.GetPerInCacheReq{
		UserID: req.UserID,
		DepID:  req.DepID,
		Path:   req.Path,
		AppID:  req.AppID,
	})
	if err != nil {
		return nil, err
	}

	if cache == nil {
		return nil, error2.New(code.ErrNotPermit)
	}

	return &consensus.Permit{
		Params:    cache.Params,
		Response:  cache.Response,
		Condition: cache.Condition,
		Types:     cache.Types,
	}, nil
}
