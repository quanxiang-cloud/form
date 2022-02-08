package main

import (
	"flag"
	"github.com/quanxiang-cloud/cabin/logger"
	"github.com/quanxiang-cloud/form/api"
	"github.com/quanxiang-cloud/form/pkg/misc/config"
	"os"
	"os/signal"
	"syscall"
)

var (
	configPath = flag.String("config", "../configs/config.yml", "-config 配置文件地址")
)

func main() {
	flag.Parse()
	conf, err := config.NewConfig(*configPath)
	if err != nil {
		panic(err)
	}
	logger.Logger = logger.New(&conf.Log)
	if err != nil {
		panic(err)
	}
	router, err := api.NewRouter(conf)
	if err != nil {
		panic(err)
	}
	go router.Run()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			//	router.Close()
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}

}
