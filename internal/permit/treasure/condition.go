package treasure

import (
	"context"
	"github.com/quanxiang-cloud/cabin/logger"
	"github.com/quanxiang-cloud/cabin/tailormade/header"
	"github.com/quanxiang-cloud/form/pkg/misc/client/lowcode"
	"reflect"

	"github.com/quanxiang-cloud/form/internal/permit"
	"github.com/quanxiang-cloud/form/internal/service/types"
	"github.com/quanxiang-cloud/form/pkg/misc/config"
)

const (
	_bool  = "bool"
	_terms = "terms"
	_match = "match"
)

// Condition condition.
type Condition struct {
	parsers   map[string]Parser
	searchAPI lowcode.SearchAPI
	formAPI   *lowcode.Form
}

// NewCondition new condition.
func NewCondition(conf *config.Config) *Condition {
	return &Condition{
		parsers:   make(map[string]Parser),
		searchAPI: lowcode.NewSearchAPI(conf),
		formAPI:   lowcode.NewForm(conf.InternalNet),
	}
}

// SetParse set parser.
func (c *Condition) SetParse(ctx context.Context, req *permit.Request) {
	for _, parse := range parsers {
		parse.Build(ctx, c, req)
		c.parsers[parse.Tag()] = parse
	}
}

var parsers = []Parser{
	&user{},
	&subordinate{},
	&project{},
}

// Parser parse param.
type Parser interface {
	Tag() string
	Build(context.Context, *Condition, *permit.Request)
	Parse(string, map[string]interface{}) error
}

type project struct {
	ctx  context.Context
	cond *Condition
	req  *permit.Request
}

func (p project) Tag() string {
	return "$project"
}

func (p *project) Build(ctx context.Context, cond *Condition, req *permit.Request) {
	p.ctx = ctx
	p.cond = cond
	p.req = req
}

func (p *project) getValue() ([]string, error) {
	logger.Logger.WithName("project user data").Infow("data",
		"req id ", header.GetRequestIDKV(p.ctx).Fuzzy()[1], "resp total ", p.req.UserID)
	resp, err := p.cond.formAPI.UserProject(p.ctx, p.req.UserID)

	if err != nil {
		return nil, err
	}

	ids := make([]string, resp.Total)
	for index, value := range resp.List {
		ids[index] = value.ProjectID
	}

	logger.Logger.WithName("project user data").Infow("data",
		"req id ", header.GetRequestIDKV(p.ctx).Fuzzy()[1], "resp total ", resp.Total, "ids", ids)
	return ids, nil
}

func (p *project) Parse(key string, params map[string]interface{}) error {
	value, err := p.getValue()
	if err != nil {
		return err
	}
	params[_terms] = types.M{
		"project_id": value,
	}
	delete(params, p.Tag())
	return nil
}

type user struct {
	ctx  context.Context
	cond *Condition
	req  *permit.Request
}

// Tag tag.
func (u *user) Tag() string {
	return "$user"
}

// Build build user.
func (u *user) Build(ctx context.Context, cond *Condition, req *permit.Request) {
	u.ctx = ctx
	u.cond = cond
	u.req = req
}

func (u *user) getValue() string {
	return u.req.UserID
}

// Parse parse tag.
func (u *user) Parse(key string, params map[string]interface{}) error {
	value := u.getValue()

	params[_match] = types.M{
		key: value,
	}
	delete(params, u.Tag())

	return nil
}

type subordinate struct {
	ctx  context.Context
	cond *Condition
	req  *permit.Request
}

// Tag tag.
func (s *subordinate) Tag() string {
	return "$subordinate"
}

// Build build.
func (s *subordinate) Build(ctx context.Context, cond *Condition, req *permit.Request) {
	s.ctx = ctx
	s.cond = cond
	s.req = req
}

func (s *subordinate) getValue() ([]string, error) {
	resp, err := s.cond.searchAPI.Subordinate(s.ctx, s.req.UserID)

	if err != nil {
		return nil, err
	}
	ids := make([]string, resp.Total)
	for index, value := range resp.Users {
		ids[index] = value.ID
	}
	logger.Logger.WithName("subordinate data").Infow("data",
		"req id ", header.GetRequestIDKV(s.ctx).Fuzzy()[1], "resp total ", resp.Total, "ids", ids)
	return ids, nil
}

// Parse parse param.
func (s *subordinate) Parse(key string, params map[string]interface{}) error {
	value, err := s.getValue()
	if err != nil {
		return err
	}

	params[_terms] = types.M{
		key: value,
	}
	delete(params, s.Tag())

	return nil
}

// ParseCondition parse param.
func (c *Condition) ParseCondition(condition interface{}) error {
	condType := reflect.TypeOf(condition)
	switch condType.Kind() {
	case reflect.Ptr:
		return c.ParseCondition(reflect.ValueOf(condition).Elem().Interface())
	case reflect.Map:
		condValue := reflect.ValueOf(condition)

		if !c.checkMapLen(condValue) {
			return nil
		}

		if key := condValue.MapKeys()[0]; key.String() == _bool {
			bool2 := condValue.MapIndex(key)
			if bool2.CanInterface() {
				return c.parseBool(bool2.Interface())
			}
		}

		return c.parse(condition)
	}

	return nil
}

func (c *Condition) parseBool(bool2 interface{}) error {
	paramType := reflect.TypeOf(bool2)
	switch paramType.Kind() {
	case reflect.Ptr:
		return c.parseBool(reflect.ValueOf(bool2).Elem().Interface())
	case reflect.Map:
		paramVal := reflect.ValueOf(bool2)

		if !c.checkMapLen(paramVal) {
			return nil
		}

		key := paramVal.MapKeys()[0]
		if !paramVal.MapIndex(key).CanInterface() {
			return nil
		}

		return c.parseParam(paramVal.MapIndex(key).Interface())
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

			err := c.parse(elemVal.Index(index).Interface())
			if err != nil {
				return err
			}
		}

		return nil
	}

	return nil
}

func (c *Condition) parse(elem interface{}) error {
	switch _parseType := reflect.TypeOf(elem); _parseType.Kind() {
	case reflect.Ptr:
		return c.parse(reflect.ValueOf(elem).Elem().Interface())
	case reflect.Map:
		parseVal := reflect.ValueOf(elem)

		if !c.checkMapLen(parseVal) {
			return nil
		}

		if key := parseVal.MapKeys()[0]; key.String() != _bool {
			data := parseVal.MapIndex(key)

			parser, ok := c.parsers[key.String()]
			if !ok {
				return nil
			}

			params, ok := elem.(map[string]interface{})
			if !ok {
				return nil
			}

			return parser.Parse(data.Elem().String(), params)
		}

		return c.ParseCondition(elem)
	}

	return nil
}

func (c *Condition) checkMapLen(value reflect.Value) bool {
	return len(value.MapKeys()) == 1
}
