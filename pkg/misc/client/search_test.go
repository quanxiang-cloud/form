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
	subordinate, err := api.Subordinate(context.Background(), "e0fbb3a2-ad03-4fc5-af7d-96a59267fdf8")
	if err != nil {

	}
	logger.Logger.Infow("sub is      ------", subordinate)

}
