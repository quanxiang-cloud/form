package inform

import "github.com/Shopify/sarama"

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
