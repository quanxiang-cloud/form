package condition

import (
	"context"
	"reflect"

	"github.com/quanxiang-cloud/form/internal/service/types"
)

const (
	_query     = "query"
	_condition = "condition"
	_bool      = "bool"
	_must      = "must"
	_terms     = "terms"
	_match     = "match"
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
	if val, ok := req.BodyData[_query]; ok {
		dataes = append(dataes, val)
	}

	condition := req.BodyData[_condition]
	err = c.parseCondition(condition)
	if err != nil {
		return err
	}

	if condition != nil {
		dataes = append(dataes, condition)
	}

	query := types.Query{
		_bool: types.M{
			_must: dataes,
		},
	}

	req.BodyData[_query] = query
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
	m, ok := valueSet.(map[string]interface{})
	if ok {
		m[_match] = types.M{
			key: u.value,
		}
		delete(m, u.GetTag())
	}
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
	m, ok := valueSet.(map[string]interface{})
	if ok {
		m[_terms] = types.M{
			key: s.value,
		}

		delete(m, s.GetTag())
	}
}

func (c *Condition) parseCondition(condition interface{}) error {
	if condition == nil {
		return nil
	}

	switch condType := reflect.TypeOf(condition); condType.Kind() {
	case reflect.Ptr:
		return c.parseCondition(reflect.ValueOf(condition).Elem().Interface())
	case reflect.Map:
		condValue := reflect.ValueOf(condition)
		if len(condValue.MapKeys()) == 0 {
			return nil
		}

		bool2 := condValue.MapIndex(condValue.MapKeys()[0])
		if !bool2.CanInterface() {
			return nil
		}

		return c.parseBool(bool2.Interface())
	}
	return nil
}

func (c *Condition) parseBool(bool2 interface{}) error {
	switch paramType := reflect.TypeOf(bool2); paramType.Kind() {
	case reflect.Ptr:
		return c.parseBool(reflect.ValueOf(bool2).Elem().Interface())
	case reflect.Map:
		paramVal := reflect.ValueOf(bool2)
		for _, key := range paramVal.MapKeys() {

			if !paramVal.MapIndex(key).CanInterface() {
				return nil
			}

			param := paramVal.MapIndex(key).Interface()
			err := c.parseParam(param)
			if err != nil {
				return err
			}

		}
	}
	return nil
}

func (c *Condition) parseParam(param interface{}) error {
	switch elemType := reflect.TypeOf(param); elemType.Kind() {
	case reflect.Ptr:
		return c.parseParam(reflect.ValueOf(param).Elem().Interface())
	case reflect.Slice, reflect.Array:
		elemVal := reflect.ValueOf(param)
		for index := 0; index < elemVal.Len(); index++ {
			if !elemVal.Index(index).CanInterface() {
				return nil
			}

			elem := elemVal.Index(index).Interface()
			err := c.parse(elem)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (c *Condition) parse(elem interface{}) error {
	switch _parseType := reflect.TypeOf(elem); _parseType.Kind() {
	case reflect.Ptr:
		return c.parse(reflect.ValueOf(elem).Elem().Interface())
	case reflect.Map:
		parseVal := reflect.ValueOf(elem)

		if len(parseVal.MapKeys()) == 0 {
			return nil
		}

		parseKey := parseVal.MapKeys()[0]
		if parseKey.String() == _bool {
			return c.parseCondition(elem)
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
