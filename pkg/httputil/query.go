package httputil

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
)

// BodyToQuery convert a JSON body to http GET query parameter
func BodyToQuery(data string) string {
	var d interface{}
	if err := json.Unmarshal([]byte(data), &d); err != nil {
		return ""
	}
	buf := bytes.NewBuffer(nil)
	if err := genQuery("", d, buf, 0); err != nil {
		return ""
	}
	s := buf.String()
	if len(s) > 0 {
		s = s[1:] // NOTE: remove prefix &
	}
	return s
}

func genQuery(name string, d interface{}, buf *bytes.Buffer, depth int) error {
	if depth >= 20 {
		return errors.New("buildQuery out of recursion")
	}
	writeSingle := func(v interface{}) {
		buf.WriteString(fmt.Sprintf("&%s=", name))
		buf.WriteString(url.QueryEscape(fmt.Sprint(v)))
	}
	switch v := d.(type) {
	case string:
		writeSingle(v)
	case float64:
		writeSingle(v)
	case bool:
		writeSingle(v)
	case map[string]interface{}:
		for k, vv := range v {
			n := k
			if name != "" {
				n = fmt.Sprintf("%s.%s", name, k)
			}
			genQuery(n, vv, buf, depth+1)
		}
	case []interface{}:
		if name != "" {
			for i, vv := range v {
				n := fmt.Sprintf("%s.%d", name, i+1)
				genQuery(n, vv, buf, depth+1)
			}
		}
	}
	return nil
}

// ObjectBodyToQuery convert a JSON body to http GET query parameter
func ObjectBodyToQuery(d interface{}) string {
	buf := bytes.NewBuffer(nil)
	if err := genQuery("", d, buf, 0); err != nil {
		return ""
	}
	s := buf.String()
	if len(s) > 0 {
		s = s[1:] // NOTE: remove prefix &
	}
	return s
}
