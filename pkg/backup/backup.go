package backup

import (
	"context"
	"net/http"
	"os"

	"github.com/quanxiang-cloud/cabin/tailormade/client"
	"github.com/quanxiang-cloud/form/internal/models"
	"github.com/quanxiang-cloud/form/pkg/backup/internal/aide"
	"github.com/quanxiang-cloud/form/pkg/backup/internal/aide/impl"
)

var formHost string

// nolint:gochecknoinits
func init() {
	formHost = os.Getenv("FORM_HOST")
	if formHost == "" {
		formHost = "http://127.0.0.1:8080"
	}
}

// Backup backup.
type Backup struct {
	formHost string
	client   http.Client
}

// NewBackup create a backup instance.
func NewBackup(conf client.Config) *Backup {
	return &Backup{
		client:   client.New(conf),
		formHost: formHost,
	}
}

var aides = []aide.Aide{
	&impl.Table{},
	&impl.TableRelation{},
	&impl.TableSchema{},
	&impl.Role{},
}

// Result is the result of export.
type Result struct {
	Permits        []*models.Permit        `json:"permits"`
	TableSchemas   []*models.TableSchema   `json:"tableSchemas"`
	Roles          []*models.Role          `json:"roles"`
	Tables         []*models.Table         `json:"tables"`
	TableRelations []*models.TableRelation `json:"tableRelations"`
}

// Export export.
func (b *Backup) Export(ctx context.Context, opts *aide.ExportOption) (*Result, error) {
	result := &Result{}
	dataes := make(map[string]interface{})

	opts.Client = b.client
	opts.Host = b.formHost

	for _, a := range aides {
		objs, err := a.Export(ctx, opts)
		if err != nil {
			return nil, err
		}

		for key, val := range objs {
			dataes[key] = val
		}
	}

	err := aide.Serialize(dataes, result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// Import import.
func (b *Backup) Import(ctx context.Context, result *Result, opts *aide.ImportOption) (map[string]string, error) {
	ids := make(map[string]string)

	var objs map[string]aide.Object
	err := aide.Serialize(result, &objs)
	if err != nil {
		return nil, err
	}

	opts.Client = b.client
	opts.Host = b.formHost

	for _, a := range aides {
		idMap, err := a.Import(ctx, objs, opts)
		if err != nil {
			return nil, err
		}

		for key, val := range idMap {
			ids[key] = val
		}
	}

	return ids, nil
}
