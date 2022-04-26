package httputil

import (
	"encoding/json"
	"fmt"
	"net/url"
	"testing"
)

func TestBodyToQuery(t *testing.T) {
	var v = map[string]interface{}{
		"a": "foo",
		"b": []string{"foo", "bar"},
		"c": 123,
		"d": 123.456,
		"e": true,
		"f": map[string]interface{}{
			"x": "xx",
			"y": []string{"foo2", "bar2"},
			"z": 123.4,
		},
	}
	b, _ := json.Marshal(v)
	s := string(b)
	fmt.Println(s)
	q := BodyToQuery(s)
	fmt.Println(q)
	qq, err := url.ParseQuery(q)
	fmt.Println(qq, err)
	fmt.Println(QueryToBody(qq, true))
}
