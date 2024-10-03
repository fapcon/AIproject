package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/IBM/sarama"
)

func main() {
	brokers := []string{"localhost:9092"}
	topic := "torque.events.private"
	group := "api-local.torque.events.private"

	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Version = sarama.V2_0_0_0                      // Adjust the Kafka version accordingly.
	config.Consumer.Offsets.Initial = sarama.OffsetOldest // Start consuming from the earliest message

	// Create new consumer group
	consumerGroup, err := sarama.NewConsumerGroup(brokers, group, config)
	if err != nil {
		log.Fatalf("Error creating consumer group: %v", err)
	}
	defer consumerGroup.Close()

	consume(consumerGroup, topic)
}

func consume(consumerGroup sarama.ConsumerGroup, topic string) {
	// Track errors
	go func() {
		for err := range consumerGroup.Errors() {
			fmt.Printf("Error: %s\n", err.Error())
		}
	}()

	// Handle session lifecycle and consume messages
	consumer := Consumer{
		ready: make(chan bool),
	}
	ctx := context.Background()
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			// 'Consume' should be called inside an infinite loop, it manages rebalance and returns errors
			if err := consumerGroup.Consume(ctx, []string{topic}, &consumer); err != nil {
				log.Fatalf("Error from consumer: %v", err)
			}
			// Check if context was cancelled, signaling that the consumer should stop
			if ctx.Err() != nil {
				return
			}
			consumer.ready = make(chan bool)
		}
	}()

	<-consumer.ready // Await till the consumer has been set up
	log.Println("Sarama consumer up and running!...")

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
	<-sigterm // Await a signal before exiting
}

// Consumer represents a Sarama consumer group consumer
type Consumer struct {
	ready chan bool
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (consumer *Consumer) Setup(sarama.ConsumerGroupSession) error {
	// Mark the consumer as ready
	close(consumer.ready)
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (consumer *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (consumer *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		fmt.Printf("Message claimed: value = %s, timestamp = %v, topic = %s\n", string(message.Value), message.Timestamp, message.Topic)
		session.MarkMessage(message, "") // Mark message as processed
	}

	return nil
}
