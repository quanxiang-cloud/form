package httputil

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	error2 "github.com/quanxiang-cloud/cabin/error"
	"github.com/quanxiang-cloud/cabin/logger"
)

var log = logger.Logger.WithName("httputil")

const (
	timeout = 10
)

var defaultClient = &http.Client{
	Transport: &http.Transport{
		Dial: func(netw, addr string) (net.Conn, error) {
			deadline := time.Now().Add(timeout * time.Second)
			c, err := net.DialTimeout(netw, addr, time.Second*timeout)
			if err != nil {
				return nil, err
			}
			c.SetDeadline(deadline)
			return c, nil
		},
		MaxIdleConns:      0,
		DisableKeepAlives: true,
	},
}

// MakeRefer get refer from request
func MakeRefer(req *http.Request) string {
	return fmt.Sprintf("%s://%s%s", "http", req.Host, req.URL.EscapedPath())
}

// HTTPRequest implements general http request methods, and return the response string body
// BUG: return (string, http.Header, error) run fail
func HTTPRequest(reqURL, method, data string, header http.Header, owner string) (string, *http.Response, error) {
	// BUG: when client=http.DefaultClient
	/*
		httpRequest url="https://api.66mz8.com:443/api/translation.php?info=apple" body={
		    "info": "apple",
		  } header=map[Accept-Encoding:[gzip] Connection:[upgrade]
		        Content-Length:[84] Content-Type:[application/json] Request-Id:[req_stWlg0JTlCZw] User-Agent:[Go-http-client/1.1] X-Forwarded-For:[192.168.200.2, 127.0.0.1, 127.0.0.1] X-Real-Ip:[127.0.0.1]]
		http.Do	{"error": "Get \"https://api.66mz8.com:443/api/translation.php?info=apple":
		        http2: invalid Connection request header: [\"upgrade\"]"}
	*/
	body, resp, err := httpRequest(reqURL, method, data, header, owner, defaultClient)
	if err != nil {
		err = error2.NewErrorWithString(error2.Internal, err.Error())
	}
	return body, resp, err
}

func httpRequest(reqURL, method, data string, header http.Header, owner string, client *http.Client) (string, *http.Response, error) {
	uri, err := url.Parse(reqURL)
	if err != nil {
		return "", nil, err
	}

	var reader io.Reader
	// build URI parameter
	if method == http.MethodGet || method == http.MethodDelete {
		reader = nil
		uri.RawQuery = BodyToQuery(data)
	} else {
		reader = strings.NewReader(data)
	}

	putErrLog := func(err error, stage string) {
		log.PutError(err, stage, "url", uri.String(), "header", header, "body", data)
	}

	logger.Logger.Debugf("httpRequest url=%q body=%s header=%+v ", uri.String(), data, header)

	req, err := http.NewRequest(method, uri.String(), reader)
	if err != nil {
		putErrLog(err, "http.NewRequest")
		return "", nil, err
	}
	req.Header = header

	resp, err := client.Do(req) // http.Get(url)
	if err != nil {
		putErrLog(err, "http.Do")
		return "{}", nil, err
	}

	var respBody []byte
	if resp.StatusCode != http.StatusOK {
		putErrLog(nil, "http.Resp:"+resp.Status)
		resp.Status = fmt.Sprintf("proxy:%s", resp.Status)
		respBody = []byte(resp.Status)
	} else {
		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			putErrLog(err, "http.ReadResp")
			return "{}", nil, err
		}
		respBody = body
	}

	logger.Logger.Debugf("httpResponse body=%s header=%+v status=\"%s\" statusCode=%d", string(respBody), resp.Header, resp.Status, resp.StatusCode)

	return string(respBody), resp, nil
}

// AllowMethod check if http method of request is allow
func AllowMethod(expect string, got string) bool {
	//NOTE: allow POST for every request
	if got == http.MethodPost {
		return true
	}

	return expect == got
}
