package impl

import (
	"context"
	"fmt"

	id2 "github.com/quanxiang-cloud/cabin/id"
	time2 "github.com/quanxiang-cloud/cabin/time"
	"github.com/quanxiang-cloud/form/internal/models"
	"github.com/quanxiang-cloud/form/pkg/backup/internal/aide"
)

var (
	exportTableSchemaURL = "%s/api/v1/form/%s/internal/backup/export/tableSchema"
	importTableSchemaURL = "%s/api/v1/form/%s/internal/backup/import/tableSchema"
)

// TableSchema tableSchema.
type TableSchema struct{}

func (ts *TableSchema) tag() string {
	return "tableSchemas"
}

// Export export.
func (ts *TableSchema) Export(ctx context.Context, opts *aide.ExportOption) (map[string]aide.Object, error) {
	url := fmt.Sprintf(exportTableSchemaURL, opts.Host, opts.AppID)

	obj, err := aide.ExportObject(ctx, url, opts)
	if err != nil {
		return nil, err
	}

	return map[string]aide.Object{
		ts.tag(): obj,
	}, nil
}

// Import import.
// nolint: dupl
func (ts *TableSchema) Import(ctx context.Context, objs map[string]aide.Object, opts *aide.ImportOption) (map[string]string, error) {
	obj := objs[ts.tag()]

	var tables []*models.TableSchema
	err := aide.Serialize(obj, &tables)
	if err != nil {
		return nil, err
	}

	ids := ts.replaceParam(tables, opts)

	data := make(aide.Object, len(obj))
	for i := 0; i < len(obj); i++ {
		data[i] = tables[i]
	}

	url := fmt.Sprintf(importTableSchemaURL, opts.Host, opts.AppID)

	err = aide.ImportObject(ctx, url, data, opts)
	if err != nil {
		return nil, err
	}

	return ids, nil
}

func (ts *TableSchema) replaceParam(tableSchemas []*models.TableSchema, opts *aide.ImportOption) map[string]string {
	ids := make(map[string]string)

	for i := 0; i < len(tableSchemas); i++ {
		id := id2.StringUUID()
		ids[tableSchemas[i].ID] = id

		tableSchemas[i].ID = id
		tableSchemas[i].AppID = opts.AppID
		tableSchemas[i].CreatorID = opts.UserID
		tableSchemas[i].CreatorName = opts.UserName
		tableSchemas[i].EditorID = opts.UserID
		tableSchemas[i].EditorName = opts.UserName
		tableSchemas[i].CreatedAt = time2.NowUnix()
		tableSchemas[i].UpdatedAt = time2.NowUnix()
	}

	return ids
}
