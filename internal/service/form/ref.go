package form

import (
	"context"
	"github.com/quanxiang-cloud/form/internal/service/consensus"
)

type refs struct {
	next consensus.Guidance
}

func (c *refs) Do(ctx context.Context, bus *consensus.Bus) (consensus.Response, error) {
	return nil, nil
}
