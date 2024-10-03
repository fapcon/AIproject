package kafka

import (
	"fmt"

	"github.com/IBM/sarama"
	"studentgit.kata.academy/quant/torque/pkg/logster"
	"studentgit.kata.academy/quant/torque/pkg/scram"
)

type ProducerConfig struct {
	ClientID string   `yaml:"client_id"`
	Topic    string   `yaml:"topic"`
	Brokers  []string `yaml:"brokers"`

	EnableAuth bool   `yaml:"enable_auth"`
	User       string `yaml:"user"`
	Password   string `yaml:"password"`
}

func (c ProducerConfig) String() string {
	return fmt.Sprintf(
		`{ClientID:%s Topic:%v Brokers:%v EnableAuth: %t}`,
		c.ClientID, c.Topic, c.Brokers, c.EnableAuth)
}

func (c ProducerConfig) toSarama() *sarama.Config {
	config := sarama.NewConfig()

	config.Version = sarama.V2_2_0_0
	config.ClientID = c.ClientID

	config.Metadata.Full = false

	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 10
	config.Producer.Return.Errors = true

	if c.EnableAuth {
		config.Net.SASL.Mechanism = sarama.SASLTypeSCRAMSHA512
		config.Net.SASL.Enable = true
		config.Net.SASL.User = c.User
		config.Net.SASL.Password = c.Password
		config.Net.SASL.SCRAMClientGeneratorFunc = scram.NewSHA512Client
	}
	return config
}

type SyncProducerConfig ProducerConfig

func (c SyncProducerConfig) toSarama() *sarama.Config {
	config := ProducerConfig(c).toSarama()
	config.Producer.Return.Successes = true
	return config
}

type Producer[K, T any] struct {
	producer sarama.AsyncProducer
	topic    string
	encode   func(K, T) (key, value []byte, _ []sarama.RecordHeader, _ error)
}

func NewProducer[K, T any](
	logger logster.Logger,
	config ProducerConfig,
	encode func(K, T) (key []byte, value []byte, _ []sarama.RecordHeader, _ error),
) (*Producer[K, T], error) {
	producer, err := sarama.NewAsyncProducer(config.Brokers, config.toSarama())
	if err != nil {
		return nil, fmt.Errorf("can't create async producer: %w", err)
	}

	p := &Producer[K, T]{
		producer: producer,
		topic:    config.Topic,
		encode:   encode,
	}

	logger = logger.WithFields(logster.Fields{
		"topic":        config.Topic,
		"message_type": messageType[T](),
	})

	go func() {
		for e := range producer.Errors() {
			key, _ := e.Msg.Key.Encode()
			logger.WithError(e).
				WithField("key", key).
				Errorf("%T: can't send message", *p)
		}
	}()

	return p, nil
}

func (p *Producer[K, T]) Topic() string {
	return p.topic
}

func (p *Producer[K, T]) Close() error {
	return p.producer.Close()
}

func (p *Producer[K, T]) Send(k K, message T) {
	key, value, headers, err := p.encode(k, message)
	if err != nil {
		p.producer.Input() <- &sarama.ProducerMessage{
			Topic: p.topic,
			Value: failedEncoder{err: fmt.Errorf("encode: %w", err)},
		}
		return
	}

	msg := &sarama.ProducerMessage{
		Topic:   p.topic,
		Value:   sarama.ByteEncoder(value),
		Headers: headers,
	}
	if len(key) > 0 {
		msg.Key = sarama.ByteEncoder(key)
	}

	p.producer.Input() <- msg
}

type SyncProducer[K, T any] struct {
	producer sarama.SyncProducer
	topic    string
	encode   func(K, T) (key, value []byte, _ []sarama.RecordHeader, _ error)
}

func NewSyncProducer[K, T any](
	config SyncProducerConfig,
	encode func(K, T) (key []byte, value []byte, _ []sarama.RecordHeader, _ error),
) (*SyncProducer[K, T], error) {
	producer, err := sarama.NewSyncProducer(config.Brokers, config.toSarama())
	if err != nil {
		return nil, fmt.Errorf("can't create sync producer: %w", err)
	}

	return &SyncProducer[K, T]{
		producer: producer,
		topic:    config.Topic,
		encode:   encode,
	}, nil
}

func (p *SyncProducer[K, T]) Topic() string {
	return p.topic
}

func (p *SyncProducer[K, T]) Close() error {
	return p.producer.Close()
}

func (p *SyncProducer[K, T]) Send(k K, message T) error {
	key, value, headers, err := p.encode(k, message)
	if err != nil {
		return fmt.Errorf("encode: %w", err)
	}

	msg := &sarama.ProducerMessage{
		Topic:   p.topic,
		Value:   sarama.ByteEncoder(value),
		Headers: headers,
	}
	if len(key) > 0 {
		msg.Key = sarama.ByteEncoder(key)
	}

	_, _, err = p.producer.SendMessage(msg)

	return err
}

type failedEncoder struct {
	err error
}

func (e failedEncoder) Encode() ([]byte, error) {
	return nil, e.err
}

func (e failedEncoder) Length() int {
	return 0
}
