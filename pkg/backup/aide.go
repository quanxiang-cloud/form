package backup

import (
	"context"
	"fmt"
	"net/http"

	"github.com/quanxiang-cloud/cabin/logger"
	"github.com/quanxiang-cloud/cabin/tailormade/client"
)

type table struct {
	appID     string
	exportURL string
	importURL string
	client    http.Client
}

func (t *table) Initialize(appID string, client http.Client) {
	t.client = client
	t.appID = appID
	t.exportURL = fmt.Sprintf(exportTableURL, formHost, appID)
	t.importURL = fmt.Sprintf(importTableURL, formHost, appID)
}

func (t *table) Export(ctx context.Context, result *Result) error {
	data, err := export(ctx, &t.client, t.exportURL, BackupReq{
		AppID: t.appID,
		Page:  startPage,
		Size:  maxSize,
	})
	if err != nil {
		logger.Logger.WithName("export table").Errorf("send http request failed: %v", err)

		return err
	}

	result.Tables = data

	return nil
}

func (t *table) Import(ctx context.Context, result *Result) error {
	err := import2(ctx, &t.client, t.importURL, ImportReq{
		Data: result.Tables,
	})
	if err != nil {
		logger.Logger.WithName("import table").Errorf("import table failed: %v", err)

		return err
	}

	return nil
}

type tableRelation struct {
	appID     string
	exportURL string
	importURL string
	client    http.Client
}

func (t *tableRelation) Initialize(appID string, client http.Client) {
	t.client = client
	t.appID = appID
	t.exportURL = fmt.Sprintf(exportTableRelationURL, formHost, appID)
	t.importURL = fmt.Sprintf(importTableRelationURL, formHost, appID)
}

func (t *tableRelation) Export(ctx context.Context, result *Result) error {
	data, err := export(ctx, &t.client, t.exportURL, BackupReq{
		AppID: t.appID,
		Page:  startPage,
		Size:  maxSize,
	})
	if err != nil {
		logger.Logger.WithName("export tableRelation").Errorf("export tableRelation failed: %v", err)
		return err
	}

	result.TableRelations = data
	return nil
}

func (t *tableRelation) Import(ctx context.Context, result *Result) error {
	err := import2(ctx, &t.client, t.importURL, ImportReq{
		Data: result.TableRelations,
	})
	if err != nil {
		logger.Logger.WithName("import tableRelation").Errorf("import tableRelation failed: %v", err)

		return err
	}

	return nil
}

type tableScheme struct {
	appID     string
	exportURL string
	importURL string
	client    http.Client
}

func (t *tableScheme) Initialize(appID string, client http.Client) {
	t.client = client
	t.appID = appID
	t.exportURL = fmt.Sprintf(exportTableSchemaURL, formHost, appID)
	t.importURL = fmt.Sprintf(importTableSchemaURL, formHost, appID)
}

func (t *tableScheme) Export(ctx context.Context, result *Result) error {
	data, err := export(ctx, &t.client, t.exportURL, BackupReq{
		AppID: t.appID,
		Page:  startPage,
		Size:  maxSize,
	})
	if err != nil {
		logger.Logger.WithName("export tableScheme").Errorf("export tableScheme failed: %v", err)
		return err
	}

	result.Schemas = data

	return nil
}

func (t *tableScheme) Import(ctx context.Context, result *Result) error {
	err := import2(ctx, &t.client, t.importURL, ImportReq{
		Data: result.Schemas,
	})
	if err != nil {
		logger.Logger.WithName("import tableScheme").Errorf("import tableScheme failed: %v", err)

		return err
	}

	return nil
}

type permit struct {
	appID     string
	exportURL string
	importURL string
	client    http.Client
}

func (t *permit) Initialize(appID string, client http.Client) {
	t.client = client
	t.appID = appID
	t.exportURL = fmt.Sprintf(exportPermitURL, formHost, appID)
	t.importURL = fmt.Sprintf(importPermitURL, formHost, appID)
}

func (t *permit) Export(ctx context.Context, result *Result) error {
	data, err := export(ctx, &t.client, t.exportURL, BackupReq{
		AppID: t.appID,
		Page:  startPage,
		Size:  maxSize,
	})
	if err != nil {
		logger.Logger.WithName("export permit").Errorf("export permit failed: %v", err)

		return err
	}

	result.Permits = data

	return nil
}

func (t *permit) Import(ctx context.Context, result *Result) error {
	err := import2(ctx, &t.client, t.importURL, ImportReq{
		Data: result.Permits,
	})
	if err != nil {
		logger.Logger.WithName("import permit").Errorf("import permit failed: %v", err)

		return err
	}

	return nil
}

type role struct {
	appID     string
	exportURL string
	importURL string

	client http.Client
}

func (t *role) Initialize(appID string, client http.Client) {
	t.client = client
	t.appID = appID
	t.exportURL = fmt.Sprintf(exportRoleURL, formHost, appID)
	t.importURL = fmt.Sprintf(importRoleURL, formHost, appID)
}

func (t *role) Export(ctx context.Context, result *Result) error {
	data, err := export(ctx, &t.client, t.exportURL, BackupReq{
		AppID: t.appID,
		Page:  startPage,
		Size:  maxSize,
	})
	if err != nil {
		logger.Logger.WithName("export role").Errorf("export role failed: %v", err)

		return err
	}

	result.Roles = data

	return nil
}

func (t *role) Import(ctx context.Context, result *Result) error {
	err := import2(ctx, &t.client, t.importURL, ImportReq{
		Data: result.Roles,
	})
	if err != nil {
		logger.Logger.WithName("import role").Errorf("import role failed: %v", err)

		return err
	}

	return nil
}

func export(ctx context.Context, cli *http.Client, url string, req BackupReq) (Object, error) {
	totalPage := 0
	data := make(Object, 0)
	for {
		resp := &BackupResp{}

		err := client.POST(ctx, cli, url, req, resp)
		if err != nil {
			logger.Logger.WithName("export request").Errorf("send http request failed: %v", err)

			return nil, err
		}

		data = append(data, resp.Data...)

		if resp.Count <= req.Size {
			break
		}

		if totalPage == 0 {
			if resp.Count%req.Size == 0 {
				totalPage = resp.Count / req.Size
			} else {
				totalPage = resp.Count/req.Size + 1
			}
		}

		if totalPage <= req.Page {
			break
		}

		req.Page++
	}

	return data, nil
}

func import2(ctx context.Context, cli *http.Client, url string, req ImportReq) error {
	var index int
	if len(req.Data)%maxSize == 0 {
		index = len(req.Data) / maxSize
	} else {
		index = len(req.Data)/maxSize + 1
	}

	for i := 0; i < index; i++ {
		if i == index-1 {
			req.Data = req.Data[i*maxSize:]
		} else {
			req.Data = req.Data[i*maxSize : (i+1)*maxSize]
		}

		err := client.POST(ctx, cli, url, req, &ImportResp{})
		if err != nil {
			logger.Logger.WithName("export request").Errorf("send http request failed: %v", err)

			return err
		}
	}

	return nil
}
