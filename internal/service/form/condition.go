package form

import (
	"context"
	"github.com/quanxiang-cloud/form/internal/service/consensus"
	"github.com/quanxiang-cloud/form/pkg/misc/config"
)

type condition struct {
	next consensus.Guidance
}

func NewCondition(conf *config.Config) (consensus.Guidance, error) {
	newRefs, err := NewRefs(conf)
	if err != nil {
		return nil, err
	}
	return &condition{
		next: newRefs,
	}, nil
}

// Do Do
func (c *condition) Do(ctx context.Context, bus *consensus.Bus) (*consensus.Response, error) {
	return c.next.Do(ctx, bus)
}
