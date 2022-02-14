package guidance

import (
	"context"
	"github.com/quanxiang-cloud/form/internal/service/consensus"
	"github.com/quanxiang-cloud/form/internal/service/form"
)

type structor struct {
	form form.Form
}

func newStructor() (Guidance, error) {
	form, err := form.NewForm()
	if err != nil {
		return nil, err
	}
	return &structor{
		form: form,
	}, nil
}

func (s *structor) Do(ctx context.Context, bus *consensus.Bus) (consensus.Response, error) {
	// TODO
	switch bus.Foundation.Method {
	case "get":
		req := &form.GetReq{}
		req.TableID = bus.TableID
		req.AppID = bus.AppID
		req.Query = bus.Query
		req.UserID = bus.UserID
		return s.form.Get(ctx, req)

	case "search":
		req := &form.SearchReq{
			Sort:  bus.Sort,
			Page:  bus.Page,
			Size:  bus.Size,
			Query: bus.Query,
		}
		req.TableID = bus.TableID
		req.AppID = bus.AppID
		req.UserID = bus.UserID
		return s.form.Search(ctx, req)
	case "create":
		req := &form.CreateReq{
			Entity: bus.CreatedOrUpdate.Entity,
		}
		req.TableID = bus.TableID
		req.AppID = bus.AppID
		req.UserID = bus.UserID
		return s.form.Create(ctx, req)
	case "update":
		req := &form.UpdateReq{
			Entity: bus.CreatedOrUpdate.Entity,
			Query:  bus.Query,
		}
		req.TableID = bus.TableID
		req.AppID = bus.AppID
		req.UserID = bus.UserID
		return s.form.Update(ctx, req)
	case "delete":
		req := &form.UpdateReq{
			Query: bus.Query,
		}
		req.TableID = bus.TableID
		req.AppID = bus.AppID
		req.UserID = bus.UserID
		return s.form.Update(ctx, req)
	}
	return nil, nil
}
