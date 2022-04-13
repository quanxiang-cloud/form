package backup

import (
	"context"
	"fmt"

	"github.com/quanxiang-cloud/cabin/logger"
	"github.com/quanxiang-cloud/cabin/tailormade/client"
)

type role struct{}

func (t *role) Export(ctx context.Context, result *Result, opts *ExportOption) error {
	url := fmt.Sprintf(exportRoleURL, formHost, opts.AppID)

	return exportRoleFunc(ctx, url,
		ExportReq{
			AppID: opts.AppID,
			Page:  startPage,
			Size:  maxSize,
		},
		result, opts,
	)
}

func (t *role) Import(ctx context.Context, result *Result, opts *ImportOption) error {
	var (
		index int
		req   ImportReq
		url   = fmt.Sprintf(importRoleURL, formHost, opts.AppID)
	)

	if len(result.Tables)%maxSize == 0 {
		index = len(result.Tables) / maxSize
	} else {
		index = len(result.Tables)/maxSize + 1
	}

	for i := 0; i < index; i++ {
		if i == index-1 {
			req.Tables = result.Tables[i*maxSize:]
		} else {
			req.Tables = result.Tables[i*maxSize : (i+1)*maxSize]
		}

		err := client.POST(ctx, &opts.Client, url, req, &ImportResp{})
		if err != nil {
			logger.Logger.WithName("export request").Errorf("send http request failed: %v", err)

			return err
		}
	}

	return nil
}

func exportRoleFunc(ctx context.Context, url string, req ExportReq, result *Result, opts *ExportOption) error {
	totalPage := 0

	for {
		resp := &ExportResp{}

		err := client.POST(ctx, &opts.Client, url, req, resp)
		if err != nil {
			logger.Logger.WithName("export request").Errorf("send http request failed: %v", err)

			return err
		}

		result.Roles = append(result.Roles, resp.Roles...)

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
