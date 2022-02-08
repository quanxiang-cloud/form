package service

import (
	"context"
	"github.com/quanxiang-cloud/cabin/logger"
	"github.com/quanxiang-cloud/form/internal/models"
	"time"
)

func (per *permission) getPerMatch(ctx context.Context, userID, depID, appID string) (*models.PermissionMatch, error) {
	for i := 0; i < 5; i++ {
		perMatch, err := per.permissionRepo.GetPerMatch(ctx, userID, appID)
		if err != nil {
			logger.Logger.Errorw("msg", err.Error())
			return nil, err
		}
		if perMatch != nil {
			return perMatch, nil
		}
		lock, err := per.permissionRepo.Lock(ctx, lockPerMatch, 1, lockTimeout) // 抢占分布式锁
		if err != nil {
			logger.Logger.Errorw("msg", err.Error())
			return nil, err
		}
		if !lock {
			<-time.After(timeSleep)
			continue
		}
		break
	}
	defer per.permissionRepo.UnLock(ctx, lockPerMatch) // 删除锁
	permitGroup, err := per.perGroupRepo.Find(per.db, &models.PerGroupQuery{
		AppID: appID,
	})
	if err != nil {
		logger.Logger.Errorw(err.Error())
	}
	permit := make([]string, len(permitGroup))

	for index, value := range permitGroup {
		permit[index] = value.ID
	}
	perGrant, err := per.permitGrantRepo.Find(per.db, &models.PerGrantQuery{
		Owners:  []string{userID, depID},
		Permits: permit,
	})
	if err != nil {
		return nil, err
	}
	if len(perGrant) == 0 { // 数据库没有数据
		return nil, nil
	}
	perMatch := &models.PermissionMatch{
		UserID:     userID,
		AppID:      appID,
		PerGroupID: perGrant[0].PerGroupID,
	}
	err = per.permissionRepo.CreatePerMatch(ctx, perMatch)
	if err != nil {
		logger.Logger.Errorw("msg", err.Error())
	}
	return perMatch, nil
}

func (per *permission) getPerInfo(ctx context.Context, perGroupID, formID string) (*models.Permission, error) {
	for i := 0; i < 5; i++ {
		// 1. 去redis 查询
		permission, err := per.permissionRepo.Get(ctx, perGroupID, formID)
		if err != nil {
			logger.Logger.Errorw("msg", err.Error())
			return nil, err
		}
		if permission != nil {
			return permission, nil
		}
		lock, err := per.permissionRepo.Lock(ctx, lockPermission, 1, lockTimeout) // 抢占分布式锁
		if err != nil {
			logger.Logger.Errorw("msg", err.Error())
			return nil, err
		}

		if !lock {
			<-time.After(timeSleep)
			continue
		}
		break
	}
	defer per.permissionRepo.UnLock(ctx, lockPermission) // 删除锁
	perGroup, err := per.perGroupRepo.Get(per.db, perGroupID)
	if err != nil {
		logger.Logger.Errorw("msg", err.Error())
		return nil, err
	}
	if perGroup == nil { // 数据库为空，直接返回，
		return nil, nil
	}
	permission := &models.Permission{
		PerGroupID: perGroup.ID,
		AppID:      perGroup.AppID,
		FormID:     formID,
		Name:       perGroup.Name,
		Type:       perGroup.Types,
	}
	permitForm, err := per.perFormRepo.Get(per.db, perGroupID, formID)
	if err != nil {
		return nil, err
	}
	if permitForm != nil {
		permission.Authority = permitForm.Authority
		permission.Conditions = permitForm.Conditions
		permission.Filter = permitForm.FieldJSON
	}
	err = per.permissionRepo.Create(ctx, permission, perTime)
	if err != nil {
		logger.Logger.Errorw("msg", err.Error())
	}
	return permission, nil
}

type SaveUserPerMatchReq struct {
	UserID     string `json:"userID"`
	AppID      string `json:"appID"`
	PerGroupID string `json:"perGroupID"`
}

type SaveUserPerMatchResp struct {
}

func (per *permission) SaveUserPerMatch(ctx context.Context, req *SaveUserPerMatchReq) (*SaveUserPerMatchResp, error) {
	matchPer := &models.PermissionMatch{
		UserID:     req.UserID,
		PerGroupID: req.PerGroupID,
		AppID:      req.AppID,
	}
	err := per.permissionRepo.CreatePerMatch(ctx, matchPer)
	if err != nil {
		return nil, err
	}
	return &SaveUserPerMatchResp{}, nil
}

type GetPerInCacheReq struct {
	UserID string `json:"userID"`
	DepID  string `json:"depID"`
	FormID string `json:"formID"`
	AppID  string `json:"appId"`
}

type GetPerInCacheResp struct {
	ID            string                  `json:"id"`
	AppID         string                  `json:"appID"`
	FormID        string                  `json:"formID"`
	Name          string                  `json:"name"`
	CreatedBy     string                  `json:"createdBy"`
	Description   string                  `json:"description"`
	Authority     int64                   `json:"authority"`
	DataAccessPer map[string]models.Query `json:"data_access_per"`
	Type          models.PerType          `json:"type"`
	Filter        map[string]interface{}  `json:"filter"`
}

