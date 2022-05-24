package impl

import (
	"context"
	"fmt"

	id2 "github.com/quanxiang-cloud/cabin/id"
	time2 "github.com/quanxiang-cloud/cabin/time"
	"github.com/quanxiang-cloud/form/internal/models"
	"github.com/quanxiang-cloud/form/pkg/backup/aide"
)

var (
	exportTableRelationURL = "%s/api/v1/form/%s/internal/backup/export/tableRelation"
	importTableRelationURL = "%s/api/v1/form/%s/internal/backup/import/tableRelation"
)

// TableRelation tableRelation.
type TableRelation struct{}

func (tr *TableRelation) tag() string {
	return "tableRelations"
}

// Export export.
func (tr *TableRelation) Export(ctx context.Context, opts *aide.ExportOption) (map[string]aide.Object, error) {
	url := fmt.Sprintf(exportTableRelationURL, opts.Host, opts.AppID)

	obj, err := aide.ExportObject(ctx, url, opts)
	if err != nil {
		return nil, err
	}

	return map[string]aide.Object{
		tr.tag(): obj,
	}, nil
}

// Import import.
// nolint: dupl
func (tr *TableRelation) Import(ctx context.Context, objs map[string]aide.Object, opts *aide.ImportOption) (map[string]string, error) {
	obj := objs[tr.tag()]

	var tables []*models.TableRelation
	err := aide.Serialize(obj, &tables)
	if err != nil {
		return nil, err
	}

	ids := tr.replaceParam(tables, opts)

	data := make(aide.Object, len(obj))
	for i := 0; i < len(obj); i++ {
		data[i] = tables[i]
	}

	url := fmt.Sprintf(importTableRelationURL, opts.Host, opts.AppID)

	err = aide.ImportObject(ctx, url, data, opts)
	if err != nil {
		return nil, err
	}

	return ids, nil
}

func (tr *TableRelation) replaceParam(tableRelations []*models.TableRelation, opts *aide.ImportOption) map[string]string {
	ids := make(map[string]string)

	for i := 0; i < len(tableRelations); i++ {
		id := id2.HexUUID(true)
		ids[tableRelations[i].ID] = id

		tableRelations[i].ID = id
		tableRelations[i].AppID = opts.AppID
		tableRelations[i].CreatedAt = time2.NowUnix()
	}

	return ids
}
