package tables

import (
	"context"

	"github.com/quanxiang-cloud/cabin/logger"
	"github.com/quanxiang-cloud/form/internal/service/tables/swagger"
	"github.com/quanxiang-cloud/form/internal/service/tables/util"
	"github.com/quanxiang-cloud/form/pkg/misc/client"
	"github.com/quanxiang-cloud/form/pkg/misc/config"
)

type registerSwagger struct {
	conf    *config.Config
	polyAPI client.PolyAPI
}

func (reg *registerSwagger) Do(ctx context.Context, bus *Bus) (*DoResponse, error) {
	schema := util.GetSpecSchema(bus.ConvertSchema)
	swagger, err := swagger.DoSchemas(bus.AppID, bus.TableID, bus.Title, schema)
	if err != nil {
		return nil, err
	}
	regSwagger, err := reg.polyAPI.RegSwagger(ctx, "form", swagger, bus.AppID, bus.TableID, bus.Title)
	if err != nil {
		return nil, err
	}
	logger.Logger.Errorw("msg", "request-id", regSwagger)
	return nil, nil
}

func newRegisterSwagger(conf *config.Config) (Guidance, error) {
	return &registerSwagger{
		conf:    conf,
		polyAPI: client.NewPolyAPI(conf.InternalNet),
	}, nil
}
