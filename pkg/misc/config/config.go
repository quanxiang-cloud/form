package config

import (
	mongo2 "github.com/quanxiang-cloud/cabin/tailormade/db/mongo"
	"io/ioutil"
	"time"

	"github.com/quanxiang-cloud/cabin/logger"
	"github.com/quanxiang-cloud/cabin/tailormade/client"
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
	Mongo       mongo2.Config `yaml:"mongo"`
	Redis       redis2.Config `yaml:"redis"`
	Endpoint    Endpoint      `yaml:"endpoint"`
	Transport   Transport     `yaml:"transport"`
	Dapr        Dapr          `yaml:"dapr"`
}

type Dapr struct {
	PubSubName string `yaml:"pubSubName"`
	TopicFlow  string `yaml:"topicFlow"`
}

type Endpoint struct {
	Poly      string `yaml:"poly"`
	Form      string `yaml:"form"`
	FormInner string `yaml:"formInner"`
	PolyInner string `yaml:"polyInner"`
	Org       string `yaml:"org"`
	AppCenter string `yaml:"appCenter"`
	Search    string `yaml:"search"`
	Structor  string `yaml:"structor"`
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
