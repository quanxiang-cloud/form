package config

import (
	"io/ioutil"
	"time"

	"github.com/quanxiang-cloud/cabin/logger"
	"github.com/quanxiang-cloud/cabin/tailormade/client"
	"github.com/quanxiang-cloud/cabin/tailormade/db/kafka"
	mysql2 "github.com/quanxiang-cloud/cabin/tailormade/db/mysql"
	redis2 "github.com/quanxiang-cloud/cabin/tailormade/db/redis"
	"gopkg.in/yaml.v2"
)

// Conf 配置文件
var Conf *Config

// DefaultPath 默认配置路径
var DefaultPath = "./configs/config.yml"

// Config 配置文件
type Config struct {
	InternalNet client.Config `yaml:"internalNet"`
	PortInner   string        `yaml:"portInner"`
	Port        string        `yaml:"port"`
	Model       string        `yaml:"model"`
	Log         logger.Config `yaml:"log"`
	Mysql       mysql2.Config `yaml:"mysql"`
	PubSubName  string        `yaml:"pubSubName"`
	Redis       redis2.Config `yaml:"redis"`
	Kafka       kafka.Config  `yaml:"kafka"`
	SwaggerPath string        `yaml:"swaggerPath"`
	Endpoint    Endpoint      `yaml:"endpoint"`
	Transport   Transport     `yaml:"transport"`
}

// Service service config
type Service struct {
	DB string `yaml:"db"`
}

type Endpoint struct {
	Poly string `yaml:"poly"`
	Form string `yaml:"form"`
}

type Transport struct {
	Timeout               time.Duration `yaml:"timeout"`
	KeepAlive             time.Duration `yaml:"keepAlive"`
	MaxIdleConns          int           `yaml:"maxIdleConns"`
	IdleConnTimeout       time.Duration `yaml:"idleConnTimeout"`
	TLSHandshakeTimeout   time.Duration `yaml:"tlsHandshakeTimeout"`
	ExpectContinueTimeout time.Duration `yaml:"expectContinueTimeout"`
}

// NewConfig 获取配置配置
func NewConfig(path string) (*Config, error) {
	if path == "" {
		path = DefaultPath
	}

	file, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(file, &Conf)
	if err != nil {
		return nil, err
	}

	return Conf, nil
}
