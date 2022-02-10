package client

import (
	"context"
	"encoding/json"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/structpb"

	pb "github.com/quanxiang-cloud/structor/api/proto"
)

const (
	target = "localhost:80"
)

type FormAPI struct {
	client pb.DSLServiceClient
}

func NewFormAPI() (*FormAPI, error) {
	client, err := connect(target)
	if err != nil {
		return nil, err
	}

	return &FormAPI{
		client: client,
	}, nil
}

func connect(target string) (pb.DSLServiceClient, error) {
	conn, err := grpc.Dial(target, grpc.WithInsecure())

	if err != nil {
		return nil, err
	}
	return pb.NewDSLServiceClient(conn), nil
}

type SearchResp struct {
	Aggregations interface{}              `json:"aggregations"`
	Entities     []map[string]interface{} `json:"entities"`
	Total        int64                    `json:"total"`
}

func (f *FormAPI) Search(ctx context.Context, options FindOptions, dsl interface{}, tableName string) (*SearchResp, error) {
	marshal, err := json.Marshal(dsl)
	if err != nil {
		return nil, err
	}
	any, err := rawToAny(marshal)
	if err != nil {
		return nil, err
	}

	searchResp, err := f.client.Find(ctx, &pb.FindReq{
		TableName: tableName,
		Dsl:       any,
		Page:      options.Page,
		Size:      options.Size,
		Sort:      options.Sort,
	})
	if err != nil {
		return nil, err
	}
	data, err := anyToRaw(searchResp.GetData())
	if err != nil {
		return nil, err
	}
	var entity []map[string]interface{}
	err = json.Unmarshal(data, &entity)
	if err != nil {
		return nil, err
	}
	return &SearchResp{
		Entities: entity,
	}, nil
}

type InsertResp struct {
	SuccessCount int64 `json:"count"`
}

// Insert insert
func (f *FormAPI) Insert(ctx context.Context, tableName string, entity interface{}) (*InsertResp, error) {
	marshal, err := json.Marshal(entity)
	if err != nil {
		return nil, err
	}
	any, err := rawToAny(marshal)
	if err != nil {
		return nil, err
	}
	anyArr := make([]*anypb.Any, 0)
	anyArr = append(anyArr, any)
	insert, err := f.client.Insert(ctx, &pb.InsertReq{
		TableName: tableName,
		Entities:  anyArr,
	})
	if err != nil {
		return nil, err
	}
	return &InsertResp{
		SuccessCount: insert.Count,
	}, err
}

type UpdateResp struct {
	Count int64 `json:"count"`
}

func (f *FormAPI) Update(ctx context.Context, entity map[string]interface{}, dsl map[string]interface{}) (*UpdateResp, error) {
	//dslAny, err := MarshalAny(dsl)
	//if err != nil {
	//	return nil ,err
	//}
	//entityAny, err := MarshalAny(entity)
	//if err != nil {
	//	return nil ,err
	//}
	update, err := f.client.Update(ctx, &pb.UpdateReq{
		//Entity: entityAny,
		//Dsl:    dslAny,
	})
	if err != nil {
		return nil, err
	}
	return &UpdateResp{
		Count: update.Count,
	}, nil

}

func (f *FormAPI) Delete(ctx context.Context, dsl map[string]interface{}) {

}

type FindOptions struct {
	Page int64    `json:"page"`
	Size int64    `json:"size"`
	Sort []string `json:"sort"`
}

func anyToRaw(any *anypb.Any) (json.RawMessage, error) {
	out := structpb.NewNullValue()
	err := any.UnmarshalTo(out)
	if err != nil {
		return nil, err
	}

	body, err := out.MarshalJSON()
	if err != nil {
		return nil, err
	}
	return body, nil
}

func rawToAny(raw []byte) (*anypb.Any, error) {
	in := structpb.NewNullValue()
	err := in.UnmarshalJSON(raw)
	if err != nil {
		return nil, err
	}
	any := &anypb.Any{}
	err = any.MarshalFrom(in)
	return any, err
}
