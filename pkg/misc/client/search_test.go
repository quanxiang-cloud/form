package client

import (
	"context"
	"flag"
	"github.com/quanxiang-cloud/cabin/logger"
	"github.com/quanxiang-cloud/form/pkg/misc/config"
	"testing"
)

var (
	configPath = flag.String("config", "../../../configs/config.yml", "-config 配置文件地址")
)

func TestName(t *testing.T) {
	flag.Parse()
	conf, err := config.NewConfig(*configPath)
	if err != nil {
		panic(err)
	}
	api := NewSearchAPI(conf)
	subordinate, err := api.Subordinate(context.Background(), "f44e3fd5-1f4e-4c01-a799-d43ac1212ef8")
	if err != nil {

	}
	logger.Logger.Infow("sub is      ------", subordinate)

}
