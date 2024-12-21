package kafka

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Shopify/sarama"
)

var (
	consumerGroup             sarama.ConsumerGroup
	SendTransactionKafkaTopic string
)

func init() {
	kafkaURL := os.Getenv("KAFKA_BROKER_URL")
	if kafkaURL == "" {
		kafkaURL = "kafka:9092"
	}
	kafkaGroupID := os.Getenv("KAFKA_GROUP_ID")
	if kafkaGroupID == "" {
		kafkaGroupID = "payment-gateway-service"
	}
	SendTransactionKafkaTopic = os.Getenv("SEND_TRANSACTION_KAFKA_TOPIC")
	if SendTransactionKafkaTopic == "" {
		SendTransactionKafkaTopic = "process-transaction"
	}
	kafkaClientID := os.Getenv("KAFKA_CLIENT_ID")
	if kafkaClientID == "" {
		log.Fatalf("KAFKA_CLIENT_ID environment variable is not set")
	}

	config := sarama.NewConfig()
	config.ClientID = kafkaClientID
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRange
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	brokers := []string{kafkaURL}
	var err error
	consumerGroup, err = sarama.NewConsumerGroup(brokers, kafkaGroupID, config)
	if err != nil {
		log.Fatalf("Failed to initialize Kafka consumer group: %v", err)
	}

	log.Println("Kafka consumer group initialized successfully.")
}

type ConsumerHandler struct{}

func (ConsumerHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (ConsumerHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

func (ConsumerHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		log.Printf("Message claimed: value = %s, topic = %s, partition = %d, offset = %d", string(message.Value), message.Topic, message.Partition, message.Offset)
		session.MarkMessage(message, "")

		switch message.Topic {
		case SendTransactionKafkaTopic:
			// process
		}
	}
	return nil
}

func StartConsuming(topics []string) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	consumer := ConsumerHandler{}
	go func() {
		for {
			if err := consumerGroup.Consume(ctx, topics, consumer); err != nil {
				log.Fatalf("Error from consumer: %v", err)
			}
		}
	}()

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
	<-sigterm
	cancel()

	if err := consumerGroup.Close(); err != nil {
		log.Fatalf("Error closing consumer group: %v", err)
	}
}
