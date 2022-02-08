package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/quanxiang-cloud/cabin/logger"
	"github.com/quanxiang-cloud/form/pkg/misc/config"
	pb "github.com/quanxiang-cloud/structor/api/proto"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/structpb"
	"testing"
)

func TestSearch(t *testing.T) {

	newConfig, err := config.NewConfig("../../configs/config.yml")
	if err != nil {
		logger.Logger.Errorw("error")
	}
	newForm, err := NewForm(newConfig)
	if err != nil {
		logger.Logger.Errorw("error")
	}
	search, err := newForm.Search(context.Background(), &SearchReq{
		IsAuth:  false,
		TableID: "7p9rb",
		AppID:   "kg4r5",
		FindOptions: FindOptions{
			Page: 1,
			Size: 10,
		},
		Query: map[string]interface{}{
			"term": map[string]interface{}{
				"creator_id": "f253a657-367e-4d7f-a815-94c43e327b04",
			},
		},
	})
	if err != nil {
		logger.Logger.Errorw("error")
	}
	logger.Logger.Info(search)
}

func Test1(t *testing.T) {
	c, err := getConn("localhost:80")
	if err != nil {
		panic(err)
	}

	find(c)
	// insert(c)
	// update(c)
	// delete(c)
	// findOne(c)
}

func getConn(addr string) (pb.DSLServiceClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	dslConn := pb.NewDSLServiceClient(conn)
	return dslConn, nil
}

func find(c pb.DSLServiceClient) {

	any, err := rawToAny([]byte(`
	{
		"query": {}
	}
	`))
	if err != nil {
		panic(err)
	}

	resp, err := c.Find(context.Background(), &pb.FindReq{
		TableName: "user",
		Dsl:       any,
		Page:      1,
		Size:      3,
	})
	if err != nil {
		panic(err)
	}

	data, err := anyToRaw(resp.GetData())
	if err != nil {
		panic(err)
	}

	fmt.Println(string(data))
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
