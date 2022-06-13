package service

import (
	"fmt"
	"github.com/quanxiang-cloud/form/internal/models"
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

func TestD(t *testing.T) {
	s := make(models.FiledPermit)
	d := make(models.FiledPermit)
	s["dke"] = models.Key{
		Type: "string",
	}
	s["object"] = models.Key{
		Type: "object",
		Properties: map[string]models.Key{
			"sss": {
				Type: "string",
			},
			"sss1": {
				Type: "string",
			},
		},
	}
	//d["object"] = models.Key{
	//	Type: "object",
	//	Properties: map[string]models.Key{
	//		"sss2": {
	//			Type: "string",
	//		},
	//	},
	//}
	FiledPermitPoly(s, d)
	fmt.Println(s, d)

}
