package treasure

import (
	"context"
	"github.com/quanxiang-cloud/cabin/logger"

	redis2 "github.com/quanxiang-cloud/cabin/tailormade/db/redis"
	"github.com/quanxiang-cloud/cabin/tailormade/header"
	"github.com/quanxiang-cloud/form/internal/models"
	"github.com/quanxiang-cloud/form/internal/models/redis"
	"github.com/quanxiang-cloud/form/internal/permit"
	"github.com/quanxiang-cloud/form/internal/service/consensus"
	"github.com/quanxiang-cloud/form/pkg/misc/client/lowcode"
	"github.com/quanxiang-cloud/form/pkg/misc/config"
)

type Auth struct {
	redis     models.LimitsRepo
	form      *lowcode.Form
	appCenter *lowcode.AppCenter
}

func NewAuth(conf *config.Config) (*Auth, error) {
	redisClient, err := redis2.NewClient(conf.Redis)
	if err != nil {
		return nil, err
	}
	return &Auth{
		redis:     redis.NewLimitRepo(redisClient),
		form:      lowcode.NewForm(conf.InternalNet),
		appCenter: lowcode.NewAppCenter(conf.InternalNet),
	}, nil
}

func (a *Auth) Auth(ctx context.Context, req *permit.Request) (*consensus.Permit, error) {
	if req.UserID == "" {
		logger.Logger.Errorw("userID is blank", header.GetRequestIDKV(ctx).Fuzzy()...)
		return nil, nil
	}
	// 判断 app  是否聚合
	app, err := a.appCenter.GetOne(ctx, req.AppID)
	if err != nil {
		return nil, err
	}
	if app.PerPoly { // 要聚合权限
		// 要聚合权限
		poly, errs := a.form.PerPoly(ctx, req.AppID, req.Path, req.UserID, req.DepID)
		if errs != nil {
			return nil, errs
		}
		return &consensus.Permit{
			Types:       poly.Types,
			Params:      poly.Params,
			Response:    poly.Response,
			Condition:   poly.Condition,
			ParamsAll:   poly.ParamsAll,
			ResponseAll: poly.ResponseAll,
		}, nil
	}

	match, err := a.getUserRole(ctx, req)
	if err != nil || match == nil {
		return nil, err
	}
	if match.RoleID == models.RoleInit {
		return &consensus.Permit{
			Types: models.InitType,
		}, nil
	}
	permits, err := a.getCachePermit(ctx, match.RoleID, req)
	if err != nil || permits == nil {
		return nil, err
	}
	return &consensus.Permit{
		Params:      permits.Params,
		Response:    permits.Response,
		Condition:   permits.Condition,
		ParamsAll:   permits.ParamsAll,
		ResponseAll: permits.ResponseAll,
	}, nil
}

func (a *Auth) getUserRole(ctx context.Context, req *permit.Request) (*models.UserRoles, error) {
	resp, err := a.form.GetCacheMatchRole(ctx, req.UserID, req.DepID, req.AppID)
	if err != nil || resp == nil {
		return nil, err
	}
	if resp.Types == models.InitType {
		resp.RoleID = models.RoleInit
	}
	return &models.UserRoles{
		RoleID: resp.RoleID,
	}, nil
}

func (a *Auth) getCachePermit(ctx context.Context, roleID string, req *permit.Request) (*models.Limits, error) {
	resp, err := a.form.GetPermit(ctx, req.AppID, roleID, req.Path, req.Echo.Request().Method)
	if err != nil || resp == nil {
		return nil, err
	}
	getPermit := &models.Limits{
		Path:        resp.Path,
		Condition:   resp.Condition,
		Params:      resp.Params,
		Response:    resp.Response,
		ParamsAll:   resp.ParamsAll,
		ResponseAll: resp.ResponseAll,
	}
	return getPermit, nil
}
