package httputil

import (
	"fmt"
	"testing"
)

func TestQueryToBody(t *testing.T) {
	var v = map[string][]string{}
	fmt.Println(v)
	raw := QueryToBody(v, true)
	fmt.Println(raw)
}
