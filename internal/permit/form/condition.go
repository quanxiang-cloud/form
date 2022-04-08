package guard

import (
	"context"

	"github.com/quanxiang-cloud/cabin/logger"
	"github.com/quanxiang-cloud/cabin/tailormade/header"
	"github.com/quanxiang-cloud/form/internal/permit"
	"github.com/quanxiang-cloud/form/internal/permit/treasure"
	"github.com/quanxiang-cloud/form/internal/service/types"
	"github.com/quanxiang-cloud/form/pkg/misc/config"
)

const (
	_query     = "query"
	_condition = "condition"
	_bool      = "bool"
	_must      = "must"
)

// Condition is a guard for permit.
type Condition struct {
	cond *treasure.Condition
	next permit.Permit
}

// NewCondition returns a new guard for permit.
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

// Do is a guard for permit.
func (c *Condition) Do(ctx context.Context, req *permit.Request) (*permit.Response, error) {
	query := req.Body[_query]

	err := c.cond.SetParseValue(ctx, req)
	if err != nil {
		logger.Logger.WithName("form condition").Errorw(err.Error(), header.GetRequestIDKV(ctx).Fuzzy()...)
		return nil, err
	}

	dataes := make([]interface{}, 0, 2)
	if query != nil {
		dataes = append(dataes, query)
	}
	condition := req.Permit.Condition
	if condition != nil {
		err = c.cond.ParseCondition(condition)
		if err != nil {
			logger.Logger.WithName("form condition").Errorw(err.Error(), header.GetRequestIDKV(ctx).Fuzzy()...)
			return nil, err
		}
		dataes = append(dataes, condition)
	}

	if len(dataes) != 0 {
		req.Body[_query] = permit.Body{
			_bool: types.M{
				_must: dataes,
			},
		}
	}

	return c.next.Do(ctx, req)
}
