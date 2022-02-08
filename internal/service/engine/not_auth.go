package engine

import (
	"github.com/quanxiang-cloud/form/pkg/client"
	config2 "github.com/quanxiang-cloud/form/pkg/misc/config"
)

type noAuth struct {
	comet
}

// NewNoAuth NewNoAuth
func NewNoAuth(conf *config2.Config) (Plugs, error) {
	formApi, err := client.NewFormAPI()
	if err != nil {
		return nil, err
	}
	a := &noAuth{
		comet{
			formClient: formApi,
		},
	}
	return a, nil
}
