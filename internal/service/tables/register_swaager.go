package tables

import (
	"context"
	"github.com/quanxiang-cloud/form/pkg/misc/client"
	"github.com/quanxiang-cloud/form/pkg/misc/config"
)

type registerSwagger struct {
	conf    *config.Config
	polyAPI client.PolyAPI
}

func (reg *registerSwagger) Do(ctx context.Context, bus *Bus) (*DoResponse, error) {
	//genSwagger, err := swagger.GenSwagger(reg.conf, bus.ConvertSchemas.ConvertSchema, bus.Title, bus.AppID, bus.TableID)
	//if err != nil {
	//	return nil, err
	//}
	//content := "form"
	//if bus.Source == models.ModelSource {
	//	content = "custom"
	//}
	//_, err = reg.polyAPI.RegSwagger(ctx, swagger.Service, genSwagger, bus.AppID, content)
	//if err != nil {
	//	return nil, err
	//}
	return nil, nil
}

func newRegisterSwagger(conf *config.Config) (Guidance, error) {
	return &registerSwagger{
		conf:    conf,
		polyAPI: client.NewPolyAPI(conf.InternalNet),
	}, nil
}
