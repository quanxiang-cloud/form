package backup

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"

	"github.com/quanxiang-cloud/cabin/tailormade/client"
)

var formHost string

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

func init() {
	formHost = os.Getenv("FORM_HOST")
	if formHost == "" {
		formHost = "http://form:8080"
	}
}

type Backup struct {
	client http.Client
}

// NewBackup create a backup instance.
func NewBackup(conf client.Config) *Backup {
	return &Backup{
		client: client.New(conf),
	}
}

// Object is the backup object.
type Object []interface{}

// BackupReq is the request of export.
type BackupReq struct {
	AppID string `json:"appID"`
	Page  int    `json:"page"`
	Size  int    `json:"size"`
}

// BackupResp is the response of export.
type BackupResp struct {
	Result  `json:",inline"`
	Data    Object `json:"data"`
	Count   int    `json:"count"`
	HasNext bool   `json:"hasNext"`
}

// Result is the result of export.
type Result struct {
	Permits        Object `json:"permits"`
	Schemas        Object `json:"schemas"`
	Roles          Object `json:"roles"`
	Tables         Object `json:"tables"`
	TableRelations Object `json:"tableRelations"`
}

const (
	startPage = 1
	maxSize   = 999
)

func (b *Backup) Export(ctx context.Context, appID string, w io.Writer) error {
	result := &Result{}

	for _, backup := range backups {
		backup.Initialize(appID, b.client)

		err := backup.Export(ctx, result)
		if err != nil {
			return err
		}
	}

	dataBytes, err := json.Marshal(result)
	if err != nil {
		return err
	}

	_, err = io.Copy(w, bytes.NewReader(dataBytes))

	return err
}

type ImportReq struct {
	Data Object `json:"data"`
}

type ImportResp struct{}

func (b *Backup) Import(ctx context.Context, appID string, r io.Reader) error {
	buf := bytes.Buffer{}
	_, err := io.Copy(&buf, r)
	if err != nil {
		return err
	}

	result := &Result{}
	err = json.Unmarshal(buf.Bytes(), result)
	if err != nil {
		return err
	}

	for _, backup := range backups {
		backup.Initialize(appID, b.client)

		err := backup.Import(ctx, result)
		if err != nil {
			return err
		}
	}

	return nil
}

type backup interface {
	Initialize(string, http.Client)
	Export(context.Context, *Result) error
	Import(context.Context, *Result) error
}

var backups = []backup{
	&table{},
	&tableRelation{},
	&tableScheme{},
	&permit{},
	&role{},
}
