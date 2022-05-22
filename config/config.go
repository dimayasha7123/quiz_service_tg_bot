package config

import (
	"gopkg.in/yaml.v2"
)

type Config struct {
	Socket         Socket
	TelegramAPIKey string
	Dsn            string
}

type Socket struct {
	Host     string
	GrpcPort string
}

type configFile struct {
	Socket struct {
		Host     string `yaml:"host"`
		GrpcPort string `yaml:"grpcPort"`
	} `yaml:"socket"`
	TelegramAPIKey string `yaml:"telegramApiKey"`
	Dsn            string `yaml:"dsn"`
}

func ParseConfig(fileBytes []byte) (*Config, error) {
	cf := configFile{}
	err := yaml.Unmarshal(fileBytes, &cf)
	if err != nil {
		return nil, err
	}

	c := Config{}

	c.TelegramAPIKey = cf.TelegramAPIKey
	c.Dsn = cf.Dsn
	c.Socket = Socket{
		Host:     cf.Socket.Host,
		GrpcPort: cf.Socket.GrpcPort,
	}

	return &c, nil
}
