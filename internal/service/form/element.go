package form

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/quanxiang-cloud/cabin/logger"
	"github.com/quanxiang-cloud/cabin/tailormade/header"
	"github.com/quanxiang-cloud/form/internal/models"
	"reflect"

	"github.com/quanxiang-cloud/form/internal/models/redis"
	"github.com/quanxiang-cloud/form/internal/service/consensus"
	"github.com/quanxiang-cloud/form/internal/service/types"
	"github.com/quanxiang-cloud/form/pkg/misc/utils"
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
	err = c.update(ctx, refData)
	if err != nil {
		logger.Logger.Errorw(err.Error(), header.GetRequestIDKV(ctx).Fuzzy()...)
	}
	err = c.new(ctx, refData, extraData)
	if err != nil {
		logger.Logger.Errorw(err.Error(), header.GetRequestIDKV(ctx).Fuzzy()...)
	}
	err = c.delete(ctx, refData, extraData, c.tag == "sub_table") // 等于sub table 需要删除
	if err != nil {
		logger.Logger.Errorw(err.Error(), header.GetRequestIDKV(ctx).Fuzzy()...)
	}
	return nil
}

func (c *common) subGet(ctx context.Context, isReplace bool) error {
	extraData := &ExtraData{}
	refData := &RefData{}
	err := c.perInitData(ctx, refData, extraData)
	if err != nil {
		return err
	}
	data := make([]interface{}, 0)
	id, err := getPrimaryID(c.primaryEntity)
	if err != nil {
		return err
	}
	idConditions := consensus.GetSimple(consensus.TermsKey, primitiveID, id)
	keyCondition := consensus.GetSimple(consensus.TermKey, fieldName, c.key)
	boolQuery := consensus.GetBool(consensus.Must, idConditions, keyCondition)

	universal := consensus.Universal{
		UserID:   c.userID,
		UserName: c.userName,
	}
	foundation := consensus.Foundation{
		AppID:   refData.AppID,
		TableID: getRelationName(extraData.TableID, refData.TableID),
		Method:  "search",
	}
	bus := new(consensus.Bus)
	bus.Universal = universal
	bus.Foundation = foundation
	bus.Get.Query = boolQuery
	list := consensus.List{
		Size: 1000,
		Page: 1,
		Sort: []string{"created_at"},
	}
	bus.List = list

	searchResp1, err := c.ref.Do(ctx, bus)
	if err != nil {
		return err
	}
	for _, value := range searchResp1.Entities {
		_, ok := value[subIDs]
		if !ok {
			continue
		}
		data = append(data, value[subIDs])
	}
	if !isReplace {
		setValue(c.primaryEntity, c.key, data)
		return nil
	}

	idsQuery := consensus.GetSimple(consensus.TermsKey, "_id", data)
	bus1 := new(consensus.Bus)
	bus1.Universal = universal
	bus1.Foundation = consensus.Foundation{
		TableID: refData.TableID,
		AppID:   refData.AppID,
		Method:  "search",
	}
	bus1.Get.Query = idsQuery
	bus.List = list

	subResp, err := c.ref.Do(ctx, bus1)
	if err != nil {
		return err
	}
	err = c.findOnePost(ctx, &params{
		ctx:       ctx,
		subResp:   subResp,
		refData:   refData,
		extraData: extraData,
		data:      data,
	})
	if err != nil {
		return err
	}
	return nil

}

type params struct {
	ctx       context.Context
	subResp   *consensus.Response
	refData   *RefData
	extraData *ExtraData
	data      []interface{}
}

// sub

func (c *common) findOnePost(ctx context.Context, param *params) error {
	replaceData := make([]interface{}, 0)
	relation, _, err := c.ref.relationRepo.List(c.ref.db, &models.TableRelationQuery{
		AppID:        param.refData.AppID,
		TableID:      param.extraData.TableID,
		FieldName:    c.key,
		SubTableType: c.tag,
		SubTableID:   param.refData.TableID,
	}, 1, 10)

	if err != nil {
		return err
	}

	if len(relation) != 1 {
		logger.Logger.Infow("relation  is not one ", header.GetRequestIDKV(ctx).Fuzzy()...)
		return nil
	}

	mapID := make(map[string]types.M)
	for i := 0; i < len(param.subResp.Entities); i++ {
		e1 := param.subResp.Entities[i]
		id := e1["_id"].(string)
		propertiesFilter(e1, relation[0].Filter)
		mapID[id] = e1
	}
	for _, value := range param.data {
		id, ok := value.(string)
		if !ok {
			return errors.New("par is error ")
		}
		if e, ok := mapID[id]; ok {
			replaceData = append(replaceData, e)
		}
	}
	setValue(c.primaryEntity, c.key, replaceData)
	return nil
}

