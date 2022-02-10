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
	components *Component
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
		components: NewCom(),
	}, nil
}

func (c *comet) Search(ctx context.Context, req *SearchReq) (*SearchResp, error) {
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
	searchResp, err := c.formClient.Search(ctx, client.FindOptions{
		Sort: req.Sort,
		Size: req.Size,
		Page: req.Page,
	}, dsl, getTableID(req.AppID, req.TableID))
	if err != nil {
		return nil, err
	}

	return &SearchResp{
		Total:        searchResp.Total,
		Entities:     searchResp.Entities,
		Aggregations: searchResp.Aggregations,
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
			com, err := c.components.GetCom(reflect.ValueOf(t).String(), req)
			if err != nil {
				continue
			}
			err = com.HandlerFunc(ctx, method)
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
		WithCreated(req.UserID, req.CreatorName))

	insert, err := c.formClient.Insert(ctx, getTableID(req.AppID, req.TableID), req.Entity)

	if err != nil {
		return nil, err
	}
	return &CreateResp{
		ErrorCount: insert.SuccessCount,
	}, nil
}

func getTableID(appID, tableID string) string {
	if len(appID) == 36 {
		return fmt.Sprintf("%s%s%s", "A", appID, tableID)
	}
	return fmt.Sprintf("%s%s%s", "a", appID, tableID)
}
