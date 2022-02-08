package form

import (
	"context"

	"github.com/quanxiang-cloud/form/pkg/client"
)

type comet struct {
	formClient *client.FormAPI
}

func NewForm() Form {
	return &comet{}
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
	}, dsl, req.TableID)
	if err != nil {
		return nil, err
	}

	return &SearchResp{
		Total:        searchResp.Total,
		Entities:     searchResp.Entities,
		Aggregations: searchResp.Aggregations,
	}, nil
}
