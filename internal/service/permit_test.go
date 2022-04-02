package service

import (
	"fmt"
	"testing"
)

func TestName(t *testing.T) {
	cases := []string{
		"/api/v1/polyapi/request/system/app/vck4k/raw/inner/form/5tz5n/5tz5n_create.r",
		"/api/v1/polyapi/request/system/app/vck4k/raw/inner/form/form/5tz5n/5tz5n_create.r",
	}
	for i, v := range cases {
		fmt.Println(i+1, IsFormAPI(v), v)
	}

}
