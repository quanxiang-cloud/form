package backup

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	id2 "github.com/quanxiang-cloud/cabin/id"
	"github.com/quanxiang-cloud/cabin/logger"
	"github.com/quanxiang-cloud/cabin/tailormade/client"
	"github.com/quanxiang-cloud/form/internal/models"
)

type table struct{}

func (t *table) Export(ctx context.Context, result *Result, opts *ExportOption) error {
	url := fmt.Sprintf(exportTableURL, formHost, opts.AppID)

	return exportTableFunc(ctx, url,
		ExportReq{
			AppID: opts.AppID,
			Page:  startPage,
			Size:  maxSize,
		},
		result, opts)
}

func (t *table) Import(ctx context.Context, result *Result, opts *ImportOption) error {
	var (
		index int
		req   ImportReq
		url   = fmt.Sprintf(importTableURL, formHost, opts.AppID)
	)
	tables, err := importTableFunc(result.Tables, opts.AppID)
	if err != nil {
		return err
	}

	if len(tables)%maxSize == 0 {
		index = len(tables) / maxSize
	} else {
		index = len(tables)/maxSize + 1
	}

	for i := 0; i < index; i++ {
		if i == index-1 {
			req.Tables = tables[i*maxSize:]
		} else {
			req.Tables = tables[i*maxSize : (i+1)*maxSize]
		}

		err := client.POST(ctx, &opts.Client, url, req, &ImportResp{})
		if err != nil {
			logger.Logger.WithName("export request").Errorf("send http request failed: %v", err)

			return err
		}
	}

	return nil
}

func importTableFunc(tables []*models.Table, appID string) ([]*models.Table, error) {
	if len(tables) == 0 {
		return nil, nil
	}
	oldAppID := tables[0].AppID

	b, err := json.Marshal(tables)
	if err != nil {
		return nil, err
	}

	t := make([]*models.Table, 0, len(tables))

	str := strings.ReplaceAll(string(b), fmt.Sprintf(`"appID":"%s"`, oldAppID), fmt.Sprintf(`"appID":"%s"`, appID))

	err = json.Unmarshal([]byte(str), &t)

	for i := 0; i < len(t); i++ {
		t[i].AppID = appID
		t[i].ID = id2.StringUUID()
	}

	return t, nil
}

func exportTableFunc(ctx context.Context, url string, req ExportReq, result *Result, opts *ExportOption) error {
	totalPage := 0

	for {
		resp := &ExportResp{}

		err := client.POST(ctx, &opts.Client, url, req, resp)
		if err != nil {
			logger.Logger.WithName("export request").Errorf("send http request failed: %v", err)

			return err
		}

		result.Tables = append(result.Tables, resp.Tables...)

		if resp.Count <= req.Size {
			break
		}

		if totalPage == 0 {
			if (resp.Count % req.Size) == 0 {
				if resp.Count/req.Size < 1 {
					totalPage = 1
				} else {
					totalPage = resp.Count / req.Size
				}
			} else {
				totalPage = resp.Count/req.Size + 1
			}
		}

		if totalPage <= req.Page {
			break
		}

		req.Page++
	}

	return nil
}
