package form

import (
	"context"
	"fmt"
	"github.com/quanxiang-cloud/cabin/logger"
	"github.com/quanxiang-cloud/form/internal/service/form/base"
	"github.com/quanxiang-cloud/form/internal/service/form/inform"
	"github.com/quanxiang-cloud/form/internal/service/types"
	"reflect"

	"github.com/quanxiang-cloud/form/pkg/client"
)

type comet struct {
	formClient *client.FormAPI
	components *component
	hook       *inform.HookManger
}

func NewForm() (Form, error) {
	return newForm()
}

func newForm() (*comet, error) {
	manger, err := inform.NewHookManger(context.Background())
	if err != nil {
		return nil, err
	}
	go manger.Start(context.Background())

	formApi, err := client.NewFormAPI()
	if err != nil {
		return nil, err
	}

	return &comet{
		formClient: formApi,
		components: newFormComponent(),
		hook:       manger,
	}, nil
}

func (c *comet) Search(ctx context.Context, req *SearchReq, opts ...inform.Options) (*SearchResp, error) {
	defer func() {
		after(ctx, &inform.OptionReq{}, opts...)
	}()
	return c.callSearch(ctx, req)
}

func after(ctx context.Context, req *inform.OptionReq, opts ...inform.Options) {
	for _, opt := range opts {
		err := opt(ctx, req)
		if err != nil {
			logger.Logger.Errorw(err.Error())
		}
	}
}
func (c *comet) callSearch(ctx context.Context, req *SearchReq) (*SearchResp, error) {
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
	formReq := &client.FormReq{
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

	return &SearchResp{
		Total:    searchResp.Total,
		Entities: searchResp.Entities,
	}, nil
}

//Create Create
func (c *comet) Create(ctx context.Context, req *CreateReq, opts ...inform.Options) (*CreateResp, error) {
	defer func() {
		after(ctx, &inform.OptionReq{}, opts...)
	}()

	resp, err := c.callCreate(ctx, req)
	if err != nil {
		return nil, err
	}
	// 处理ref 等高级字段的数据   //
	comReq := &comReq{
		comet:         c,
		userID:        req.UserID,
		depID:         req.DepID,
		primaryEntity: req.Entity,
		refValue:      req.Ref,
		oldValue: types.M{
			appIDKey:   req.AppID,
			tableIDKey: req.TableID,
		},
	}

	err = c.getManyCom(ctx, comReq, post)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *comet) Update(ctx context.Context, req *UpdateReq, opts ...inform.Options) (*UpdateResp, error) {
	defer func() {
		after(ctx, &inform.OptionReq{}, opts...)
	}()
	return c.callUpdate(ctx, req)
}

func (c *comet) callUpdate(ctx context.Context, req *UpdateReq) (*UpdateResp, error) {
	req.Entity = base.DefaultField(req.Entity,
		base.WithUpdated(req.UserID, req.UserName))
	dsl := make(map[string]interface{})
	if req.Query != nil {
		dsl["query"] = req.Query
	}
	if len(dsl) == 0 {
		dsl = nil
	}

	formReq := &client.FormReq{
		Entity:   req.Entity,
		TableID:  getTableID(req.AppID, req.TableID),
		DslQuery: dsl,
	}
	update, err := c.formClient.Update(ctx, formReq)

	if err != nil {
		return nil, err
	}
	return &UpdateResp{
		Count: update.SuccessCount,
	}, nil
}

func (c *comet) getManyCom(ctx context.Context, req *comReq, method string) error {
	for fieldKey, value := range req.refValue {
		optionValue, ok := value.(map[string]interface{})
		if !ok {
			continue
		}
		t := optionValue[_type]
		if reflect.ValueOf(t).Kind() == reflect.String {
			req.tag = reflect.ValueOf(t).String()
			req.key = fieldKey
			req.refValue = optionValue
			com, err := c.components.getCom(reflect.ValueOf(t).String(), req)
			if err != nil {
				continue
			}
			err = com.handlerFunc(ctx, method)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (c *comet) callCreate(ctx context.Context, req *CreateReq) (*CreateResp, error) {
	req.Entity = base.DefaultField(req.Entity,
		base.WithID(),
		base.WithCreated(req.UserID, req.UserName))

	formReq := &client.FormReq{
		Entity:  req.Entity,
		TableID: getTableID(req.AppID, req.TableID),
	}
	insert, err := c.formClient.Insert(ctx, formReq)

	if err != nil {
		return nil, err
	}
	return &CreateResp{
		Count: insert.SuccessCount,
	}, nil
}

func getTableID(appID, tableID string) string {
	if len(appID) == 36 {
		return fmt.Sprintf("%s%s%s", "A", appID, tableID)
	}
	return fmt.Sprintf("%s%s%s", "a", appID, tableID)
}

func (c *comet) Get(ctx context.Context, req *GetReq, opts ...inform.Options) (*GetResp, error) {
	defer func() {
		after(ctx, &inform.OptionReq{}, opts...)
	}()
	return c.callGet(ctx, req)
}

func (c *comet) callGet(ctx context.Context, req *GetReq) (*GetResp, error) {
	dsl := make(map[string]interface{})
	if req.Query != nil {
		dsl["query"] = req.Query
	}
	if len(dsl) == 0 {
		dsl = nil
	}

	formReq := &client.FormReq{
		DslQuery: dsl,
		TableID:  getTableID(req.AppID, req.TableID),
	}
	resp, err := c.formClient.Get(ctx, formReq)
	if err != nil {
		return nil, err
	}
	return &GetResp{
		Entity: resp.Entity,
	}, nil

}

func (c *comet) Delete(ctx context.Context, req *DeleteReq, opts ...inform.Options) (*DeleteResp, error) {
	defer func() {
		after(ctx, &inform.OptionReq{}, opts...)
	}()
	return c.callDelete(ctx, req)
}

func (c *comet) callDelete(ctx context.Context, req *DeleteReq) (*DeleteResp, error) {
	dsl := make(map[string]interface{})
	if req.Query != nil {
		dsl["query"] = req.Query
	}
	if len(dsl) == 0 {
		dsl = nil
	}
	formReq := &client.FormReq{
		DslQuery: dsl,
		TableID:  getTableID(req.AppID, req.TableID),
	}
	resp, err := c.formClient.Delete(ctx, formReq)
	if err != nil {
		return nil, err
	}
	return &DeleteResp{
		Count: resp.SuccessCount,
	}, nil
}
