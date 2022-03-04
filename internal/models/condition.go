package models

import "github.com/quanxiang-cloud/form/internal/service/types"

type Condition struct {
	Bool BOOL `json:"bool,omitempty"`
}

type BOOL struct {
	Must   types.Entities `json:"must,omitempty"`
	Should types.Entities `json:"should,omitempty"`
}
