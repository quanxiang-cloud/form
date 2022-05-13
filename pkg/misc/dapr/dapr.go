package dapr

import (
	daprd "github.com/dapr/go-sdk/client"
	"github.com/quanxiang-cloud/cabin/logger"
	"os"
	"sync"
	"time"
)

const (
	daprPortEnvVarName = "DAPR_GRPC_PORT" /* #nosec */
	daprPortDefault    = "50001"
)

var (
	daprClient daprd.Client
	doOnce     sync.Once
	mu         sync.Mutex
)

func InitDaprClientIfNil() (daprd.Client, error) {
	port := os.Getenv(daprPortEnvVarName)
	if port == "" {
		port = daprPortDefault
	}
	if daprClient == nil {
		var err error
		mu.Lock()
		defer mu.Unlock()
		for attempts := 120; attempts > 0; attempts-- {
			c, e := daprd.NewClientWithPort(port)
			if e == nil {
				daprClient = c
				break
			}
			err = e
			time.Sleep(500 * time.Millisecond)
		}
		if daprClient == nil {
			logger.Logger.Errorf("failed to init dapr client: %v", err)
			return nil, err
		}
	}
	return daprClient, nil
}
