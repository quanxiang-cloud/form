package form

import (
	"context"
	"github.com/quanxiang-cloud/form/internal/service/form/base"
	"github.com/quanxiang-cloud/form/internal/service/types"
	"reflect"
)

type common struct {
	comet         *comet
	userID        string
	depID         string
	userName      string
	tag           string
	key           string
	refValue      types.Ref   // ref 结构
	primaryEntity base.Entity // 主表的entity
	oldValue      types.M     // 透传的数据
}

type subTable struct {
	common
}

type comReq struct {
	comet         *comet
	userID        string
	userName      string
	depID         string
	tag           string
	key           string
	refValue      types.Ref
	primaryEntity base.Entity
	oldValue      types.M
}

func (s *subTable) getTag() string {
	return "sub_table"
}

func (c *common) setValue(com *comReq) {
	c.comet = com.comet
	c.userID = com.userID
	c.depID = com.depID
	c.userName = com.userName
	c.tag = com.tag
	c.key = com.key
	c.refValue = com.refValue
	c.primaryEntity = com.primaryEntity
	c.oldValue = com.oldValue
}

func (s *subTable) handlerFunc(ctx context.Context, action string) error {

	switch action {
	case post: // 更新
		return s.common.subCreate(ctx)
	case put:
		return s.common.subUpdate(ctx)
	case get:
		return s.common.subGet(ctx, true)
	}
	return nil
}

func (c *common) subCreate(ctx context.Context) error {
	originalData := &OriginalData{}
	refData := &RefData{}
	err := c.perInitData(ctx, refData, originalData)
	if err != nil {
		return err
	}
	err = c.new(ctx, refData, originalData)
	if err != nil {
		return err
	}
	return nil
}

func (c *common) subUpdate(ctx context.Context) error {
	originalData := &OriginalData{}
	refData := &RefData{}
	err := c.perInitData(ctx, refData, originalData)
	if err != nil {
		return err
	}

	return nil
}

func (c *common) subGet(ctx context.Context, isReplace bool) error {
	return nil
}

// new :component create options
func (c *common) new(ctx context.Context, refData *RefData, originalData *OriginalData) error {
	for _, _subTable := range refData.New {
		subTable, ok := _subTable.(map[string]interface{})
		if !ok {
			continue
		}
		opt := &OptionEntity{}
		err := initOptionEntity(ctx, subTable, opt)
		if err != nil {
			return err
		}
		req := &CreateReq{
			Entity: opt.Entity,
			Ref:    opt.Ref,
		}
		req.AppID = refData.AppID
		req.TableID = refData.TableID
		req.UserName = c.userName
		req.UserID = c.userID
		_, err = c.comet.Create(ctx, req)
		if err != nil {
			return err
		}
		//createResp.Entity
		//e, ok := createResp.Entity.(map[string]interface{})
		//if !ok {
		//	continue

		//}
		//ids = append(ids, e[IDKey].(string))
	}
	return nil
}

func (c *common) update(ctx context.Context, refData *RefData) error {
	for _, _subTable := range refData.Updated {
		subTable, ok := _subTable.(map[string]interface{})
		if !ok {
			continue
		}
		opt := &OptionEntity{}
		err := initOptionEntity(ctx, subTable, opt)
		if err != nil {
			return err
		}
		//req := &UpdateFormReq{
		//	IsAuth:  client.NoAuth,
		//	AppID:   refData.AppID,
		//	TableID: refData.TableID,
		//	Entity:  opt.Entity,
		//	Query:   opt.Query,
		//	Profile: c.Profile,
		//}
		//_, err = c.form.UpdateForm(ctx, req, UpdateOption(c.form))
		//if err != nil {
		//	return err
		//}
	}
	return nil
}

func (c *common) delete(ctx context.Context, refData *RefData, originalData *OriginalData, isDelete bool) error {
	return nil
}

func (c *common) perInitData(ctx context.Context, refData *RefData, original *OriginalData) error {
	var err error
	if refData != nil {
		err = initRefData(ctx, c.refValue, refData)
	}

	if err != nil {
		return err
	}
	if original != nil {
		err = initOriginalData(ctx, c.oldValue, original)
	}
	if err != nil {
		return err
	}
	return nil
}

func getPrimaryID(e base.Entity) (string, error) {
	if e == nil {
		return "", nil
	}
	value := reflect.ValueOf(e)
	switch _t := reflect.TypeOf(e); _t.Kind() {
	case reflect.Map:
		if value := value.MapIndex(reflect.ValueOf("_id")); value.IsValid() {
			return value.String(), nil
		}
	default:
		// 返回错误
		return "", nil
	}

	return "", nil
}
