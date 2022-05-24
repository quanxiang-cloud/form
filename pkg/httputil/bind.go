package httputil

import (
	"encoding/json"
	"net/http"

	"github.com/labstack/echo/v4"
)

// BindBody bind gin body
func BindBody(c echo.Context, d interface{}) error {
	db := new(echo.DefaultBinder)
	if err := db.BindBody(c, d); err != echo.ErrUnsupportedMediaType {
		return err
	}
	return nil
}

// GetRequestArgs get request args
func GetRequestArgs(c echo.Context, d interface{}) error {
	if d == nil {
		d = &json.RawMessage{}
	}
	method := c.Request().Method
	if method == http.MethodGet || method == http.MethodDelete || method == http.MethodHead {
		q := c.QueryParams()
		raw := QueryToBody(q, false)
		// FIXME return json.Unmarshal([]byte(raw), d)
		err := json.Unmarshal([]byte(raw), d)
		return err
	}
	// FIXME return BindBody(c, d)
	err := BindBody(c, d)
	return err
}

// IsQueryMethod check if http method is query
func IsQueryMethod(method string) bool {
	switch method {
	// NOTE: parameter is in query GET, DELETE, HEAD
	case http.MethodGet, http.MethodDelete, http.MethodHead:
		return true
	}
	return false
}
