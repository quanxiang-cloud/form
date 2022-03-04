package inform

import (
	"context"
	"encoding/json"
	"git.internal.yunify.com/qxp/misc/logger"
	"github.com/Shopify/sarama"
	kafka2 "github.com/quanxiang-cloud/cabin/tailormade/db/kafka"
	"github.com/quanxiang-cloud/form/pkg/misc/config"
)

const (
	topic = "form-data"
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
	Send chan *FormData // 增删改数据后，放到这个信道

	kafkaProducer sarama.SyncProducer
}

// NewHookManger NewHookManger
func NewHookManger(ctx context.Context) (*HookManger, error) {
	producer, err := kafka2.NewSyncProducer(config.Conf.Kafka)
	if err != nil {
		return nil, err
	}
	m := &HookManger{
		Send:          make(chan *FormData),
		kafkaProducer: producer,
	}
	go m.Start(ctx)
	return m, nil

}

// Start Start
func (manager *HookManger) Start(ctx context.Context) {
	for {
		select {
		case sendData := <-manager.Send:
			value, err := json.Marshal(sendData)
			if err != nil {
				logger.Logger.Error(err.Error()+sendData.TableID, logger.STDRequestID(ctx))
			}
			message := sarama.ProducerMessage{
				Topic: topic,
				Value: sarama.ByteEncoder(value),
			}
			_, _, err = manager.kafkaProducer.SendMessage(&message)
			if err != nil {
				logger.Logger.Error(err.Error()+sendData.TableID, logger.STDRequestID(ctx))
			}

		case <-ctx.Done():

		}
	}

}
