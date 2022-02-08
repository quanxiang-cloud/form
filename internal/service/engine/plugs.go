package engine

import (
	"context"
	"github.com/quanxiang-cloud/form/pkg/client"
)

type SchemaHandle interface {
	SchemaHandle(ctx context.Context) error
}

type Search interface {
	Search(ctx context.Context, req *SearchReq) (*client.SearchResp, error)
}

type Create interface {
	Insert(ctx context.Context) (client.InsertResp, error)
}

type Update interface {
	Update(ctx context.Context) (client.UpdateResp, error)
}

type Delete interface {
	Delete(ctx context.Context) (client.SearchResp, error)
}

type Get interface {
	Get(ctx context.Context) (client.SearchResp, error)
}

type Pre interface {
	Pre(ctx context.Context, bus2 *bus, method string, opts ...PreOption) error
}

type Post interface {
	Postfix(ctx context.Context, data interface{}, bus *bus, opt ...FilterOption) error
}
