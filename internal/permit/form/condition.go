package guard

import (
	"context"
	"github.com/quanxiang-cloud/form/internal/permit"
	"github.com/quanxiang-cloud/form/internal/permit/treasure"
	"github.com/quanxiang-cloud/form/pkg/misc/config"
)

const (
	_query     = "query"
	_condition = "condition"
	_bool      = "bool"
	_must      = "must"
)

type Condition struct {
	cond *treasure.Condition
	next permit.Permit
}

func NewCondition(conf *config.Config) (*Condition, error) {
	next, err := NewProxy(conf)
	if err != nil {
		return nil, err
	}
	return &Condition{
		cond: treasure.NewCondition(),
		next: next,
	}, nil
}

func (c *Condition) Do(ctx context.Context, req *permit.Request) (*permit.Response, error) {
	//var (
	//	query     = req.Body[_query]
	//	condition = req.Body[_condition]
	//)
	//
	//err := c.cond.SetParseValue(ctx, req)
	//if err != nil {
	//	return nil, err
	//}
	//
	//dataes := make([]interface{}, 0, 2)
	//if query == nil {
	//	return nil, fmt.Errorf("must have query dsl")
	//}
	//dataes = append(dataes, query)
	//
	//if condition != nil {
	//	err = c.cond.ParseCondition(condition)
	//	if err != nil {
	//		return nil, err
	//	}
	//	dataes = append(dataes, condition)
	//}
	//
	//req.Body[_query] = permit.Body{
	//	_bool: types.M{
	//		_must: dataes,
	//	},
	//}

	return c.next.Do(ctx, req)
}
