package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	error2 "github.com/quanxiang-cloud/cabin/error"
	"github.com/quanxiang-cloud/cabin/logger"
	"github.com/quanxiang-cloud/cabin/tailormade/client"
	redis2 "github.com/quanxiang-cloud/cabin/tailormade/db/redis"
	"github.com/quanxiang-cloud/cabin/tailormade/header"
	"github.com/quanxiang-cloud/form/internal/auth/filters"
	"github.com/quanxiang-cloud/form/internal/auth/lowcode"
	"github.com/quanxiang-cloud/form/internal/models"
	"github.com/quanxiang-cloud/form/internal/models/redis"
	"github.com/quanxiang-cloud/form/internal/service/consensus"
	"github.com/quanxiang-cloud/form/pkg/misc/code"
	"github.com/quanxiang-cloud/form/pkg/misc/config"
)

const (
	lockPermission = "lockPermission"
	lockPerMatch   = "lockPerMatch"
	lockTimeout    = time.Duration(30) * time.Second // 30秒
	timeSleep      = time.Millisecond * 500          // 0.5 秒
)

type FormAuth interface {
	Auth(context.Context, *FormAuthReq) (*FormAuthResp, error)

	Filter(*http.Response, string) error
}

type formAuth struct {
	redis   models.LimitsRepo
	lowcode lowcode.Form
	permit  *consensus.Permit
}

func NewFormAuth(conf *config.Config) (FormAuth, error) {
	redisClient, err := redis2.NewClient(conf.Redis)
	if err != nil {
		return nil, err
	}

	return &formAuth{
		redis:   redis.NewLimitRepo(redisClient),
		lowcode: lowcode.NewForm(client.Config{}),
	}, nil
}

type FormAuthReq struct {
	AppID  string      `json:"appID,omitempty"`
	UserID string      `json:"userID,omitempty"`
	DepID  string      `json:"depID,omitempty"`
	Path   string      `json:"path,omitempty"`
	Entity interface{} `json:"entity,omitempty"`
}

type FormAuthResp struct {
	IsPermit bool `json:"isPermit,omitempty"`
}

func (f *formAuth) Auth(ctx context.Context, req *FormAuthReq) (*FormAuthResp, error) {
	// get the role information owned by the user
	match, err := f.getCacheMatch(ctx, req)
	if err != nil || match == nil {
		return &FormAuthResp{}, err
	}

	if match.Types == models.InitType {
		return &FormAuthResp{true}, nil
	}

	permits, err := f.getCachePermit(ctx, match.RoleID, req)
	if err != nil {
		return &FormAuthResp{}, err
	}

	// access judgment
	if !filters.Pre(req.Entity, permits.Params) {
		return &FormAuthResp{}, error2.New(code.ErrNotPermit)
	}

	f.permit = &consensus.Permit{
		Params:    permits.Params,
		Response:  permits.Response,
		Condition: permits.Condition,
		Types:     match.Types,
	}

	return &FormAuthResp{true}, nil
}

func (f *formAuth) Filter(resp *http.Response, method string) error {
	respDate, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	conResp := &consensus.Response{}

	err = json.Unmarshal(respDate, conResp)
	if err != nil {
		return err
	}

	var entity interface{}
	switch method {
	case "get":
		entity = conResp.GetResp.Entity
	case "search":
		entity = conResp.ListResp.Entities
	}
	filters.Post(entity, f.permit.Response)

	data, err := json.Marshal(entity)
	if err != nil {
		logger.Logger.Errorf("entity json marshal failed: %s", err.Error())
		return err
	}

	resp.Body = io.NopCloser(bytes.NewReader(data))
	return nil
}

func (f *formAuth) getCacheMatch(ctx context.Context, req *FormAuthReq) (*models.PermitMatch, error) {
	// relese lock
	defer f.redis.UnLock(ctx, lockPerMatch)
	for i := 0; i < 5; i++ {
		perMatch, err := f.redis.GetPerMatch(ctx, req.UserID, req.AppID)
		if err != nil {
			logger.Logger.Errorw(req.UserID, header.GetRequestIDKV(ctx).Fuzzy(), err.Error())
			return nil, err
		}
		if perMatch != nil {
			return perMatch, nil
		}

		// acquire distributed locks
		lock, err := f.redis.Lock(ctx, lockPerMatch, 1, lockTimeout)
		if err != nil {
			return nil, err
		}
		if !lock {
			<-time.After(timeSleep)
			continue
		}
		break
	}

	resp, err := f.lowcode.GetCacheMatchRole(ctx, req.UserID, req.DepID, req.AppID)
	if err != nil || resp == nil {
		return nil, err
	}

	return &models.PermitMatch{
		RoleID: resp.RoleID,
		Types:  models.RoleType(resp.Types),
	}, nil
}

func (f *formAuth) getCachePermit(ctx context.Context, roleID string, req *FormAuthReq) (*models.Limits, error) {
	// relese lock
	defer f.redis.UnLock(ctx, lockPermission)
	for i := 0; i < 5; i++ {
		exist := f.redis.ExistsKey(ctx, roleID)
		if exist {
			// judge path
			getPermit, err := f.redis.GetPermit(ctx, roleID, req.Path)
			if err != nil {
				return nil, err
			}
			if getPermit.Path == "" {
				return nil, error2.New(code.ErrNotPermit)
			}
			return getPermit, nil
		}

		// acquire distributed locks
		lock, err := f.redis.Lock(ctx, lockPermission, 1, lockTimeout)
		if err != nil {
			return nil, err
		}
		if !lock {
			<-time.After(timeSleep)
			continue
		}
		break
	}

	resp, err := f.lowcode.GetRoleMatchPermit(ctx, roleID)
	if err != nil || resp == nil {
		return nil, err
	}

	limits := make([]*models.Limits, len(resp.List))
	var getPermit *models.Limits
	for index, value := range resp.List {
		per := &models.Limits{
			Path:      value.Path,
			Condition: value.Condition,
			Params:    value.Params,
			Response:  value.Response,
		}
		if value.Path == req.Path {
			getPermit = per
		}
		limits[index] = per
	}
	err = f.redis.CreatePermit(ctx, roleID, limits...)
	if err != nil {
		logger.Logger.Errorw("create permit err", roleID, err.Error())
	}

	if getPermit == nil {
		return nil, error2.New(code.ErrNotPermit)
	}

	return getPermit, nil
}
