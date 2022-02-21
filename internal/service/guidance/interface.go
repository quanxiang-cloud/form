package guidance

import (
	"github.com/quanxiang-cloud/form/internal/service/consensus"
	"github.com/quanxiang-cloud/form/pkg/misc/config"
)

func New(conf *config.Config) (consensus.Guidance, error) {
	return newCertifier(conf)
}
