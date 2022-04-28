package lowcode

import (
	"context"
	"encoding/json"
	"fmt"
	e "github.com/quanxiang-cloud/cabin/error"
	"github.com/quanxiang-cloud/cabin/logger"
	"github.com/quanxiang-cloud/cabin/tailormade/client"
	"github.com/quanxiang-cloud/cabin/tailormade/header"
	"github.com/quanxiang-cloud/cabin/tailormade/resp"
	"github.com/quanxiang-cloud/form/pkg/misc/config"
	"io/ioutil"
	"net/http"
	"net/url"
)

type searchAPI struct {
	client http.Client
	conf   *config.Config
}

func (s *searchAPI) Subordinate(ctx context.Context, userID string) (*SubordinateResp, error) {
	resp := &SubordinateResp{}
	values := url.Values{}
	values.Set("query", "{query(page:0,size:100){users{id},total}}")
	url := fmt.Sprintf("%s%s", s.conf.Endpoint.Search, "/api/v1/search/subordinate")
	err := Get(ctx, &s.client, url, values, userID, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// SearchAPI SearchAPI
type SearchAPI interface {
	Subordinate(ctx context.Context, userID string) (*SubordinateResp, error)
}

type SubordinateResp struct {
	Users []*users `json:"users"`
	Total int64    `json:"total"`
}

type users struct {
	ID string `json:"id"`
}

func NewSearchAPI(conf *config.Config) SearchAPI {
	return &searchAPI{
		client: client.New(conf.InternalNet),
		conf:   conf,
	}
}

//Get Get
func Get(ctx context.Context, client *http.Client, uri string, values url.Values, userID string, entity interface{}) error {
	u, err := url.ParseRequestURI(uri)
	if err != nil {
		return err
	}
	if values != nil {
		u.RawQuery = values.Encode() // URL encode
	}
	logger.Logger.Info("url is", u.String())
	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	// header 封装request user-id user-name department-id role
	req.Header.Add(header.GetRequestIDKV(ctx).Wreck())
	req.Header.Add(header.GetTimezone(ctx).Wreck())
	req.Header.Add("User-Id", userID)
	response, err := client.Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("expected state value is 200, actually %d", response.StatusCode)
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	return decomposeBody(body, entity)
}

func decomposeBody(body []byte, entity interface{}) error {
	r := new(resp.Resp)
	r.Data = entity

	err := json.Unmarshal(body, r)
	if err != nil {
		return err
	}

	if r.Code != e.Success {
		return r.Error
	}

	return nil
}
