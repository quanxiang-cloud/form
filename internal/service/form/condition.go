package form

import (
	"context"
	"github.com/quanxiang-cloud/form/internal/service/consensus"
)

type condition struct {
	next consensus.Guidance
}

// Do Do
func (c *condition) Do(ctx context.Context, bus *consensus.Bus) (consensus.Response, error) {
	return nil, nil
}
