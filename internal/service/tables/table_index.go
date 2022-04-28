package tables

import (
	"context"
	"github.com/quanxiang-cloud/form/internal/service/consensus"
	"github.com/quanxiang-cloud/form/pkg/misc/client"
	"github.com/quanxiang-cloud/form/pkg/misc/config"
)

type tableIndex struct {
	conf       *config.Config
	FormDDLAPI *client.FormDDLAPI
}

func (t *tableIndex) Do(ctx context.Context, bus *Bus) (*DoResponse, error) {
	_, err := t.FormDDLAPI.Index(ctx, consensus.GetTableID(bus.AppID, bus.TableID), "created_at", "created_at")
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func newTableIndex(conf *config.Config) (Guidance, error) {
	formDDLAPI, err := client.NewFormDDLAPI(conf)
	if err != nil {
		return nil, err
	}
	return &tableIndex{
		conf:       conf,
		FormDDLAPI: formDDLAPI,
	}, nil
}
