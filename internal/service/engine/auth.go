package engine

import (
	"context"
	error2 "github.com/quanxiang-cloud/cabin/error"
	comet2 "github.com/quanxiang-cloud/form/internal/comet"
	"github.com/quanxiang-cloud/form/internal/filters"
	"github.com/quanxiang-cloud/form/internal/models"
	"github.com/quanxiang-cloud/form/pkg/client"
	"github.com/quanxiang-cloud/form/pkg/misc/code"
	config2 "github.com/quanxiang-cloud/form/pkg/misc/config"
)

type auth struct {
	comet
}

type Plugs interface{}

// NewAuth NewAuth
func NewAuth(conf *config2.Config) (Plugs, error) {
	formApi, err := client.NewFormAPI()
	if err != nil {
		return nil, err
	}
	a := &auth{
		comet{
			formClient: formApi,
		},
	}
	return a, nil
}

func (c1 *auth) Pre(ctx context.Context, bus *bus, method string, opts ...PreOption) error {
	for _, opt := range opts {
		if !opt(method, bus, bus.entity) {
			return error2.New(code.ErrNotPermit)
		}
	}
	return nil
}

func getFormID(tableID string) string {
	switch tableID[0:1] {
	case "A":
		if len(tableID) > 37 {
			return tableID[37:]
		}
		return tableID
	case "a":
		if len(tableID) > 6 {
			return tableID[6:]
		}
	}
	return tableID
}

// PreOption PreOption
type PreOption func(method string, bus *bus, entity comet2.Entity) bool

// CheckOperate CheckOperate
func CheckOperate() PreOption {
	return func(method string, bus *bus, entity comet2.Entity) bool {
		if bus.permit.permitTypes == models.InitType {
			return true
		}
		var op int64
		switch method {
		case "find", "findOne":
			op = models.OPRead
		case "create":
			op = models.OPCreate
		case "update":
			op = models.OPUpdate
		case "delete":
			op = models.OPDelete
		default:
			return false
		}

		return op&bus.permit.Authority != 0
	}
}

// CheckData  CheckData
func CheckData() PreOption {
	return func(method string, bus *bus, entity comet2.Entity) bool {
		if entity == nil {
			return false
		}
		if bus.permit.permitTypes == models.InitType {
			return true
		}
		return filters.FilterCheckData(entity, bus.permit.Filter)
	}
}

//FilterOption FilterOption
type FilterOption func(data interface{}, filter map[string]interface{}) error
