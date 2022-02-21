package guidance

import (
	"context"
	"github.com/quanxiang-cloud/form/internal/service/consensus"
	"github.com/quanxiang-cloud/form/internal/service/form"
)

type structor struct {
	form form.Form
}

func newStructor() (consensus.Guidance, error) {
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
	base := form.Base{
		AppID:   bus.AppID,
		TableID: bus.TableID,
		UserID:  bus.UserID,
	}
	switch bus.Foundation.Method {
	case "get":
		req := &form.GetReq{
			Base:  base,
			Query: bus.Query,
		}
		req.Base = base
		req.Query = bus.Query
		return s.form.Get(ctx, req)

	case "search":
		req := &form.SearchReq{
			Sort:  bus.List.Sort,
			Page:  bus.List.Page,
			Size:  bus.List.Size,
			Query: bus.Query,
			Base:  base,
		}
		return s.form.Search(ctx, req)
	case "create":
		req := &form.CreateReq{
			Entity: bus.CreatedOrUpdate.Entity,
			Base:   base,
		}
		return s.form.Create(ctx, req)
	case "update":
		req := &form.UpdateReq{
			Entity: bus.CreatedOrUpdate.Entity,
			Query:  bus.Query,
			Base:   base,
		}
		return s.form.Update(ctx, req)
	case "delete":
		req := &form.DeleteReq{
			Query: bus.Query,
			Base:  base,
		}
		return s.form.Delete(ctx, req)
	}
	return nil, nil
}
