package form

// import (
// 	"context"
// 	"encoding/json"
// 	"github.com/quanxiang-cloud/form/internal/models"
// 	"github.com/quanxiang-cloud/form/internal/service/consensus"
// 	"github.com/quanxiang-cloud/form/internal/service/types"
// 	"github.com/quanxiang-cloud/form/pkg/misc/config"
// )

// type condition struct {
// 	next    consensus.Guidance
// 	parsers map[string]Parser
// }

// func NewCondition(conf *config.Config) (consensus.Guidance, error) {
// 	newRefs, err := NewRefs(conf)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &condition{
// 		next:    newRefs,
// 		parsers: make(map[string]Parser),
// 	}, nil
// }

// // Do Do
// func (c *condition) Do(ctx context.Context, bus *consensus.Bus) (*consensus.Response, error) {
// 	//err := c.SetParsers(ctx, bus)
// 	//if err != nil {
// 	//	return nil, err
// 	//}
// 	//dataes := make([]interface{}, 0, 2)
// 	//if bus.Get.Query != nil {
// 	//	dataes = append(dataes, bus.Get.Query)
// 	//}
// 	//
// 	//err = c.parse(bus.Permit.Condition)
// 	//if err != nil {
// 	//	return nil, err
// 	//}
// 	//
// 	//if bus.Permit.Condition != nil {
// 	//	dataes = append(dataes, bus.Permit.Condition)
// 	//}
// 	//
// 	//query := types.Query{
// 	//	"bool": types.M{
// 	//		"must": dataes,
// 	//	},
// 	//}
// 	//bus.Get.Query = query

// 	return c.next.Do(ctx, bus)
// }

// func (c *condition) SetParsers(ctx context.Context, bus *consensus.Bus) error {
// 	for _, parse := range parsers {
// 		err := parse.SetValue(ctx, c, bus)
// 		if err != nil {
// 			return err
// 		}

// 		c.parsers[parse.GetTag()] = parse
// 	}
// 	return nil
// }

// var parsers = []Parser{
// 	&user{},
// 	&subordinate{},
// }

// type Parser interface {
// 	GetTag() string
// 	SetValue(context.Context, *condition, *consensus.Bus) error
// 	Parse(map[string]interface{}, string)
// }

// type user struct {
// 	value interface{}
// }

// func (u *user) GetTag() string {
// 	return "$user"
// }

// func (u *user) SetValue(ctx context.Context, c *condition, bus *consensus.Bus) error {
// 	u.value = bus.UserID
// 	return nil
// }

// func (u *user) Parse(must map[string]interface{}, key string) {
// 	must["match"] = types.M{
// 		key: u.value,
// 	}
// 	delete(must, u.GetTag())
// }

// type subordinate struct {
// 	value interface{}
// }

// func (s *subordinate) GetTag() string {
// 	return "$subordinate"
// }

// func (s *subordinate) SetValue(ctx context.Context, c *condition, bus *consensus.Bus) error {
// 	// TODO set subordinate value
// 	return nil
// }

// func (s *subordinate) Parse(must map[string]interface{}, key string) {
// 	must["terms"] = types.M{
// 		key: s.value,
// 	}

// 	delete(must, s.GetTag())
// }

// func (c *condition) parse(cond *models.Condition) error {
// 	if cond == nil {
// 		return nil
// 	}

// 	for _, must := range cond.Bool.Must {
// 		for tag, value := range must {

// 			bool2 := models.BOOL{}
// 			if tag == "bool" {
// 				boolBytes, err := json.Marshal(value)
// 				if err != nil {
// 					return err
// 				}

// 				err = json.Unmarshal(boolBytes, &bool2)
// 				if err != nil {
// 					return err
// 				}

// 				err = c.parse(&models.Condition{Bool: bool2})
// 				if err != nil {
// 					return err
// 				}
// 			} else {
// 				parser, ok := c.parsers[tag]
// 				if !ok {
// 					continue
// 				}
// 				parser.Parse(must, value.(string))
// 			}
// 		}
// 	}

// 	return nil
// }
