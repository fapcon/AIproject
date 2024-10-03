package main

import (
	"context"
	"fmt"
	"github.com/IBM/sarama"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"studentgit.kata.academy/quant/torque/internal/ports/types/api"
)

func main() {
	// Set up gRPC client
	conn, err := grpc.Dial("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client := api.NewAPIServiceClient(conn)

	config := sarama.NewConfig()
	config.Version = sarama.V2_5_0_0
	group, err := sarama.NewConsumerGroup([]string{"localhost:9092"}, "api-local.torque.events.private", config)
	if err != nil {
		log.Fatalf("Error creating consumer group: %v", err)
	}

	// Consumer logic; define a ConsumerGroupHandler
	consumer := Consumer{
		ready: make(chan bool),
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Consume messages
	go func() {
		for {
			if err := group.Consume(ctx, []string{"torque.events.private"}, &consumer); err != nil {
				log.Fatalf("Error consuming from topic: %v", err)
			}
			if ctx.Err() != nil {
				return
			}
			consumer.ready = make(chan bool)
		}
	}()

	<-consumer.ready // Await till the consumer has been set up
	log.Println("Sarama consumer up and running!...")

	_, err = client.SubscribeInstruments(context.Background(), &api.Instruments{Instruments: []*api.Instrument{{
		Exchange:       1,
		LocalSymbol:    "ETH/USDT",
		ExchangeSymbol: "ETH-USDT",
	}}})
	if err != nil {
		fmt.Println(err)
		return
	}
}

type Consumer struct {
	ready chan bool
}

func (consumer *Consumer) Setup(sarama.ConsumerGroupSession) error {
	close(consumer.ready)
	return nil
}

func (consumer *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (consumer *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		log.Printf("Message claimed: value = %s, timestamp = %v, topic = %s", string(message.Value), message.Timestamp, message.Topic)
		// Implement your verification logic here
		session.MarkMessage(message, "")
	}
	return nil
}
