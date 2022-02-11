package form

import (
	"context"
	"fmt"
	"github.com/quanxiang-cloud/form/internal/service/types"
	"reflect"

	"github.com/quanxiang-cloud/form/pkg/client"
)

type comet struct {
	formClient *client.FormAPI
	components *component
}

func NewForm() (Form, error) {
	return newForm()
}

func newForm() (*comet, error) {
	formApi, err := client.NewFormAPI()
	if err != nil {
		return nil, err
	}

	return &comet{
		formClient: formApi,
		components: newFormComponent(),
	}, nil
}

func (c *comet) Search(ctx context.Context, req *SearchReq) (*SearchResp, error) {
	return c.callSearch(ctx, req)
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
func (c *comet) Create(ctx context.Context, req *CreateReq) (*CreateResp, error) {
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

func (c *comet) Update(ctx context.Context, req *UpdateReq) (*UpdateResp, error) {
	return c.callUpdate(ctx, req)
}

func (c *comet) callUpdate(ctx context.Context, req *UpdateReq) (*UpdateResp, error) {
	req.Entity = DefaultField(req.Entity,
		WithID(),
		WithUpdated(req.UserID, req.UserName))

	formReq := &client.FormReq{
		Entity:  req.Entity,
		TableID: getTableID(req.AppID, req.TableID),
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
	req.Entity = DefaultField(req.Entity,
		WithID(),
		WithCreated(req.UserID, req.UserName))

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

func (c *comet) Get(ctx context.Context, req *GetReq) (*GetResp, error) {
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

func (c *comet) Delete(ctx context.Context, req *DeleteReq) (*DeleteResp, error) {
	return nil, nil
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
