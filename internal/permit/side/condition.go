package side

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/quanxiang-cloud/cabin/logger"
	"github.com/quanxiang-cloud/cabin/tailormade/header"
	"github.com/quanxiang-cloud/form/internal/models"
	"github.com/quanxiang-cloud/form/internal/permit"
	"github.com/quanxiang-cloud/form/internal/permit/treasure"
	"github.com/quanxiang-cloud/form/internal/service/types"
	"github.com/quanxiang-cloud/form/pkg/misc/config"
)

const (
	_query = "query"
	_bool  = "bool"
	_must  = "must"
)

// Condition is a guard for permit.
type Condition struct {
	cond *treasure.Condition
	next permit.Permit
}

// NewCondition returns a new guard for permit.
func NewCondition(conf *config.Config, rawurl string) (*Condition, error) {
	next, err := NewProxy(conf, rawurl)
	if err != nil {
		return nil, err
	}

	return &Condition{
		cond: treasure.NewCondition(conf),
		next: next,
	}, nil
}

// Do is a guard for permit.
func (c *Condition) Do(ctx context.Context, req *permit.Request) (*permit.Response, error) {

	if req.Permit.Types == models.InitType {
		return c.next.Do(ctx, req)
	}

	// 判单来源。
	// TODO
	oldQuery := req.Data[_query]
	var query permit.Object
	switch req.Echo.Request().Method {
	case http.MethodGet:
		//query = req.Query
	case http.MethodPost:
		bytes, err := json.Marshal(oldQuery)
		if err != nil {
			return nil, err
		}
		json.Unmarshal(bytes, &query)
	}
	c.cond.SetParse(ctx, req)
	dataes := make([]interface{}, 0, 2)
	if query != nil && len(query) != 0 {
		dataes = append(dataes, query)
	}
	condition := req.Permit.Condition
	i, ok := condition[_query]
	if !ok {
		logger.Logger.Infow("query is nil ", header.GetRequestIDKV(ctx).Fuzzy()...)
	}
	s, ok := i.(map[string]interface{})
	if !ok {
		logger.Logger.Infow("query is not map ", header.GetRequestIDKV(ctx).Fuzzy()...)
	}
	if s != nil && len(s) != 0 {
		err := c.cond.ParseCondition(condition[_query])
		if err != nil {
			logger.Logger.WithName("form condition").Errorw(err.Error(), header.GetRequestIDKV(ctx).Fuzzy()...)
			return nil, err
		}
		dataes = append(dataes, condition[_query])
	}

	var newQuery permit.Object
	if len(dataes) != 0 {
		newQuery = permit.Object{
			_bool: types.M{
				_must: dataes,
			},
		}
	}
	switch req.Echo.Request().Method {
	case http.MethodGet:
		queryBytes, err := json.Marshal(newQuery)
		if err != nil {
			return nil, err
		}

		str, err := url.QueryUnescape(req.Echo.Request().URL.RawQuery)
		if err != nil {
			return nil, err
		}

		v, err := url.ParseQuery(str)
		if err != nil {
			return nil, err
		}

		v.Set("query", string(queryBytes))

		req.Echo.Request().URL.RawQuery = v.Encode()
	case http.MethodPost:
		req.Data[_query] = newQuery
	}
	req.Data["oldQuery"] = oldQuery
	return c.next.Do(ctx, req)
}
