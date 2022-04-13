package aide

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/quanxiang-cloud/cabin/logger"
	"github.com/quanxiang-cloud/cabin/tailormade/client"
)

const (
	startPage = 1
	maxSize   = 1000
)

var defaultReq = ExportReq{
	Page: startPage,
	Size: maxSize,
}

// Aide is the interface of aide.
type Aide interface {
	Export(ctx context.Context, opts *ExportOption) (map[string]Object, error)
	Import(ctx context.Context, objs map[string]Object, opts *ImportOption) (map[string]string, error)
}

// Object is a slice of interface{}.
type Object []interface{}

// ExportOption is the option of export.
type ExportOption struct {
	AppID string `required:"true"`

	// these parameters do not need to be passed
	Host   string
	Client http.Client
}

// ExportReq is the request of export.
type ExportReq struct {
	Page int `json:"page"`
	Size int `json:"size"`
}

// ExportResp is the response of export.
type ExportResp struct {
	Data  Object `json:"data"`
	Count int    `json:"count"`
}

// ImportOption is the option of import.
type ImportOption struct {
	AppID    string `required:"true"`
	UserID   string `required:"true"`
	UserName string `required:"true"`

	// these parameters do not need to be passed
	Host   string
	Client http.Client
}

// ExportObject export object.
func ExportObject(ctx context.Context, url string, opts *ExportOption) (Object, error) {
	var (
		totalPage = 0
		req       = defaultReq
		result    = make(Object, 0)
	)

	for {
		resp := &ExportResp{}
		err := client.POST(ctx, &opts.Client, url, req, resp)
		if err != nil {
			logger.Logger.WithName("export request").Errorf("send http request failed: %v", err)

			return nil, err
		}

		result = append(result, resp.Data...)

		if resp.Count <= req.Size {
			break
		}

		// count is greater than size
		if totalPage == 0 {
			totalPage = callTotalPage(resp.Count, req.Size)
		}

		req.Page++
	}

	return result, nil
}

// ImportReq is the request of import.
type ImportReq struct {
	Data Object `json:"data"`
}

// ImportResp is the response of import.
type ImportResp struct{}

// ImportObject import object.
func ImportObject(ctx context.Context, url string, data Object, opts *ImportOption) error {
	var (
		totalPage = 1
		req       = ImportReq{}
	)

	if len(data) > maxSize {
		totalPage = callTotalPage(len(data), maxSize)
	}

	for i := 0; i < totalPage; i++ {
		req.Data = callImportData(data, i, totalPage)

		err := client.POST(ctx, &opts.Client, url, req, &ImportResp{})
		if err != nil {
			logger.Logger.WithName("export request").Errorf("send http request failed: %v", err)

			return err
		}
	}

	return nil
}

func callTotalPage(count, size int) int {
	if (count % size) == 0 {
		return count / size
	}

	return (count / size) + 1
}

func callImportData(data Object, index, page int) Object {
	if index == page-1 {
		return data[index*maxSize:]
	}

	return data[(index * maxSize):((index + 1) * maxSize)]
}

// Serialize serialize object.
func Serialize(data interface{}, obj interface{}) error {
	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return json.Unmarshal(bytes, obj)
}
