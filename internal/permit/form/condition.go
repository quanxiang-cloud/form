package guard

import (
	"context"
	"encoding/json"
	"github.com/quanxiang-cloud/form/internal/models"
	"net/http"
	"net/url"

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
	if req.Permit.Types == models.InitType {
		return c.next.Do(ctx, req)
	}
	var query permit.Object
	switch req.Request.Method {
	case http.MethodGet:
		query = req.Query
	case http.MethodPost:
		bytes, err := json.Marshal(req.Body[_query])
		if err != nil {
			return nil, err
		}
		json.Unmarshal(bytes, &query)
	}

	err := c.cond.SetParseValue(ctx, req)
	if err != nil {
		logger.Logger.WithName("form condition").Errorw(err.Error(), header.GetRequestIDKV(ctx).Fuzzy()...)
		return nil, err
	}

	dataes := make([]interface{}, 0, 2)
	if query != nil && len(query) != 0 {
		dataes = append(dataes, query)
	}
	condition := req.Permit.Condition
	println(len(condition))
	if condition != nil && len(condition) != 0 {
		err = c.cond.ParseCondition(condition)
		if err != nil {
			logger.Logger.WithName("form condition").Errorw(err.Error(), header.GetRequestIDKV(ctx).Fuzzy()...)
			return nil, err
		}
		dataes = append(dataes, condition)
	}

	var newQuery permit.Object
	if len(dataes) != 0 {
		newQuery = permit.Object{
			_bool: types.M{
				_must: dataes,
			},
		}
	}
	switch req.Request.Method {
	case http.MethodGet:
		queryBytes, err := json.Marshal(newQuery)
		if err != nil {
			return nil, err
		}

		str, err := url.QueryUnescape(req.Request.URL.RawQuery)
		if err != nil {
			return nil, err
		}

		v, err := url.ParseQuery(str)
		if err != nil {
			return nil, err
		}

		v.Set("query", string(queryBytes))

		req.Request.URL.RawQuery = v.Encode()
	case http.MethodPost:
		req.Body[_query] = newQuery
	}

	return c.next.Do(ctx, req)
}
