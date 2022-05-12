package side

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/quanxiang-cloud/cabin/lib/httputil"
	"github.com/quanxiang-cloud/cabin/logger"
	"github.com/quanxiang-cloud/form/internal/models"
	"github.com/quanxiang-cloud/form/internal/permit/treasure"
	"github.com/quanxiang-cloud/form/internal/service/consensus"
	httputil2 "github.com/quanxiang-cloud/form/pkg/httputil"

	"github.com/quanxiang-cloud/form/internal/permit"
	"github.com/quanxiang-cloud/form/pkg/misc/config"
)

type Proxy struct {
	url *url.URL

	transport http.RoundTripper

	next permit.Permit

	isPermit bool
}

func NewProxy(conf *config.Config, rawurl string) (*Proxy, error) {
	url, err := url.Parse(rawurl)
	if err != nil {
		return nil, err
	}
	return &Proxy{
		url:       url,
		transport: httputil2.Transport(conf),
		isPermit:  true,
	}, nil
}

func NewNilModifyProxy(conf *config.Config, rawurl string) (*Proxy, error) {
	url, err := url.Parse(rawurl)
	if err != nil {
		return nil, err
	}
	return &Proxy{
		url:       url,
		transport: httputil2.Transport(conf),
		isPermit:  false,
	}, nil
}

func (p *Proxy) Do(ctx context.Context, req *permit.Request) (*permit.Response, error) {
	var filters httputil2.ModifyResponse
	if p.isPermit {
		filters = Filter(req.Permit)
	}
	// FIXME typo
	err := httputil2.DoPoxy(ctx, req, &httputil2.Proxys{
		Url:       p.url,
		Transport: p.transport,
	}, filters)
	if err != nil {
		return nil, err
	}
	return &permit.Response{}, nil
}

const (
	mimeApplicationJSON = "application/json"
)

func Filter(permit *consensus.Permit) httputil2.ModifyResponse {
	return func(resp *http.Response) error {
		return filter(resp, permit)
	}
}

func filter(resp *http.Response, permit *consensus.Permit) (err error) {
	if resp.StatusCode < http.StatusOK || resp.StatusCode > http.StatusIMUsed {
		return nil
	}

	response := httputil.NewResponse(resp)

	logger.Logger.Info("content-type ", response.ContentType())
	if strings.HasPrefix(response.ContentType(), mimeApplicationJSON) {
		return doFilterJSON(response, permit)
	}

	// FIXME we do not care other Content-Type.
	_, err = response.ReadRawBody(http.DefaultMaxHeaderBytes)
	if err != nil {
		return err
	}

	buf := strings.NewReader("")
	resp.Body = io.NopCloser(buf)
	resp.ContentLength = buf.Size()
	resp.Header.Set("Content-Length", fmt.Sprint(buf.Size()))

	return nil
}

func doFilterJSON(resp *httputil.Response, permit *consensus.Permit) (err error) {
	if permit == nil {
		return nil
	}

	if permit.Types == models.InitType || permit.ResponseAll {
		return nil
	}
	// FIXME type
	respDate, err := resp.DecodeCloseBody(http.DefaultMaxHeaderBytes)
	if err != nil {
		return err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(respDate, &result); err != nil {
		return err
	}

	treasure.Filter(result, permit.Response)

	data, err := json.Marshal(result)
	if err != nil {
		logger.Logger.Errorf("entity json marshal failed: %s", err.Error())
		return err
	}

	err = resp.EncodeWriteBody(data, false)
	if err != nil {
		return err
	}

	resp.ContentLength = int64(len(data))
	resp.Header.Set("Content-Length", strconv.Itoa(len(data)))

	return nil
}
