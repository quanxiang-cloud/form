package side

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

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
	contentType         = "Content-Type"
	contentEncoding     = "Content-Encoding"
	mimeApplicationJSON = "application/json"
	mimeGzip            = "gzip"
)

func Filter(permit *consensus.Permit) httputil2.ModifyResponse {
	return func(resp *http.Response) error {
		return filter(resp, permit)
	}
}

func filter(resp *http.Response, permit *consensus.Permit) (err error) {
	if resp.StatusCode != http.StatusOK {
		return nil
	}

	var (
		cTypeFlag     = false
		cEncodingFlag = false
	)

	ctype := resp.Header.Get(contentType)
	if strings.HasPrefix(strings.ToLower(ctype), mimeApplicationJSON) {
		cTypeFlag = true
	}

	cEncoding := resp.Header.Get(contentEncoding)
	if strings.Contains(strings.ToLower(cEncoding), mimeGzip) {
		cEncodingFlag = true
	}

	if cTypeFlag {
		return doFilterJSON(resp, permit, cEncodingFlag)
	}

	defer resp.Body.Close()
	_, err = io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	buf := []byte("")
	resp.Body = io.NopCloser(bytes.NewReader(buf))
	resp.ContentLength = int64(len(buf))
	resp.Header.Set("Content-Length", fmt.Sprint(len(buf)))

	return nil
}

func doFilterJSON(resp *http.Response, permit *consensus.Permit, cEncodingFlag bool) (err error) {
	if permit == nil {
		return nil
	}

	if permit.Types == models.InitType || permit.ResponseAll {
		return nil
	}

	var respDate []byte
	if cEncodingFlag {
		reder, err := gzip.NewReader(resp.Body)
		if err != nil {
			return err
		}

		defer reder.Close()
		respDate, err = io.ReadAll(reder)
	} else {
		respDate, err = io.ReadAll(resp.Body)
	}
	defer resp.Body.Close()
	if err != nil {
		return err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(respDate, &result); err != nil {
		return err
	}

	//if result["code"] != error2.Success {
	//	return nil
	//}

	if !permit.ResponseAll {
		treasure.Filter(result, permit.Response)
	}
	data, err := json.Marshal(result)
	if err != nil {
		logger.Logger.Errorf("entity json marshal failed: %s", err.Error())
		return err
	}

	buf := bytes.Buffer{}
	w := gzip.NewWriter(&buf)

	if cEncodingFlag {
		_, err = w.Write(data)
		w.Close()
	} else {
		_, err = buf.Write(data)
	}

	if err != nil {
		return err
	}

	resp.Body = io.NopCloser(bytes.NewReader(buf.Bytes()))
	resp.ContentLength = int64(len(data))
	resp.Header.Set("Content-Length", fmt.Sprint(len(data)))

	return nil
}
