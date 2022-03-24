package form

import (
	"context"
	"reflect"

	"github.com/quanxiang-cloud/cabin/logger"
	redis2 "github.com/quanxiang-cloud/cabin/tailormade/db/redis"
	"github.com/quanxiang-cloud/form/internal/models"
	"github.com/quanxiang-cloud/form/internal/models/redis"
	"github.com/quanxiang-cloud/form/internal/service/consensus"
	"github.com/quanxiang-cloud/form/internal/service/types"
	"github.com/quanxiang-cloud/form/pkg/misc/config"
)

type refs struct {
	next       consensus.Guidance
	component  *component
	serialRepo models.SerialRepo
}

func NewRefs(conf *config.Config) (consensus.Guidance, error) {
	appriseFlows, err := NewAppriseFlow(conf)
	if err != nil {
		return nil, err
	}
	redisClient, err := redis2.NewClient(conf.Redis)
	if err != nil {
		return nil, err
	}
	return &refs{
		next:       appriseFlows,
		component:  newFormComponent(),
		serialRepo: redis.NewSerialRepo(redisClient),
	}, nil
}

// Do create update
func (c *refs) Do(ctx context.Context, bus *consensus.Bus) (*consensus.Response, error) {
	if bus.Method == "create" {
		bus.CreatedOrUpdate.Entity = consensus.DefaultField(bus.CreatedOrUpdate.Entity,
			consensus.WithID(), consensus.WithCreated(bus.UserID, bus.UserName))
	}
	if bus.Method == "update" {
		bus.CreatedOrUpdate.Entity = consensus.DefaultField(bus.CreatedOrUpdate.Entity,
			consensus.WithUpdated(bus.UserID, bus.UserName))
	}

	for fieldKey, value := range bus.Ref.Ref {
		fieldValue, ok := value.(map[string]interface{})
		if !ok {
			logger.Logger.Errorw("msg", "", "key is error ")
			continue
		}
		t := fieldValue["type"]
		if reflect.ValueOf(t).Kind() == reflect.String {
			comReqs := &comReq{
				ref:           c,
				userID:        bus.UserID,
				userName:      bus.UserName,
				depID:         bus.DepID,
				tag:           reflect.ValueOf(t).String(),
				key:           fieldKey,
				refValue:      fieldValue,
				primaryEntity: bus.CreatedOrUpdate.Entity,
				extraValue: types.M{
					appIDKey:   bus.AppID,
					tableIDKey: bus.TableID,
				},
			}
			if bus.Method == "update" {
				ids := consensus.GetIDByQuery(bus.Get.Query)
				if len(ids) > 0 {
					comReqs.primaryEntity = types.M{
						"_id": ids[0],
					}
				}
			}
			com, err := c.component.getCom(reflect.ValueOf(t).String(), comReqs)
			if err != nil {
				logger.Logger.Errorw("msg", "", "key is error ")
				continue
			}
			err = com.handlerFunc(ctx, bus.Method)
			if err != nil {
				logger.Logger.Errorw("msg", "", err.Error())
			}

		}
	}
	return c.next.Do(ctx, bus)
}
