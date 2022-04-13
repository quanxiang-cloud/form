package backup

import (
	"context"
	"net/http"
	"os"

	"github.com/quanxiang-cloud/cabin/tailormade/client"
)

var formHost string

func init() {
	formHost = os.Getenv("FORM_HOST")
	if formHost == "" {
		formHost = "http://form:8080"
	}
}

// Backup backup.
type Backup struct {
	client http.Client
}

// NewBackup create a backup instance.
func NewBackup(conf client.Config) *Backup {
	return &Backup{
		client: client.New(conf),
	}
}

// Export export.
func (b *Backup) Export(ctx context.Context, appID string) (*Result, error) {
	result := &Result{}

	for _, backup := range backups {
		err := backup.Export(ctx, result, &ExportOption{
			AppID:  appID,
			Page:   startPage,
			Size:   maxSize,
			Client: b.client,
		})
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}

// Import import.
func (b *Backup) Import(ctx context.Context, result *Result, appID string) error {
	for _, backup := range backups {
		err := backup.Import(ctx, result, &ImportOption{
			AppID:  appID,
			Client: b.client,
		})
		if err != nil {
			return err
		}
	}

	return nil
}
