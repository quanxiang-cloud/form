package api

import "github.com/quanxiang-cloud/form/internal/service"

type CustomPage struct {
	customPage service.CustomPage
	//menu       service.Menu
	permission service.Permission
}
