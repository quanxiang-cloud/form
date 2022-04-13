package impl

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	id2 "github.com/quanxiang-cloud/cabin/id"
	time2 "github.com/quanxiang-cloud/cabin/time"
	"github.com/quanxiang-cloud/form/internal/models"
	"github.com/quanxiang-cloud/form/pkg/backup/internal/aide"
)

var (
	exportTableURL = "%s/api/v1/form/%s/internal/backup/export/table"
	importTableURL = "%s/api/v1/form/%s/internal/backup/import/table"
)

// Table table.
type Table struct{}

func (t *Table) tag() string {
	return "tables"
}

// Export export.
func (t *Table) Export(ctx context.Context, opts *aide.ExportOption) (map[string]aide.Object, error) {
	url := fmt.Sprintf(exportTableURL, opts.Host, opts.AppID)

	obj, err := aide.ExportObject(ctx, url, opts)
	if err != nil {
		return nil, err
	}

	return map[string]aide.Object{
		t.tag(): obj,
	}, nil
}

// Import import.
func (t *Table) Import(ctx context.Context, objs map[string]aide.Object, opts *aide.ImportOption) (map[string]string, error) {
	obj := objs[t.tag()]

	var tables []*models.Table
	err := aide.Serialize(obj, &tables)
	if err != nil {
		return nil, err
	}

	ids, err := t.replaceParam(tables, opts)
	if err != nil {
		return nil, err
	}

	data := make(aide.Object, len(obj))
	for i := 0; i < len(obj); i++ {
		data[i] = tables[i]
	}

	url := fmt.Sprintf(importTableURL, opts.Host, opts.AppID)

	err = aide.ImportObject(ctx, url, data, opts)
	if err != nil {
		return nil, err
	}

	return ids, nil
}

func (t *Table) replaceParam(tables []*models.Table, opts *aide.ImportOption) (map[string]string, error) {
	ids := make(map[string]string)
	if len(tables) == 0 {
		return nil, nil
	}

	oldAppID := tables[0].AppID
	bytes, err := json.Marshal(tables)
	if err != nil {
		return nil, err
	}

	str := strings.ReplaceAll(
		string(bytes),
		genTrimAppID(oldAppID),
		genTrimAppID(opts.AppID),
	)

	temp := make([]*models.Table, 0, len(tables))
	err = json.Unmarshal([]byte(str), &temp)
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(tables); i++ {
		id := id2.StringUUID()
		ids[tables[i].ID] = id

		tables[i].ID = id
		tables[i].AppID = opts.AppID
		tables[i].Schema = temp[i].Schema
		tables[i].CreatedAt = time2.NowUnix()
	}

	return ids, nil
}

func genTrimAppID(appID string) string {
	return fmt.Sprintf(`"appID":"%s"`, appID)
}
