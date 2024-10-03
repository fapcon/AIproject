package kafka

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/IBM/sarama"
	"github.com/google/uuid"
	"go.uber.org/multierr"
	"studentgit.kata.academy/quant/torque/pkg/logster"
	"studentgit.kata.academy/quant/torque/pkg/scram"
)

type InitialOffset string

const (
	Newest InitialOffset = "newest"
	Oldest InitialOffset = "oldest"
)

func (o InitialOffset) toSarama() int64 {
	switch o {
	case Newest:
		return sarama.OffsetNewest
	case Oldest:
		return sarama.OffsetOldest
	default:
		return sarama.OffsetNewest
	}
}

type ConsumerConfig struct {
	ClientID      string        `yaml:"client_id"`
	RawGroupID    string        `yaml:"group_id"`
	UniqueGroup   bool          `yaml:"unique_group"`
	InitialOffset InitialOffset `yaml:"initial_offset"`
	Topics        []string      `yaml:"topics"`
	Brokers       []string      `yaml:"brokers"`

	EnableAuth bool   `yaml:"enable_auth"`
	User       string `yaml:"user"`
	Password   string `yaml:"password"`

	once    sync.Once `yaml:"-"`
	groupID string    `yaml:"-"`
}

func (c *ConsumerConfig) Topic() string {
	if len(c.Topics) == 0 {
		return ""
	}
	if len(c.Topics) == 1 {
		return c.Topics[0]
	}
	return strings.Join(c.Topics, ", ")
}

func (c *ConsumerConfig) GroupID() string {
	c.once.Do(func() {
		c.groupID = c.RawGroupID
		if c.UniqueGroup {
			c.groupID += "-" + uuid.NewString()
		}
	})
	return c.groupID
}

func (c *ConsumerConfig) String() string {
	return fmt.Sprintf(
		`{Brokers:%v Topics:%v GroupID:%s ClientID:%s EnableAuth: %t}`,
		c.Brokers, c.Topics, c.GroupID(), c.ClientID, c.EnableAuth)
}

func (c *ConsumerConfig) toSarama() *sarama.Config {
	config := sarama.NewConfig()

	config.Version = sarama.V2_2_0_0
	config.ClientID = c.ClientID

	config.Metadata.Full = false

	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.Initial = c.InitialOffset.toSarama()

	if c.EnableAuth {
		config.Net.SASL.Mechanism = sarama.SASLTypeSCRAMSHA512
		config.Net.SASL.Enable = true
		config.Net.SASL.User = c.User
		config.Net.SASL.Password = c.Password
		config.Net.SASL.SCRAMClientGeneratorFunc = scram.NewSHA512Client
	}

	return config
}

type Consumer[T any] struct {
	logger   logster.Logger
	config   *ConsumerConfig
	handler  sarama.ConsumerGroupHandler
	messages chan Ackable[T]
	client   sarama.Client
	group    sarama.ConsumerGroup
	waiter   *WaitingHandler
}

func NewConsumer[T any](
	logger logster.Logger,
	config *ConsumerConfig,
	decode func(context.Context, *sarama.ConsumerMessage) (T, Mark, error),
) (*Consumer[T], error) {
	logger = logger.WithFields(logster.Fields{
		"message_type": messageType[T](),
		"topics":       config.Topics,
		"unique_group": config.UniqueGroup,
	})
	messages := make(chan Ackable[T])
	rawHandler := NewHandler(logger, messages, decode)
	waiter := NewWaitingHandler(logger, rawHandler, config.GroupID())
	handler := waiter
	return &Consumer[T]{
		logger:   logger,
		config:   config,
		handler:  handler,
		messages: messages,
		client:   nil, // will be set in Init
		group:    nil, // will be set in Init
		waiter:   waiter,
	}, nil
}

func (c *Consumer[T]) Topic() string {
	return c.config.Topic()
}

// Init creates consumer group and fixes current offset, but doesn't process messages.
// Need to call [Consumer.Run] to start processing.
func (c *Consumer[T]) Init(context.Context) (err error) {
	config := c.config.toSarama()

	client, err := sarama.NewClient(c.config.Brokers, config)
	if err != nil {
		return fmt.Errorf("sarama.NewClient: %w", err)
	}
	defer func() {
		if err != nil {
			err = multierr.Append(err, client.Close())
		}
	}()

	group, err := sarama.NewConsumerGroupFromClient(c.config.GroupID(), client)
	if err != nil {
		return fmt.Errorf("sarama.NewConsumerGroupFromClient: %w", err)
	}
	defer func() {
		if err != nil {
			err = multierr.Append(err, group.Close())
		}
	}()

	go c.logErrors(group.Errors())

	c.waiter.SetClient(client)
	c.client = client
	c.group = group

	// If initial offset is newest, we need to pin (to fix) it now.
	// So consumption after calling [Consumer.Run] will start from that pinned offset.
	// This guarantees we don't miss messages between [Consumer.Init] and [Consumer.Run].
	//
	// We look up through every partition to find and mark an unmarked offset.
	if config.Consumer.Offsets.Initial == sarama.OffsetNewest {
		return c.pinOffsets()
	}

	c.logger.
		WithField("sarama_config", config).
		WithField("config", c.config).
		Debugf("Consumer: initialized")

	return nil
}

