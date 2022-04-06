package form

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/quanxiang-cloud/form/internal/models/redis"
	"github.com/quanxiang-cloud/form/internal/service/consensus"
	"github.com/quanxiang-cloud/form/internal/service/types"
	"github.com/quanxiang-cloud/form/pkg/misc/utils"
	"reflect"
)

type common struct {
	ref           *refs
	userID        string
	depID         string
	userName      string
	tag           string
	key           string
	refValue      types.Ref        // ref 结构
	primaryEntity consensus.Entity // 主表的entity
	extraValue    types.M          // 透传的数据
}

type comReq struct {
	ref           *refs
	userID        string
	userName      string
	depID         string
	tag           string
	key           string
	refValue      types.Ref
	primaryEntity consensus.Entity
	extraValue    types.M
}

func (c *common) setValue(com *comReq) {
	c.ref = com.ref
	c.userID = com.userID
	c.depID = com.depID
	c.userName = com.userName
	c.tag = com.tag
	c.key = com.key
	c.refValue = com.refValue
	c.primaryEntity = com.primaryEntity
	c.extraValue = com.extraValue
}

type subTable struct {
	common
}

func (s *subTable) getTag() string {
	return "sub_table"
}

type foreignTable struct {
	common
}

func (f *foreignTable) getTag() string {
	return "foreign_table"
}

func (c *common) handlerFunc(ctx context.Context, action string) error {
	switch action {
	case "create": // 创建
		return c.subCreate(ctx)
	case "update":
		return c.subUpdate(ctx)
	case "get":
		return c.subGet(ctx, true)
	}
	return nil
}

func (c *common) subCreate(ctx context.Context) error {
	extraData := &ExtraData{}
	refData := &RefData{}
	err := c.perInitData(ctx, refData, extraData)
	if err != nil {
		return err
	}
	err = c.new(ctx, refData, extraData)
	if err != nil {
		return err
	}
	return nil
}

func (c *common) subUpdate(ctx context.Context) error {
	extraData := &ExtraData{}
	refData := &RefData{}
	err := c.perInitData(ctx, refData, extraData)
	if err != nil {
		return err
	}
	// deal with update option
	err = c.update(ctx, refData)
	if err != nil {
		//	logger.Logger.Errorw("update is  err is", err.Error())
	}
	err = c.new(ctx, refData, extraData)
	if err != nil {
		//logger.Logger.Errorw("new  is  err is", err.Error())
	}
	err = c.delete(ctx, refData, extraData, c.tag == "sub_table") // 等于sub table 需要删除
	if err != nil {
		//logger.Logger.Errorw("new  is  err is", err.Error())
	}
	return nil
}

func (c *common) subGet(ctx context.Context, isReplace bool) error {
	return nil
}

// new :component create options
func (c *common) new(ctx context.Context, refData *RefData, originalData *ExtraData) error {
	ids := make([]interface{}, 0)
	for _, _subTable := range refData.New {
		subTable, ok := _subTable.(map[string]interface{})
		if !ok {
			continue
		}
		opt := &OptionEntity{}
		err := initOptionEntity(ctx, subTable, opt)
		if err != nil {
			return err
		}

		universal := consensus.Universal{
			UserID:   c.userID,
			UserName: c.userName,
		}
		foundation := consensus.Foundation{
			AppID:   refData.AppID,
			TableID: refData.TableID,
			Method:  "create",
		}
		bus := new(consensus.Bus)
		bus.Universal = universal
		bus.Foundation = foundation
		bus.CreatedOrUpdate.Entity = opt.Entity
		bus.Ref.Ref = opt.Ref

		_, err = c.ref.Do(ctx, bus)
		if err != nil {
			return err
		}
		ids = append(ids, opt.Entity["_id"])

	}
	return c.addRelationShip(ctx, refData, originalData, ids)
}

func (c *common) addRelationShip(ctx context.Context, refData *RefData, extraData *ExtraData, ids []interface{}) error {

	primitiveKey, err := getPrimaryID(c.primaryEntity)
	if err != nil {
		return err
	}
	// 往中间表单插入数据
	for _, subID := range ids {
		entity := map[string]interface{}{
			primitiveID: primitiveKey,
			subIDs:      subID,
			fieldName:   c.key,
		}

		universal := consensus.Universal{
			UserID:   c.userID,
			UserName: c.userName,
		}
		foundation := consensus.Foundation{
			AppID:   refData.AppID,
			TableID: getRelationName(extraData.TableID, refData.TableID),
			Method:  "create",
		}
		bus := new(consensus.Bus)
		bus.Universal = universal
		bus.Foundation = foundation
		bus.CreatedOrUpdate.Entity = entity
		_, err = c.ref.Do(ctx, bus)
		if err != nil {
			//logger.Logger.Errorw("add sub form is err ,err is ", logger.STDRequestID(ctx), err.Error())
		}
	}
	return nil
}

