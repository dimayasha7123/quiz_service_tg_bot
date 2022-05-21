package config

import (
	"gopkg.in/yaml.v2"
)

// я понимаю, что ключи нельзя так хранить, но в требованиях написано, что нужно использовать yaml
type ApiKeys struct {
	Telegram string
	Quiz     string
}

type Config struct {
	ApiKeys  ApiKeys
	Dsn      string
	QuizTags []string
}

type configFile struct {
	APIKeys struct {
		Telegram string `yaml:"telegram"`
		Quiz     string `yaml:"quiz"`
	} `yaml:"apiKeys"`
	QuizTags []string `yaml:"quizTags"`
	Dsn      string   `yaml:"dsn"`
}

func ParseConfig(fileBytes []byte) (*Config, error) {
	cf := configFile{}
	err := yaml.Unmarshal(fileBytes, &cf)
	if err != nil {
		return nil, err
	}

	c := Config{}

	c.ApiKeys.Telegram = cf.APIKeys.Telegram
	c.ApiKeys.Quiz = cf.APIKeys.Quiz
	c.Dsn = cf.Dsn

	c.QuizTags = make([]string, len(cf.QuizTags))
	for i, q := range cf.QuizTags {
		c.QuizTags[i] = q
	}

	return &c, nil
}

//func ConfigWriter(c *Config) ([]byte, error) {
//	cf := configFile{}
//	cf.APIKeys.Telegram = c.ApiKeys.Telegram
//	cf.APIKeys.Quiz = c.ApiKeys.Quiz
//	cf.APIKeys.Dsn = c.ApiKeys.Dsn
//
//	cf.QuizTags = make([]string, len(c.QuizTags))
//	for i, q := range c.QuizTags {
//		cf.QuizTags[i] = q
//	}
//
//	bytes, err := yaml.Marshal(cf)
//	if err != nil {
//		return nil, err
//	}
//
//	return bytes, nil
//}
