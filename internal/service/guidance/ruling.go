package guidance

import (
	"context"

	"github.com/quanxiang-cloud/form/internal/service/consensus"
)

type ruling struct {
	poly     Guidance
	structor Guidance
}

func newRuling() (Guidance, error) {
	poly, err := newPoly()
	if err != nil {
		return nil, err
	}
	structor, err := newStructor()
	if err != nil {
		return nil, err
	}
	return &ruling{
		poly:     poly,
		structor: structor,
	}, nil
}

func (r *ruling) Do(ctx context.Context, bus *consensus.Bus) (consensus.Response, error) {
	switch bus.Method {
	case "get", "search", "create", "update", "delete":
		// Specified method, distributed to form processor
		return r.structor.Do(ctx, bus)
	default:
		return r.poly.Do(ctx, bus)
	}
}
