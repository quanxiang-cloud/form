package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"

	"github.com/quanxiang-cloud/cabin/logger"
	"github.com/quanxiang-cloud/form/internal/auth/condition"
	"github.com/quanxiang-cloud/form/internal/auth/filters"
	"github.com/quanxiang-cloud/form/internal/service/consensus"
	"github.com/quanxiang-cloud/form/pkg/misc/config"
)

type formAuth struct {
	auth    *auth
	url     *url.URL
	tripper http.RoundTripper
	permit  *consensus.Permit
	cond    *condition.Condition
}

func NewFormAuth(conf *config.Config) (Auth, error) {
	auth, err := newAuth(conf)
	if err != nil {
		return nil, err
	}

	cond := condition.NewCondition()
	return &formAuth{
		auth: auth,
		tripper: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   conf.Transport.Timeout,
				KeepAlive: conf.Transport.KeepAlive,
				DualStack: true,
			}).DialContext,
			MaxIdleConns:          conf.Transport.MaxIdleConns,
			IdleConnTimeout:       conf.Transport.IdleConnTimeout,
			TLSHandshakeTimeout:   conf.Transport.TLSHandshakeTimeout,
			ExpectContinueTimeout: conf.Transport.ExpectContinueTimeout,
		},
		cond: cond,
	}, nil
}

func (f *formAuth) Forward(ctx context.Context, req *ReqParam) (*http.Response, error) {
	// 1、鉴权
	// havePermit, err := f.Auth(ctx, req)
	// if err != nil {
	// 	return nil, err
	// }

	// if !havePermit {
	// 	return nil, nil
	// }

	// 2、组装
	f.cond.Do(ctx, &condition.CondReq{
		UserID:   req.UserID,
		BodyData: req.Body,
	})
	// 3、转发

	// 4、过滤

	return nil, nil
}

func (f *formAuth) Auth(ctx context.Context, req *ReqParam) (bool, error) {
	f.cond.Do(ctx, &condition.CondReq{
		UserID:   req.UserID,
		BodyData: req.Body,
	})
	a, _ := json.Marshal(req.Body)
	fmt.Println(string(a))
	// resp, err := f.auth.Auth(ctx, req)
	// if err != nil {
	// 	return false, err
	// }

	// if resp == nil {
	// 	return false, nil
	// }

	// // access judgment
	// // if !filters.Pre(req.Entity, resp.Permit.Params) {
	// // 	return false, nil
	// // }

	// f.permit = resp.Permit
	return true, nil
}

func (f *formAuth) Filter(resp *http.Response, method string) error {
	respDate, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	conResp := &consensus.Response{}

	err = json.Unmarshal(respDate, conResp)
	if err != nil {
		return err
	}

	var entity interface{}
	switch method {
	case "get":
		entity = conResp.GetResp.Entity
	case "search":
		entity = conResp.ListResp.Entities
	}
	filters.Post(entity, f.permit.Response)

	data, err := json.Marshal(entity)
	if err != nil {
		logger.Logger.Errorf("entity json marshal failed: %s", err.Error())
		return err
	}

	resp.Body = io.NopCloser(bytes.NewReader(data))
	resp.ContentLength = int64(len(data))
	return nil
}
