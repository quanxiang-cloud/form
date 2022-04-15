package client

import (
	"context"
	pb "github.com/quanxiang-cloud/structor/api/proto"
	"google.golang.org/grpc"
)

type FormDDLAPI struct {
	client pb.DDLServiceClient
}

func NewFormDDLAPI() (*FormDDLAPI, error) {
	client, err := connectDDL(target)
	if err != nil {
		return nil, err
	}

	return &FormDDLAPI{
		client: client,
	}, nil
}

func connectDDL(target string) (pb.DDLServiceClient, error) {
	conn, err := grpc.Dial(target, grpc.WithInsecure())

	if err != nil {
		return nil, err
	}
	return pb.NewDDLServiceClient(conn), nil
}

type CreateTableResp struct {
	TableName string
}

type Field struct {
	Title string
	Type  string
	Max   int64
}

func (f *FormDDLAPI) CreateTable(ctx context.Context, tableName string, field []*Field) (*CreateTableResp, error) {
	req := &pb.CreateReq{
		Fields:    toPbField(field),
		TableName: tableName,
	}
	createResp, err := f.client.Create(ctx, req)
	if err != nil {
		return nil, err
	}
	return &CreateTableResp{
		TableName: createResp.TableName,
	}, nil
}

// 初始化的创建

// 修改表的结构

// 不支持改数据类型 ，

type AlterADDResp struct {
	TableName string
}

func (f *FormDDLAPI) AlterADD(ctx context.Context, tableName string, field []*Field) (*AlterADDResp, error) {
	req := &pb.AddReq{
		Fields:    toPbField(field),
		TableName: tableName,
	}
	createResp, err := f.client.Add(ctx, req)
	if err != nil {
		return nil, err
	}
	return &AlterADDResp{
		TableName: createResp.TableName,
	}, nil
}

type IndexResp struct {
	IndexName string
}

func (f *FormDDLAPI) Index(ctx context.Context, tableID, fieldName, indexName string) (*IndexResp, error) {
	req := &pb.IndexReq{
		TableName: tableID,
		IndexName: indexName,
		Titles:    []string{fieldName},
	}
	indexResp, err := f.client.Index(ctx, req)
	if err != nil {
		return nil, err
	}

	return &IndexResp{
		IndexName: indexResp.IndexName,
	}, nil

}

func toPbField(field []*Field) []*pb.Field {
	fields := make([]*pb.Field, len(field))
	for index, value := range field {
		fields[index] = &pb.Field{
			Title: value.Title,
			Type:  value.Type,
			Max:   value.Max,
		}
	}
	return fields
}
