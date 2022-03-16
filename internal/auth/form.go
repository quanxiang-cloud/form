package auth

import (
	"context"
	"fmt"
	"net/http"
	"time"

	error2 "github.com/quanxiang-cloud/cabin/error"
	"github.com/quanxiang-cloud/cabin/logger"
	redis2 "github.com/quanxiang-cloud/cabin/tailormade/db/redis"
	"github.com/quanxiang-cloud/cabin/tailormade/header"
	"github.com/quanxiang-cloud/form/internal/auth/cache"
	"github.com/quanxiang-cloud/form/internal/auth/lowcode"
	"github.com/quanxiang-cloud/form/internal/models"
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

	Filter(*http.Response) error
}

type formAuth struct {
	redis   *cache.LimitRepo
	lowcode *lowcode.Lowcode
}

func NewFormAuth(conf *config.Config) (FormAuth, error) {
	redisClient, err := redis2.NewClient(conf.Redis)
	if err != nil {
		return nil, err
	}

	return &formAuth{
		redis:   cache.NewLimitRepo(redisClient),
		lowcode: lowcode.NewLowcode(),
	}, nil
}

type FormAuthReq struct {
	AppID   string `json:"appID,omitempty"`
	TableID string `json:"tableID,omitempty"`
	Path    string `json:"path,omitempty"`
	UserID  string `json:"userID,omitempty"`
	DepID   string `json:"depID,omitempty"`
}

type FormAuthResp struct {
	IsPermit bool `json:"isPermit,omitempty"`
}

func (f *formAuth) Auth(ctx context.Context, req *FormAuthReq) (*FormAuthResp, error) {
	// 从redis获取权限

	return nil, nil
}

func (f *formAuth) Filter(resp *http.Response) error {
	return nil
}

func (f *formAuth) getPermit(ctx context.Context, req *FormAuthReq) {
	match, err := f.getCacheMatch(ctx, req)
	if err != nil {
		// return nil, err
	}
	if match == nil {
		// return nil, error2.New(code.ErrNotPermit)
	}

	if match.Types == models.InitType {
		// return &GetPerInCacheResp{
		// 	Types: match.Types,
		// }, nil
	}
	permits, err := f.getCachePermit(ctx, match.RoleID, req)
	if err != nil {
		// return nil, err
	}

	fmt.Println(permits)
	// return &GetPerInCacheResp{
	// 	Params:    permits.Params,
	// 	Response:  permits.Response,
	// 	Condition: permits.Condition,
	// 	Types:     match.Types,
	// }, err
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

	f.lowcode.GetCacheMatchRole()
	return nil, nil
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

	f.lowcode.GetRoleMatchPermit()

	return nil, nil
}
