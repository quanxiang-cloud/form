package condition

import (
	"context"
	"reflect"

	"github.com/quanxiang-cloud/form/internal/service/types"
)

type Condition struct {
	parsers map[string]Parser
}

func NewCondition() *Condition {
	return &Condition{
		parsers: make(map[string]Parser),
	}
}

type CondReq struct {
	UserID   string `json:"userID"`
	BodyData map[string]interface{}
}

// Do Do
func (c *Condition) Do(ctx context.Context, req *CondReq) error {
	err := c.SetParsers(ctx, req)
	if err != nil {
		return err
	}

	dataes := make([]interface{}, 0, 2)
	if val, ok := req.BodyData["query"]; ok {
		dataes = append(dataes, val)
	}

	condition := req.BodyData["condition"]
	err = c.parseCondition(condition)
	if err != nil {
		return err
	}

	if condition != nil {
		dataes = append(dataes, condition)
	}

	query := types.Query{
		"bool": types.M{
			"must": dataes,
		},
	}

	req.BodyData["query"] = query
	return nil
}

func (c *Condition) SetParsers(ctx context.Context, req *CondReq) error {
	for _, parse := range parsers {
		err := parse.SetValue(ctx, c, req)
		if err != nil {
			return err
		}

		c.parsers[parse.GetTag()] = parse
	}
	return nil
}

var parsers = []Parser{
	&user{},
	&subordinate{},
}

type Parser interface {
	GetTag() string
	SetValue(context.Context, *Condition, *CondReq) error
	Parse(string, interface{})
}

type user struct {
	value interface{}
}

func (u *user) GetTag() string {
	return "$user"
}

func (u *user) SetValue(ctx context.Context, c *Condition, req *CondReq) error {
	u.value = req.UserID
	return nil
}

func (u *user) Parse(key string, valueSet interface{}) {
	m := valueSet.(map[string]interface{})
	m["match"] = types.M{
		key: u.value,
	}
	delete(m, u.GetTag())
}

type subordinate struct {
	value interface{}
}

func (s *subordinate) GetTag() string {
	return "$subordinate"
}

func (s *subordinate) SetValue(ctx context.Context, c *Condition, req *CondReq) error {
	// TODO set subordinate value
	return nil
}

func (s *subordinate) Parse(key string, valueSet interface{}) {
	m := valueSet.(map[string]interface{})

	m["terms"] = types.M{
		key: s.value,
	}

	delete(m, s.GetTag())
}

func (c *Condition) parseCondition(condition interface{}) error {
	if condition == nil {
		return nil
	}

	switch condType := reflect.TypeOf(condition); condType.Kind() {
	case reflect.Ptr:
		c.parseCondition(reflect.ValueOf(condition).Elem().Interface())
	case reflect.Map:
		condValue := reflect.ValueOf(condition)
		if len(condValue.MapKeys()) == 0 {
			return nil
		}

		bool2 := condValue.MapIndex(condValue.MapKeys()[0])
		if !bool2.CanInterface() {
			return nil
		}

		return c.parseBool(bool2)
	}
	return nil
}

func (c *Condition) parseBool(bool2 reflect.Value) error {
	switch paramType := reflect.TypeOf(bool2.Interface()); paramType.Kind() {
	case reflect.Map:
		paramVal := reflect.ValueOf(bool2.Interface())

		keys := paramVal.MapKeys()

		for _, key := range keys {

			if !paramVal.MapIndex(key).CanInterface() {
				return nil
			}

			err := c.parseParam(paramVal, key)
			if err != nil {
				return err
			}

		}
	}
	return nil
}

func (c *Condition) parseParam(paramVal, key reflect.Value) error {
	switch elemType := reflect.TypeOf(paramVal.MapIndex(key).Interface()); elemType.Kind() {
	case reflect.Slice, reflect.Array:
		elemVal := reflect.ValueOf(paramVal.MapIndex(key).Interface())
		for index := 0; index < elemVal.Len(); index++ {
			if !elemVal.Index(index).CanInterface() {
				return nil
			}

			err := c.parse(elemVal, index)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (c *Condition) parse(elemVal reflect.Value, index int) error {
	switch _parseType := reflect.TypeOf(elemVal.Index(index).Interface()); _parseType.Kind() {
	case reflect.Map:
		parseVal := reflect.ValueOf(elemVal.Index(index).Interface())

		if len(parseVal.MapKeys()) == 0 {
			return nil
		}

		parseKey := parseVal.MapKeys()[0]
		if parseKey.String() == "bool" {
			return c.parseCondition(elemVal.Index(index).Interface())
		}

		data := parseVal.MapIndex(parseKey)

		parser, ok := c.parsers[parseKey.String()]
		if !ok {
			return nil
		}
		parser.Parse(data.Elem().String(), parseVal.Interface())
	}
	return nil
}