// getRelationName getRelationName
func getRelationName(primary, sub string) string {
	return fmt.Sprintf("%s_%s", primary, sub)
}

func (c *common) update(ctx context.Context, refData *RefData) error {
	for _, _subTable := range refData.Updated {
		subTables, ok := _subTable.(map[string]interface{})
		if !ok {
			continue
		}
		opt := &OptionEntity{}
		err := initOptionEntity(ctx, subTables, opt)
		if err != nil {
			return err
		}
		universal := consensus.Universal{
			UserID:   c.userID,
			UserName: c.userName,
		}
		foundation := consensus.Foundation{
			AppID:   refData.AppID,
			TableID: refData.TableID,
			Method:  "update",
		}
		bus := new(consensus.Bus)
		bus.Universal = universal
		bus.Foundation = foundation
		bus.CreatedOrUpdate.Entity = opt.Entity
		bus.Ref.Ref = opt.Ref
		bus.Get.Query = opt.Query
		_, err = c.ref.Do(ctx, bus)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *common) delete(ctx context.Context, refData *RefData, originalData *ExtraData, isDelete bool) error {
	if len(refData.Deleted) == 0 {
		return nil
	}
	if isDelete {
		query := consensus.GetSimple(consensus.TermsKey, consensus.IDKey, refData.Deleted)
		universal := consensus.Universal{
			UserID:   c.userID,
			UserName: c.userName,
		}
		foundation := consensus.Foundation{
			AppID:   refData.AppID,
			TableID: refData.TableID,
			Method:  "delete",
		}
		bus := new(consensus.Bus)
		bus.Universal = universal
		bus.Foundation = foundation

		bus.Get.Query = query
		_, err := c.ref.Do(ctx, bus)
		if err != nil {
			return err
		}
	}
	subConditions := consensus.GetSimple(consensus.TermsKey, subIDs, refData.Deleted)
	fieldNameConditions := consensus.GetSimple(consensus.TermKey, fieldName, c.key)
	dslQuery := consensus.GetBool("must", subConditions, fieldNameConditions)

	universal := consensus.Universal{
		UserID:   c.userID,
		UserName: c.userName,
	}
	foundation := consensus.Foundation{
		AppID:   refData.AppID,
		TableID: getRelationName(originalData.TableID, refData.TableID),
		Method:  "delete",
	}
	bus := new(consensus.Bus)
	bus.Universal = universal
	bus.Foundation = foundation

	bus.Get.Query = dslQuery
	_, err := c.ref.Do(ctx, bus)
	if err != nil {
		//logger.Logger.Errorw("add sub form is err ,err is ", logger.STDRequestID(ctx), err.Error())
	}
	return nil
}

func (c *common) perInitData(ctx context.Context, refData *RefData, original *ExtraData) error {
	var err error
	if refData != nil {
		err = initRefData(ctx, c.refValue, refData)
	}

	if err != nil {
		return err
	}
	if original != nil {
		err = initExtraData(ctx, c.extraValue, original)
	}
	if err != nil {
		return err
	}
	return nil
}

func getPrimaryID(e consensus.Entity) (string, error) {
	if e == nil {
		return "", nil
	}
	value := reflect.ValueOf(e)
	switch _t := reflect.TypeOf(e); _t.Kind() {
	case reflect.Map:
		if values := value.MapIndex(reflect.ValueOf("_id")); values.IsValid() {
			if values.CanInterface() {
				return values.Interface().(string), nil
			}
		}
	}

	return "", nil
}

type serial struct {
	common
}

func (s *serial) getTag() string {
	return "serial"
}

func (s *serial) handlerFunc(ctx context.Context, action string) error {
	if action != "create" {
		return nil
	}
	originalData := &ExtraData{}
	err := s.perInitData(ctx, nil, originalData)
	if err != nil {
		return err
	}
	fmt.Printf("%+v\n", originalData)
	serialMap := s.ref.serialRepo.GetAll(ctx, originalData.AppID, originalData.TableID, s.key)
	serialScheme, res, err := utils.ExecuteTemplate(serialMap)
	if err != nil {
		//logger.Logger.Errorw("serial template is err", logger.STDRequestID(ctx))
		return err
	}
	//
	entity, ok := s.primaryEntity.(map[string]interface{})
	if !ok {
		return nil
	}
	entity[s.key] = res
	data, err := json.Marshal(serialScheme)
	if err != nil {
		return err
	}

	err = s.ref.serialRepo.Create(ctx, originalData.AppID, originalData.TableID, s.key, map[string]interface{}{
		redis.Serials: data,
	})
	if err != nil {
		return err
	}
	return nil
}

type associatedRecords struct {
	common
}

func (a *associatedRecords) getTag() string {
	return "associated_records"
}

func (a *associatedRecords) handlerFunc(ctx context.Context, action string) error {

	return nil
}
