package tables

import (
	"context"
	"encoding/base64"

	"github.com/quanxiang-cloud/cabin/logger"
	"github.com/quanxiang-cloud/form/internal/service/tables/swagger"
	"github.com/quanxiang-cloud/form/internal/service/tables/util"

	"github.com/quanxiang-cloud/form/pkg/misc/client"
	"github.com/quanxiang-cloud/form/pkg/misc/config"
)

type registerSwagger struct {
	conf    *config.Config
	polyAPI client.PolyAPI
	next    Guidance
}

func (reg *registerSwagger) Do(ctx context.Context, bus *Bus) (*DoResponse, error) {
	schema, require := util.GetSpecSchema(bus.ConvertSchema)
	swagger, err := swagger.DoSchemas(bus.AppID, bus.TableID, bus.Title, schema, require)
	if err != nil {
		return nil, err
	}
	regSwagger, err := reg.polyAPI.RegSwagger(ctx, "form", base64.StdEncoding.EncodeToString([]byte(swagger)), bus.AppID, bus.TableID, bus.Title)
	if err != nil {
		return nil, err
	}
	logger.Logger.Errorw("msg", "request-id", regSwagger)

	return reg.next.Do(ctx, bus)
}

func newRegisterSwagger(conf *config.Config) (Guidance, error) {
	index, err := newTableIndex(conf)
	if err != nil {
		return nil, err
	}
	return &registerSwagger{
		conf:    conf,
		next:    index,
		polyAPI: client.NewPolyAPI(conf),
	}, nil
}
