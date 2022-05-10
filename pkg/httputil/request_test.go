package httputil

import (
	"fmt"
	"net/http"
	"testing"
)

func TestRequest(t *testing.T) {
	h := http.Header{}
	h.Set("Accept-Encoding", "gzip")
	h.Set("Connection", "upgrade")
	h.Set("Content-Type", "application/json")
	body, resp, err := httpRequest(
		"https://api.66mz8.com:443/api/translation.php",
		http.MethodGet,
		`{"info":"apple"}`,
		h,
		"",
		http.DefaultClient,
	)
	fmt.Println(body, resp, err)
}

func TestRequest2(t *testing.T) {
	h := http.Header{}
	//h.Set("Content-Type", "application/json")
	body, resp, err := httpRequest(
		"https://home.yunify.com:443/distributor.action",
		http.MethodGet,
		`{"serviceName":"clogin"}`,
		h,
		"",
		http.DefaultClient,
	)
	fmt.Println(body, resp, err)
}

func TestRequest3(t *testing.T) {
	h := http.Header{}
	//h.Set("Content-Type", "application/json")
	body, resp, err := httpRequest(
		"https://api.qingcloud.com:443/iaas/",
		http.MethodGet,
		`{"access_key_id": "PDOIUFDSIK",
"action": "DescribeInstances",
"limit": 15,
"offset": 0,
"signature": "iiIPYusdkwlkejkdksahdfk78=",
"signature_method": "HmacSHA256",
"signature_version": "1",
"status": ["running","stopped"],
"time_stamp": "2022-03-26T00:55:07Z",
"version": "1",
"zone": "ap2a"
}`,
		h,
		"",
		http.DefaultClient,
	)
	fmt.Println(body, resp, err)
}