// PropertiesFilter PropertiesFilter
func propertiesFilter(oldProperties map[string]interface{}, filter []string) {
	if filter == nil {
		return
	}
	// 1、 把过滤字段 切换成map，
	filters := make(map[string]int, len(filter))
	for index, filterKey := range filter {
		filters[filterKey] = index
	}
	// 2、判断 遍历的列，在不在map 中，不在，删除该列
	for column := range oldProperties {
		if _, ok := filters[column]; !ok {
			delete(oldProperties, column)
		}
	}
}

// new :component create options.
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
			logger.Logger.Errorw(err.Error(), header.GetRequestIDKV(ctx).Fuzzy()...)
		}
	}
	return nil
}

// getRelationName getRelationName.
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
	pid := consensus.GetSimple(consensus.TermsKey, "primitiveID", originalData.ID)
	dslQuery := consensus.GetBool("must", subConditions, fieldNameConditions, pid)

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
		logger.Logger.Errorw(err.Error(), header.GetRequestIDKV(ctx).Fuzzy()...)
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
		} else {
			return "", errors.New("miss id")
		}
	}

	return "", nil
}

func setValue(e consensus.Entity, key string, values interface{}) {
	if e == nil {
		return
	}
	value := reflect.ValueOf(e)
	switch _t := reflect.TypeOf(e); _t.Kind() {
	case reflect.Map:
		value.SetMapIndex(reflect.ValueOf(key), reflect.ValueOf(values))
	}
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
	serialMap := s.ref.serialRepo.GetAll(ctx, originalData.AppID, originalData.TableID, s.key)
	serialScheme, res, err := utils.ExecuteTemplate(serialMap)
	if err != nil {
		// logger.Logger.Errorw("serial template is err", logger.STDRequestID(ctx))
		return err
	}

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

	switch action {
	case "create":
		extraData := &ExtraData{}
		refData := &RefData{}
		err := a.perInitData(ctx, refData, extraData)
		if err != nil {
			return err
		}
		err = a.common.addRelationShip(ctx, refData, extraData, refData.New)
		if err != nil {
			logger.Logger.Errorw(err.Error()+"add releation ship is nil ", header.GetRequestIDKV(ctx).Fuzzy()...)
		}
	case "update":
		extraData := &ExtraData{}
		refData := &RefData{}
		err := a.perInitData(ctx, refData, extraData)
		if err != nil {
			return err
		}
		err = a.common.addRelationShip(ctx, refData, extraData, refData.New)
		if err != nil {
			logger.Logger.Errorw(err.Error()+"add releation ship is nil ", header.GetRequestIDKV(ctx).Fuzzy()...)
		}
		err = a.common.delete(ctx, refData, extraData, false)
		if err != nil {
			logger.Logger.Errorw(err.Error()+"delete releation ship is nil ", header.GetRequestIDKV(ctx).Fuzzy()...)
		}
	case "get":
		a.common.subGet(ctx, true)
	}
	return nil
}

// 聚合组件。
type aggregation struct {
	common
}

func (a *aggregation) getTag() string {
	return "aggregation"
}

func (a *aggregation) GetTag() string {
	return "aggregation"
}

