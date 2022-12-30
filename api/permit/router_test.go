package router

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"testing"
	"time"
)

func TestName(t *testing.T) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints: []string{"192.168.208.4:2379", "192.168.208.6:2379", "192.168.208.15:2379"},
		Username:  "root",
		Password:  "9yXpYLPIX4",

		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		// handle error!
		fmt.Printf("connect to etcd failed, err:%v\n", err)
		return
	}
	fmt.Println("connect to etcd success")
	defer cli.Close()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	_, err = cli.Put(ctx, "q1mi", "dsb")
	cancel()
	if err != nil {
		fmt.Printf("put to etcd failed, err:%v\n", err)
		return
	}
	// get
	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	//, err := cli.Get(ctx, "q1mi")
	//cancel()
	//if err != nil {
	//	fmt.Printf("get from etcd failed, err:%v\n", err)
	//	return
	//}

}
