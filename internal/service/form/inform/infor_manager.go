package inform

import (
	"context"
	"fmt"
	daprd "github.com/dapr/go-sdk/client"
	"github.com/go-logr/logr"
	"github.com/quanxiang-cloud/form/pkg/misc/config"
)

const (
	topic = "form-data4"
)

// FormData FormData
type FormData struct {
	TableID string      `json:"tableID"`
	Entity  interface{} `json:"entity"`
	Magic   string      `json:"magic"`
	Seq     string      `json:"seq"`
	Version string      `json:"version"`
	Method  string      `json:"method"`
}

// HookManger 管理发送kafka
type HookManger struct {
	Send       chan *FormData // 增删改数据后，放到这个信道
	conf       *config.Config
	daprClient daprd.Client
	log        logr.Logger
}

// NewHookManger NewHookManger
func NewHookManger(ctx context.Context, conf *config.Config) (*HookManger, error) {

	client, err := daprd.NewClient()
	if err != nil {
		return nil, err
	}
	m := &HookManger{
		daprClient: client,
		Send:       make(chan *FormData),
		conf:       conf,
	}
	go m.Start(ctx)
	return m, nil

}

// Start Start
func (manager *HookManger) Start(ctx context.Context) {
	for {
		select {
		case sendData := <-manager.Send:
			fmt.Println("121212121")
			if err := manager.publish(ctx, topic, sendData); err != nil {
				manager.log.Error(err, "push flow", "sendData ", sendData)
			}
			fmt.Println("1212121")
		case <-ctx.Done():

		}
	}

}
func (manager *HookManger) publish(ctx context.Context, topic string, data interface{}) error {
	//	manager.log.Info("send message", "data is ", data)
	fmt.Println("2323232")
	if err := manager.daprClient.PublishEvent(context.Background(), manager.conf.PubSubName, topic, data); err != nil {
		manager.log.Error(err, "publishEvent", "topic", topic, "pubsubName", manager.conf.PubSubName)
		return err
	}
	fmt.Println("32323232")
	return nil
}
