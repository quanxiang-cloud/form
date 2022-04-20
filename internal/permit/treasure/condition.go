package treasure

import (
	"context"
	"fmt"
	"github.com/quanxiang-cloud/form/pkg/misc/client"
	"github.com/quanxiang-cloud/form/pkg/misc/config"
	"reflect"

	"github.com/quanxiang-cloud/form/internal/permit"
	"github.com/quanxiang-cloud/form/internal/service/types"
)

const (
	_bool  = "bool"
	_terms = "terms"
	_match = "match"
)

type Condition struct {
	parsers   map[string]Parser
	searchAPI client.SearchAPI
}

func NewCondition(config *config.Config) *Condition {
	return &Condition{
		parsers:   make(map[string]Parser),
		searchAPI: client.NewSearchAPI(config),
	}
}

func (c *Condition) SetParseValue(ctx context.Context, req *permit.Request) error {
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
	SetValue(context.Context, *Condition, *permit.Request) error
	Parse(string, interface{})
}

type user struct {
	value interface{}
}

func (u *user) GetTag() string {
	return "$user"
}

func (u *user) SetValue(ctx context.Context, c *Condition, req *permit.Request) error {
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

func (s *subordinate) SetValue(ctx context.Context, c *Condition, req *permit.Request) error {
	// TODO set subordinate value
	resp, err := c.searchAPI.Subordinate(ctx, req.UserID)
	if err != nil {
		return err
	}
	ids := make([]string, resp.Total)
	for index, value := range resp.Users {
		ids[index] = value.ID
	}
	s.value = ids
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

func (c *Condition) ParseCondition(condition interface{}) error {
	if condition == nil {
		return nil
	}

	switch condType := reflect.TypeOf(condition); condType.Kind() {
	case reflect.Ptr:
		return c.ParseCondition(reflect.ValueOf(condition).Elem().Interface())
	case reflect.Map:
		condValue := reflect.ValueOf(condition)
		if len(condValue.MapKeys()) == 0 {
			return nil
		}
		if !condValue.CanInterface() {
			return nil
		}

		for _, key := range condValue.MapKeys() {
			fmt.Println(key.String())
			if key.String() != _bool {
				err := c.parse(condValue.Interface())
				if err != nil {
					return err
				}
			} else {
				bool2 := condValue.MapIndex(key)
				if !bool2.CanInterface() {
					return nil
				}
				err := c.parseBool(bool2.Interface())
				if err != nil {
					return err
				}
			}

		}
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

		for _, key := range parseVal.MapKeys() {
			fmt.Println(key.String())
			if key.String() != _bool {
				data := parseVal.MapIndex(key)

				parser, ok := c.parsers[key.String()]
				if !ok {
					return nil
				}
				parser.Parse(data.Elem().String(), parseVal.Interface())
			} else {
				if err := c.ParseCondition(elem); err != nil {
					return err
				}
			}
		}

	}
	return nil
}
