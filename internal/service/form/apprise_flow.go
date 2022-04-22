package form

import (
	"context"

	"github.com/quanxiang-cloud/form/pkg/misc/config"

	"github.com/quanxiang-cloud/cabin/logger"
	"github.com/quanxiang-cloud/form/internal/service/consensus"
	"github.com/quanxiang-cloud/form/internal/service/form/inform"
)

type appriseFlow struct {
	next   consensus.Guidance
	inform *inform.HookManger
}

func NewAppriseFlow(conf *config.Config) (consensus.Guidance, error) {
	form, err := newForm()
	if err != nil {
		return nil, err
	}
	ctx := context.Background()
	manger, err := inform.NewHookManger(ctx, conf)
	if err != nil {
		return nil, err
	}
	go manger.Start(ctx)
	return &appriseFlow{
		next:   form,
		inform: manger,
	}, nil
}

// Do 可以用策略模式改，可以先用switch.
func (a *appriseFlow) Do(ctx context.Context, bus *consensus.Bus) (*consensus.Response, error) {
	//	先去创建数据
	do, err := a.next.Do(ctx, bus)
	if err != nil {
		return nil, err
	}
	if do.Total == 0 {
		return do, err
	}
	//create update delete
	switch bus.Method {
	case "create":
		a.createApprise(ctx, bus)
	case "update":
		a.updateApprise(ctx, bus)
	case "delete":
		a.deleteApprise(ctx, bus)
	}
	return do, nil
}

func (a *appriseFlow) createApprise(ctx context.Context, bus *consensus.Bus) {
	data := new(inform.FormData)
	data.TableID = bus.TableID
	data.Entity = bus.CreatedOrUpdate.Entity
	inform.DefaultFormFiled(ctx, data, "post")
	logger.Logger.Infow("create", "data is ", data)
	a.inform.Send <- data
}

func (a *appriseFlow) deleteApprise(ctx context.Context, bus *consensus.Bus) {
	ids := make([]string, 0)
	if len(bus.Get.OldQuery) == 0 {
		ids = consensus.GetIDByQuery(bus.Get.Query)
	}
	ids = consensus.GetIDByQuery(bus.Get.OldQuery)
	data := &inform.FormData{
		TableID: bus.TableID,
		Entity: map[string]interface{}{
			"data":      ids,
			"delete_id": bus.UserID,
		},
	}
	inform.DefaultFormFiled(ctx, data, "delete")
	logger.Logger.Infow("delete", "data is ", data)
	a.inform.Send <- data
}

func (a *appriseFlow) updateApprise(ctx context.Context, bus *consensus.Bus) {
	ids := make([]string, 0)
	if len(bus.Get.OldQuery) == 0 {
		ids = consensus.GetIDByQuery(bus.Get.Query)
	}
	ids = consensus.GetIDByQuery(bus.Get.OldQuery)
	for _, id := range ids {
		entity := consensus.DefaultField(bus.CreatedOrUpdate.Entity,
			consensus.WithUpdateID(id),
		)
		data := &inform.FormData{
			TableID: bus.TableID,
			Entity:  entity,
		}
		inform.DefaultFormFiled(ctx, data, "put")
		logger.Logger.Infow("update", "data is ", data)
		a.inform.Send <- data
	}
}
