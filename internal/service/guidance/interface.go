package guidance

import (
	"context"

	"github.com/quanxiang-cloud/form/internal/service/consensus"
	"github.com/quanxiang-cloud/form/pkg/misc/config"
)

type Guidance interface {
	Do(ctx context.Context, bus *consensus.Bus) (consensus.Response, error)
}

func New(conf *config.Config) (Guidance, error) {
	return newCertifier(conf)
}
