package tables

import (
	"context"
	"github.com/quanxiang-cloud/form/internal/models"
	"github.com/quanxiang-cloud/form/pkg/misc/config"
)

type component struct {
	tableRelation models.TableRelationRepo
	next          Guidance
}

func (c *component) Do(ctx context.Context, bus *Bus) (*DoResponse, error) {
	// 关联关系表
	//for key, value := range bus.Schema {
	//	columnValue, err :=getAsMap(value)
	//	if err != nil {
	//		return nil ,err
	//	}
	//	types ,err := getMapToString(columnValue, "x-component")
	//	if err != nil {
	//		return nil ,err
	//	}
	//
	//	componentData
	//
	//
	//}

	// 流水号

	return c.next.Do(ctx, bus)
}

func newComponent(conf *config.Config) (Guidance, error) {
	swagger, err := newRegisterSwagger(conf)
	if err != nil {
		return nil, err
	}
	return &component{
		next: swagger,
	}, nil
}
