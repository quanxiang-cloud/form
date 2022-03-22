package condition

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"

	"github.com/quanxiang-cloud/form/internal/permit"
	"github.com/quanxiang-cloud/form/internal/permit/proxy"
	"github.com/quanxiang-cloud/form/internal/service/types"
	"github.com/quanxiang-cloud/form/pkg/misc/config"
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
	next    permit.Form
	parsers map[string]Parser
}

func NewCondition(conf *config.Config) (*Condition, error) {
	next, err := proxy.NewProxy(conf)
	if err != nil {
		return nil, err
	}
	return &Condition{
		parsers: make(map[string]Parser),
		next:    next,
	}, nil
}

func (c *Condition) Guard(ctx context.Context, req *permit.GuardReq) (*permit.GuardResp, error) {
	var (
		query     = req.Body[_query]
		condition = req.Body[_condition]
	)

	if req.Request.Method == http.MethodGet {
		query = req.Get.Query
		condition = req.Get.Condition
	}

	err := c.setParseValue(ctx, req)
	if err != nil {
		return nil, err
	}

	dataes := make([]interface{}, 0, 2)
	if query != nil {
		dataes = append(dataes, query)
	}

	if condition != nil {
		err = c.parseCondition(condition)
		if err != nil {
			return nil, err
		}
		dataes = append(dataes, condition)
	}

	newQuery := permit.Query{
		_bool: types.M{
			_must: dataes,
		},
	}

	b, _ := json.Marshal(newQuery)
	fmt.Println(string(b))

	if req.Request.Method == http.MethodGet {
		req.Get.Query = newQuery
	} else {
		req.Body[_query] = newQuery
	}

	return c.next.Guard(ctx, req)
}

func (c *Condition) setParseValue(ctx context.Context, req *permit.GuardReq) error {
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
	SetValue(context.Context, *Condition, *permit.GuardReq) error
	Parse(string, interface{})
}

type user struct {
	value interface{}
}

func (u *user) GetTag() string {
	return "$user"
}

func (u *user) SetValue(ctx context.Context, c *Condition, req *permit.GuardReq) error {
	u.value = req.Header.UserID
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

func (s *subordinate) SetValue(ctx context.Context, c *Condition, req *permit.GuardReq) error {
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

		// parseKey := parseVal.MapKeys()[0]
		// if parseKey.String() == _bool {
		// 	return c.parseCondition(elem)
		// }

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
				if err := c.parseCondition(elem); err != nil {
					return err
				}
			}
		}

	}
	return nil
}
