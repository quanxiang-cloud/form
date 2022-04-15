package models

import "context"

// SerialScheme SerialScheme
type SerialScheme struct {
	Bit   string `json:"bit"`
	Value string `json:"value"`
	Step  string `json:"step"`
	Date  string `json:"date"`
}

// SerialExcute SerialExcute
type SerialExcute struct {
	Date string
	Incr string
}

// SerialRepo SerialRepo
type SerialRepo interface {
	Create(ctx context.Context, appID, tableID, fieldID string, values map[string]interface{}) error
	Get(ctx context.Context, appID, tableID, fieldID, field string) string
	GetAll(ctx context.Context, appID, tableID, fieldID string) map[string]string
}
