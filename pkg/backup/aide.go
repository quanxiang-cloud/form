package backup

import (
	"context"
	"net/http"

	"github.com/quanxiang-cloud/form/internal/models"
)

type backup interface {
	Export(context.Context, *Result, *ExportOption) error
	Import(context.Context, *Result, *ImportOption) error
}

var backups = []backup{
	&table{},
	&tableRelation{},
	&tableScheme{},
	&permit{},
	&role{},
}

var (
	exportTableURL         = "%s/api/v1/form/%s/internal/backup/export/table"
	exportPermitURL        = "%s/api/v1/form/%s/internal/backup/export/permit"
	exportRoleURL          = "%s/api/v1/form/%s/internal/backup/export/role"
	exportTableRelationURL = "%s/api/v1/form/%s/internal/backup/export/tableRelation"
	exportTableSchemaURL   = "%s/api/v1/form/%s/internal/backup/export/tableSchema"
)

var (
	importTableURL         = "%s/api/v1/form/%s/internal/backup/import/table"
	importPermitURL        = "%s/api/v1/form/%s/internal/backup/import/permit"
	importRoleURL          = "%s/api/v1/form/%s/internal/backup/import/role"
	importTableRelationURL = "%s/api/v1/form/%s/internal/backup/import/tableRelation"
	importTableSchemaURL   = "%s/api/v1/form/%s/internal/backup/import/tableSchema"
)

const (
	startPage = 1
	maxSize   = 999
)

// Result is the result of export.
type Result struct {
	Permits        []*models.Permit        `json:"permits"`
	TableSchemas   []*models.TableSchema   `json:"tableSchemas"`
	Roles          []*models.Role          `json:"roles"`
	Tables         []*models.Table         `json:"tables"`
	TableRelations []*models.TableRelation `json:"tableRelations"`
}

// ExportReq is the request of export.
type ExportReq struct {
	AppID string `json:"appID"`
	Page  int    `json:"page"`
	Size  int    `json:"size"`
}

// ExportResp is the response of export.
type ExportResp struct {
	Result `json:",inline"`
	Count  int `json:"count"`
}

// ExportOption is the option of export.
type ExportOption struct {
	AppID string `json:"appID"`
	Page  int    `json:"page"`
	Size  int    `json:"size"`

	Client http.Client
}

// ImportReq import request.
type ImportReq struct {
	Result `json:",inline"`
}

// ImportResp import response.
type ImportResp struct{}

// ImportOption is the option of import.
type ImportOption struct {
	AppID string `json:"appID"`

	Client http.Client
}
