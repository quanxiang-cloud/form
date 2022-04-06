package permit

import (
	"context"
	redis2 "github.com/quanxiang-cloud/cabin/tailormade/db/redis"
	"github.com/quanxiang-cloud/form/internal/models"
	"github.com/quanxiang-cloud/form/internal/models/redis"
	"github.com/quanxiang-cloud/form/pkg/misc/config"
)

const (
	delete = "delete"
)

type Cache interface {
	UserMatch(ctx context.Context, req *UserMatchReq) (*UserMatchResp, error)
	Limit(ctx context.Context, req *LimitReq) (*LimitResp, error)
}
type cache struct {
	redis models.LimitsRepo
}

type LimitReq struct {
	RoleID    string             `json:"roleID"`
	Path      string             `json:"path"`
	Condition models.Condition   `json:"condition"`
	Params    models.FiledPermit `json:"params"`
	Response  models.FiledPermit `json:"response"`
	Action    string             `json:"action"`
}

type LimitResp struct {
}

func (c *cache) Limit(ctx context.Context, req *LimitReq) (*LimitResp, error) {
	//
	if req.Action == delete {
		err := c.redis.DeletePermit(ctx, req.RoleID)
		if err != nil {
			return nil, err
		}
		return &LimitResp{}, err
	}

	err := c.redis.CreatePermit(ctx, req.RoleID, &models.Limits{
		Path:      req.Path,
		Condition: req.Condition,
		Params:    req.Params,
		Response:  req.Response,
	})
	if err != nil {
		return nil, err
	}
	return &LimitResp{}, nil
}

func NewCache(config *config.Config) (Cache, error) {
	redisClient, err := redis2.NewClient(config.Redis)
	if err != nil {
		return nil, err
	}
	return &cache{
		redis: redis.NewLimitRepo(redisClient),
	}, nil
}

type UserMatchReq struct {
	AppID  string `json:"appID"`
	UserID string `json:"userID"`
	RoleID string `json:"roleID"`
	Action string `json:"action"`
}

type UserMatchResp struct {
}

func (c *cache) UserMatch(ctx context.Context, req *UserMatchReq) (*UserMatchResp, error) {
	if req.Action == "delete" {
		err := c.redis.DeletePerMatch(ctx, req.AppID)
		if err != nil {
			return nil, err
		}
		return &UserMatchResp{}, nil
	}
	// create
	err := c.redis.CreatePerMatch(ctx, &models.PermitMatch{
		RoleID: req.RoleID,
		UserID: req.UserID,
		AppID:  req.AppID,
	})
	if err != nil {
		return nil, err
	}
	return &UserMatchResp{}, nil

}