func (per *permission) GetPerInCache(ctx context.Context, req *GetPerInCacheReq) (*GetPerInCacheResp, error) {
	// 1、 根据用户id 、  应用id ，得到 ====》 对应的权限组ID，
	perMatch, err := per.getPerMatch(ctx, req.UserID, req.DepID, req.AppID)
	if err != nil {
		return nil, err
	}
	if perMatch == nil {
		return nil, nil
	}
	// 2、 根据用户权限组id ，得到 权限信息
	permission, err := per.getPerInfo(ctx, perMatch.PerGroupID, req.FormID)
	if err != nil {
		return nil, err
	}
	if permission == nil {
		return nil, nil
	}
	resp := &GetPerInCacheResp{}
	per.cloneField(permission, resp)

	return resp, nil
}

func (per *permission) cloneField(permission *models.Permission, resp *GetPerInCacheResp) {
	resp.ID = permission.PerGroupID
	resp.DataAccessPer = permission.Conditions
	resp.Authority = permission.Authority
	resp.FormID = permission.FormID
	resp.Type = permission.Type
	resp.Filter = permission.Filter
}

type DelPerGroupReq struct {
	ID    string `json:"id"`
	AppID string `json:"appID"`
}

type DelPerGroupResp struct {
}

func (per *permission) DelPerGroup(ctx context.Context, req *DelPerGroupReq) (*DelPerGroupResp, error) {

	forms, err := per.perFormRepo.Find(per.db, &models.PerFormQuery{
		PerGroupID: req.ID,
	})
	if err != nil {
		return nil, err
	}
	err = per.perGroupRepo.Delete(per.db, &models.PerGroupQuery{
		ID: req.ID,
	})
	if err != nil {
		return nil, err
	}
	err = per.perFormRepo.Delete(per.db, &models.PerFormQuery{PerGroupID: req.ID})
	if err != nil {
		return nil, err
	}
	err = per.permitGrantRepo.Delete(per.db, &models.PerGrantQuery{
		PerGroupID: req.ID,
	})
	// TODO 删除缓存
	err = per.permissionRepo.DeletePerMatch(ctx, req.AppID)
	if err != nil {
		logger.Logger.Errorw(err.Error())
	}
	for _, value := range forms {
		err = per.permissionRepo.Delete(ctx, &models.Permission{
			PerGroupID: req.ID,
			FormID:     value.FormID,
		})
		if err != nil {
			logger.Logger.Errorw(err.Error())
		}
	}
	return &DelPerGroupResp{}, nil
}

type AddOwnerToGroupReq struct {
	ID     string     `json:"id"`
	AppID  string     `json:"appID"`
	Scopes []*grantVO `json:"scopes"`
}

type AddOwnerToGroupResp struct {
}

func (per *permission) AddOwnerToGroup(ctx context.Context, req *AddOwnerToGroupReq) (*AddOwnerToGroupResp, error) {
	// 先删除
	err := per.permitGrantRepo.Delete(per.db, &models.PerGrantQuery{
		PerGroupID: req.ID,
	})
	if err != nil {
		return nil, err
	}
	grants := make([]*models.PerGrant, len(req.Scopes))
	for index, value := range req.Scopes {
		r := &models.PerGrant{
			PerGroupID: req.ID,
			Owner:      value.ID,
			OwnerName:  value.Name,
			Types:      value.Types,
		}
		grants[index] = r
	}
	err = per.permitGrantRepo.BatchCreate(per.db, grants...)
	if err != nil {
		return nil, err
	}
	err = per.permissionRepo.DeletePerMatch(ctx, req.AppID)
	if err != nil {
		//logger.Logger.Errorw("delete redis error ", logger.STDRequestID(ctx), err.Error())
	}
	return &AddOwnerToGroupResp{}, nil

}

type AddOwnerToAppReq struct {
	AppID string
}

type AddOwnerToAppResp struct {
}

func (per *permission) AddOwnerToApp(ctx context.Context, req *AddOwnerToAppReq) (*AddOwnerToAppResp, error) {
	//
	groups, err := per.perGroupRepo.Find(per.db, &models.PerGroupQuery{AppID: req.AppID})
	if err != nil {
		return nil, err
	}
	permits := make([]string, len(groups))

	for index, value := range groups {
		permits[index] = value.ID
	}
	grants, err := per.permitGrantRepo.Find(per.db, &models.PerGrantQuery{Permits: permits})
	if err != nil {
		return nil, err
	}
	cache := make(map[string]struct{})
	owns := make([]string, 0)
	for _, value := range grants {
		_, ok := cache[value.Owner]
		if ok {
			continue
		}
		owns = append(owns, value.Owner)
		cache[value.Owner] = struct{}{}
	}
	_, err = per.appClient.AddAppScope(ctx, req.AppID, owns)
	if err != nil {
		return nil, err
	}
	return &AddOwnerToAppResp{}, nil
}