func (c *Consumer[T]) Run(ctx context.Context) error {
	logger := c.logger.WithContext(ctx)

	defer close(c.messages)

	for {
		err := c.group.Consume(ctx, c.config.Topics, c.handler)
		switch {
		case errors.Is(err, sarama.ErrClosedConsumerGroup) || errors.Is(err, sarama.ErrClosedClient):
			return fmt.Errorf("consumption stopped unexpectedly: %w", err)
		case ctx.Err() != nil:
			logger.Infof("Consumer.Run(): shutdown")
			return nil //nolint:nilerr // it's ok to ignore error if context is done
		case err != nil:
			logger.WithError(err).Errorf("Consumer.Run(): consumption error")
		}
		logger.Infof("Consumer.Run(): consumption stopped, restarting")
	}
}

func (c *Consumer[T]) Close() error {
	return multierr.Combine(c.group.Close(), c.client.Close())
}

func (c *Consumer[T]) Messages() <-chan Ackable[T] {
	return c.messages
}

func (c *Consumer[T]) WaitNewestAcked(ctx context.Context) error {
	return c.waiter.WaitNewestMarked(ctx)
}

func (c *Consumer[T]) logErrors(errs <-chan error) {
	for err := range errs {
		logger := c.logger.WithError(err)
		var consumerError *sarama.ConsumerError
		if errors.As(err, &consumerError) {
			logger = logger.WithField("topic", consumerError.Topic)
		}
		logger.Errorf("Consumer error")
	}
}

func (c *Consumer[T]) pinOffsets() (err error) {
	coordinator, err := c.client.Coordinator(c.config.GroupID())
	if err != nil {
		return fmt.Errorf("client.Coordinator: %w", err)
	}
	defer func() {
		if err != nil {
			err = multierr.Append(err, coordinator.Close())
		}
	}()

	notCommitted, ok, err := c.getNotCommittedOffsets(coordinator)
	if err != nil {
		return fmt.Errorf("getNotCommittedOffsets: %w", err)
	}

	if !ok {
		c.logger.Debugf("Consumer: all offsets already committed")
		return nil
	}

	newestOffsets := make(map[string]map[int32]int64, len(notCommitted))
	for topic, partitions := range notCommitted {
		m := make(map[int32]int64, len(partitions))
		for partition := range partitions {
			offset, errX := c.client.GetOffset(topic, partition, sarama.OffsetNewest)
			if errX != nil {
				return fmt.Errorf("client.GetOffset: %w", errX)
			}

			m[partition] = offset
		}
		newestOffsets[topic] = m
	}

	manager, err := sarama.NewOffsetManagerFromClient(c.config.GroupID(), c.client)
	if err != nil {
		return fmt.Errorf("sarama.NewOffsetManagerFromClient: %w", err)
	}

	defer manager.Commit()

	for topic, partitions := range notCommitted {
		for partition := range partitions {
			partitionManager, errX := manager.ManagePartition(topic, partition)
			if errX != nil {
				return fmt.Errorf("manager.ManagePartition(%s/%d): %w", topic, partition, errX)
			}
			partitionManager.MarkOffset(newestOffsets[topic][partition], "")
		}
	}

	c.logger.WithField("offsets", newestOffsets).Debugf("Consumer: offsets committed")

	return nil
}

func (c *Consumer[T]) getNotCommittedOffsets(coordinator *sarama.Broker) (map[string]map[int32]struct{}, bool, error) {
	allPartitions := make(map[string]map[int32]struct{}, len(c.config.Topics))
	for _, topic := range c.config.Topics {
		partitions, err := c.client.Partitions(topic)
		if err != nil {
			return nil, false, fmt.Errorf("client.Partitions(%s): %w", topic, err)
		}
		x := make(map[int32]struct{}, len(partitions))
		for _, partition := range partitions {
			x[partition] = struct{}{}
		}
		allPartitions[topic] = x
	}

	committedOffsets, err := getCommittedOffsets(coordinator, c.config.GroupID(), allPartitions)
	if err != nil {
		return nil, false, fmt.Errorf("getCommittedOffsets: %w", err)
	}

	c.logger.WithField("offsets", committedOffsets).Debugf("Consumer: committed offsets")

	ok := false
	notCommitted := make(map[string]map[int32]struct{}, len(c.config.Topics))
	for topic, partitions := range committedOffsets {
		x := make(map[int32]struct{}, len(partitions))
		for partition, offset := range partitions {
			if offset < 0 {
				// Negative value means group haven't committed offset yet.
				// In this case need to mark (and commit) current newest offset.
				// This allows to pin (to fix) current newest offset to be consumed from.
				// Otherwise, real handler will receive the newest offset at the time of consumption start.
				ok = true
				x[partition] = struct{}{}
			}
		}
		notCommitted[topic] = x
	}
	return notCommitted, ok, nil
}