func (a *aggregation) handlerFunc(ctx context.Context, action string) error {
	if action != "get" {
		return nil
	}
	refData := &RefData{}
	extra := &ExtraData{}
	err := a.perInitData(ctx, refData, extra)
	if err != nil {
		return err
	}
	if extra.ID == "" {

	}
	data := make([]interface{}, 0)
	idCondition := consensus.GetSimple(consensus.TermsKey, "primitiveID", extra.ID)
	fieldNameCondition := consensus.GetSimple(consensus.TermKey, "fieldName", refData.SourceFieldID)
	dslQuery := consensus.GetBool("must", idCondition, fieldNameCondition)
	foundation := consensus.Foundation{
		AppID:   refData.AppID,
		TableID: getRelationName(extra.TableID, refData.TableID),
		Method:  "search",
	}
	bus := new(consensus.Bus)
	bus.Foundation = foundation
	bus.Get.Query = dslQuery
	list := consensus.List{
		Size: 1000,
		Page: 1,
		Sort: []string{"created_at"},
	}
	bus.List = list

	searchResp1, err := a.ref.Do(ctx, bus)
	if err != nil {
		return err
	}
	for _, value := range searchResp1.Entities {
		_, ok := value[subIDs]
		if !ok {
			continue
		}
		data = append(data, value[subIDs])
	}
	if refData.AggType == "sum" {
		agg, err := a.doAgg(ctx, "avg", refData, data)
		if err != nil {
			return err
		}
		if agg == nil {
			setValue(a.primaryEntity, a.key, 112)
			return nil
		}
	}
	agg, err := a.doAgg(ctx, refData.AggType, refData, data)
	if err != nil {
		return err
	}
	setValue(a.primaryEntity, a.key, agg)
	return nil
}

func (a *aggregation) doAgg(ctx context.Context, aggType string, refData *RefData, data interface{}) (interface{}, error) {
	constructQuery(refData, data)
	alias := refData.AggType + refData.FieldName
	agg := map[string]interface{}{
		alias: consensus.KeyValue{
			aggType: consensus.KeyValue{
				"field": refData.FieldName,
			},
		},
	}
	foundation := consensus.Foundation{
		AppID:   refData.AppID,
		TableID: refData.TableID,
		Method:  "search",
	}
	bus := new(consensus.Bus)
	bus.Foundation = foundation
	bus.Get.Query = refData.Query
	bus.Aggs = agg
	searchResp, err := a.ref.Do(ctx, bus)
	if err != nil {
		return nil, err
	}
	return getResult(searchResp.Entities, alias), nil
}

func getResult(data interface{}, fieldName string) interface{} {
	v := reflect.ValueOf(data)
	switch reflect.TypeOf(data).Kind() {
	case reflect.Slice, reflect.Array:
		if v.Len() >= 1 {
			return getResult(v.Index(0).Interface(), fieldName)
		}
		return nil
	case reflect.Map:
		if value := v.MapIndex(reflect.ValueOf(fieldName)); value.IsValid() {
			return value.Interface()
		}
	}
	return nil
}

// AggregationQuery AggregationQuery
func aggregationQuery(data interface{}, ids interface{}) {
	if data == nil {
		return
	}
	switch reflect.TypeOf(data).Kind() {
	case reflect.Map:
		v := reflect.ValueOf(data)
		if value := v.MapIndex(reflect.ValueOf("bool")); value.IsValid() {
			aggregationQuery(value.Elem().Interface(), ids)
			return
		}
		if value := v.MapIndex(reflect.ValueOf("must")); value.IsValid() {
			must, ok := value.Interface().([]interface{})
			if !ok {
				return
			}
			termsIDs := consensus.GetSimple(consensus.TermsKey, consensus.IDKey, ids)
			must = append(must, termsIDs)
			v.SetMapIndex(reflect.ValueOf("must"), reflect.ValueOf(must))
			return
		}
		if value := v.MapIndex(reflect.ValueOf("should")); value.IsValid() {
			shouldArr, ok := value.Interface().([]interface{})
			if !ok {
				return
			}
			termsIDs := consensus.GetSimple(consensus.TermsKey, consensus.IDKey, ids)
			arr := make([]interface{}, 0)
			arr = append(arr, consensus.KeyValue{
				"should": shouldArr,
			}, termsIDs)

			v.SetMapIndex(reflect.ValueOf("should"), reflect.Value{})
			v.SetMapIndex(reflect.ValueOf("must"), reflect.ValueOf(arr))
			return
		}
	}
}

func constructQuery(ref *RefData, ids interface{}) {
	if ref.Query == nil {
		ref.Query = consensus.GetSimple(consensus.TermsKey, consensus.IDKey, ids)
	} else {
		aggregationQuery(ref.Query, ids)
	}
}
