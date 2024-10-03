package config

import (
	"os"
	"time"

	"studentgit.kata.academy/quant/torque/pkg/kafka"

	"gopkg.in/yaml.v3"
	"studentgit.kata.academy/quant/torque/internal/domain"
	"studentgit.kata.academy/quant/torque/pkg/logster"
	"studentgit.kata.academy/quant/torque/pkg/postgres"
)

type HTTPClient struct {
	APIKey        string        `yaml:"apiKey"`
	APISecret     string        `yaml:"apiSecret"`
	APIPassphrase string        `yaml:"apiPassphrase"`
	URL           string        `yaml:"url"`
	Timeout       time.Duration `yaml:"timeout"`
}

type MarketDataConfig struct {
	URL                    string   `yaml:"url"`
	InstrumentsToSubscribe []string `yaml:"instrumentsToSubscribe"`
}

type UserStreamConfig struct {
	APIKey        string `yaml:"apiKey"`
	APISecret     string `yaml:"apiSecret"`
	APIPassphrase string `yaml:"apiPassphrase"`
	URL           string `yaml:"url"`
}

type Config struct {
	HTTPClient       HTTPClient             `yaml:"httpClient"`
	MarketDataConfig MarketDataConfig       `yaml:"marketDataConfig"`
	UserStreamConfig UserStreamConfig       `yaml:"userStreamConfig"`
	TradingWebsocket UserStreamConfig       `yaml:"tradingWebsocket"`
	ExchangeAccount  domain.ExchangeAccount `yaml:"exchangeAccount"`
	NameResolver     map[string]string      `yaml:"nameResolver"`
	Piston           string                 `yaml:"piston"`
	PrivateAddr      string                 `yaml:"privateAddr"`
	GRPCApiAddr      string                 `yaml:"GRPCApiAddr"`
	Postgres         postgres.Config        `yaml:"postgres"`
	Log              logster.Config         `yaml:"log"`
	Kafka            KafkaConfig            `yaml:"kafka"`
}

type KafkaConfig struct {
	Consumers Consumers `yaml:"consumers"`
	Producers Producers `yaml:"producers"`
}

type Consumers struct {
	Internal      *kafka.ConsumerConfig `yaml:"internal"`
	GroupInternal *kafka.ConsumerConfig `yaml:"group_internal"`
}

type Producers struct {
	Internal *kafka.ProducerConfig `yaml:"internal"`
}

func LoadConfig(filename string, cfg interface{}) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		return err
	}

	return nil
}
