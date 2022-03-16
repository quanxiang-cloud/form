package guidance

import (
	"context"

	"github.com/quanxiang-cloud/form/internal/service/consensus"
	"github.com/quanxiang-cloud/form/internal/service/form"
	"github.com/quanxiang-cloud/form/pkg/misc/config"
)

type structor struct {
	next consensus.Guidance
}

func newStructor(conf *config.Config) (consensus.Guidance, error) {
	conditions, err := form.NewCondition(conf)
	if err != nil {
		return nil, err
	}
	return &structor{
		next: conditions,
	}, nil
}

func (s *structor) Do(ctx context.Context, bus *consensus.Bus) (*consensus.Response, error) {
	return s.next.Do(ctx, bus)
}
