package kafka

import (
	"context"
	"log"
	"os"
	"os/signal"
	"payment-gateway/app/cmd"
	"payment-gateway/internal/kafka/handlers"
	"strings"
	"syscall"

	"github.com/Shopify/sarama"
)

var (
	SendTransactionKafkaTopic string
	consumerGroup             sarama.ConsumerGroup
)

func InitializeKafkaConsumer() {
	kafkaBrokerUrl := os.Getenv("KAFKA_BROKER_URL")
	if kafkaBrokerUrl == "" {
		log.Fatalf("KAFKA_BROKER_URL environment variable is not set")
	}

	SendTransactionKafkaTopic = os.Getenv("SEND_TRANSACTION_KAFKA_TOPIC")
	if SendTransactionKafkaTopic == "" {
		log.Fatalf("SEND_TRANSACTION_KAFKA_TOPIC environment variable is not set")
	}

	brokers := strings.Split(kafkaBrokerUrl, ",")

	kafkaGroupId := os.Getenv("KAFKA_GROUP_ID")
	if kafkaGroupId == "" {
		log.Fatalf("KAFKA_GROUP_ID environment variable is not set")
	}

	kafkaClientId := os.Getenv("KAFKA_CLIENT_ID")
	if kafkaClientId == "" {
		log.Fatalf("KAFKA_CLIENT_ID environment variable is not set")
	}

	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Version = sarama.V2_4_1_0
	config.ClientID = kafkaClientId
	config.Consumer.Offsets.AutoCommit.Enable = true
	config.Consumer.Offsets.Initial = sarama.OffsetNewest

	var err error
	consumerGroup, err = sarama.NewConsumerGroup(brokers, kafkaGroupId, config)
	if err != nil {
		log.Fatalf("Failed to start consumer group: %v", err)
	}

	topics := []string{SendTransactionKafkaTopic}
	startConsuming(topics)

	log.Printf("Kafka connected to brokers: %s, topic: %s\n", brokers, topics)
}

type ConsumerHandler struct{}

func (ConsumerHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	transactionHandler := handlers.NewTransactionHandler(
		*cmd.TransactionRepository,
		cmd.KafkaProducer,
		*cmd.SendTransactionClient,
		*cmd.GatewayCountryRepo,
		*cmd.GatewayRepo,
	)

	for message := range claim.Messages() {
		log.Printf("Message claimed: value = %s, topic = %s, partition = %d, offset = %d", string(message.Value), message.Topic, message.Partition, message.Offset)
		session.MarkMessage(message, "")

		switch message.Topic {
		case SendTransactionKafkaTopic:
			transactionHandler.HandleTransaction(context.Background(), message)
		}
	}
	return nil
}

func startConsuming(topics []string) {
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

func (ConsumerHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (ConsumerHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }
