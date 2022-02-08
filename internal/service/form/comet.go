package form

import (
	"context"
	"fmt"

	"github.com/quanxiang-cloud/form/pkg/client"
)

type comet struct {
	formClient *client.FormAPI
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

func getTableID(appID, tableID string) string {
	if len(appID) == 36 {
		return fmt.Sprintf("%s%s%s", "A", appID, tableID)
	}
	return fmt.Sprintf("%s%s%s", "a", appID, tableID)
}
