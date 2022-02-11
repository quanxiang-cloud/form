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

type FindOptions struct {
	Page int64    `json:"page"`
	Size int64    `json:"size"`
	Sort []string `json:"sort"`
}

type FormReq struct {
	FindOptions
	DslQuery interface{}
	Entity   interface{}
	TableID  string
}

type StructorReq struct {
	FindOptions
	TableID string
	Dsl     *anypb.Any
	Entity  *anypb.Any
}

type FormResp struct {
	Aggregations interface{}              `json:"aggregations"`
	Entities     []map[string]interface{} `json:"entities"`
	Total        int64                    `json:"total"`
	Entity       map[string]interface{}   `json:"entity"`
	SuccessCount int64                    `json:"successCount"`
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

func getStructorReq(req *FormReq) (*StructorReq, error) {
	structor := &StructorReq{}
	structor.Size = req.Size
	structor.Page = req.Page
	structor.Sort = req.Sort
	if req.DslQuery != nil {
		marshal, err := json.Marshal(req.DslQuery)
		if err != nil {
			return nil, err
		}
		any, err := rawToAny(marshal)
		if err != nil {
			return nil, err
		}
		structor.Dsl = any
	}
	if req.Entity != nil {
		marshal, err := json.Marshal(req.Entity)
		if err != nil {
			return nil, err
		}
		any, err := rawToAny(marshal)
		if err != nil {
			return nil, err
		}
		structor.Entity = any
	}

	return structor, nil
}

type SearchResp struct {
	Entities []map[string]interface{} `json:"entities"`
	Total    int64                    `json:"total"`
}

func (f *FormAPI) Search(ctx context.Context, formReq *FormReq) (*SearchResp, error) {
	req, err := getStructorReq(formReq)
	if err != nil {
		return nil, err
	}
	searchResp, err := f.client.Find(ctx, &pb.FindReq{
		TableName: req.TableID,
		Dsl:       req.Dsl,
		Page:      req.Page,
		Size:      req.Size,
		Sort:      req.Sort,
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

type UpdateResp struct {
	SuccessCount int64 `json:"successCount"`
}

func (f *FormAPI) Update(ctx context.Context, formReq *FormReq) (*UpdateResp, error) {
	req, err := getStructorReq(formReq)
	if err != nil {
		return nil, err
	}
	update, err := f.client.Update(ctx, &pb.UpdateReq{
		Entity:    req.Entity,
		Dsl:       req.Dsl,
		TableName: req.TableID,
	})
	if err != nil {
		return nil, err
	}
	return &UpdateResp{
		SuccessCount: update.Count,
	}, nil
}

type InsertResp struct {
	SuccessCount int64 `json:"successCount"`
}

func (f *FormAPI) Insert(ctx context.Context, formReq *FormReq) (*InsertResp, error) {
	req, err := getStructorReq(formReq)
	if err != nil {
		return nil, err
	}
	anyArr := make([]*anypb.Any, 0)
	anyArr = append(anyArr, req.Entity)
	insert, err := f.client.Insert(ctx, &pb.InsertReq{
		TableName: req.TableID,
		Entities:  nil,
	})
	if err != nil {
		return nil, err
	}

	return &InsertResp{
		SuccessCount: insert.Count,
	}, err
}

type GetResp struct {
	Entity map[string]interface{} `json:"entity"`
}

func (f *FormAPI) Get(ctx context.Context, formReq *FormReq) (*GetResp, error) {
	req, err := getStructorReq(formReq)
	if err != nil {
		return nil, err
	}
	getResp, err := f.client.FindOne(ctx, &pb.FindOneReq{
		TableName: req.TableID,
		Dsl:       req.Dsl,
	})
	if err != nil {
		return nil, err
	}
	data, err := anyToRaw(getResp.GetData())
	if err != nil {
		return nil, err
	}
	var entity map[string]interface{}
	err = json.Unmarshal(data, &entity)
	if err != nil {
		return nil, err
	}
	return &GetResp{
		Entity: entity,
	}, nil
}

type DeleteResp struct {
	SuccessCount int64 `json:"successCount"`
}

func (f *FormAPI) Delete(ctx context.Context, formReq *FormReq) (*DeleteResp, error) {
	req, err := getStructorReq(formReq)
	if err != nil {
		return nil, err
	}
	deleteResp, err := f.client.Delete(ctx, &pb.DeleteReq{
		Dsl:       req.Dsl,
		TableName: req.TableID,
	})
	if err != nil {
		return nil, err
	}
	return &DeleteResp{
		SuccessCount: deleteResp.Count,
	}, nil

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
