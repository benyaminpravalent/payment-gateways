package kafka

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/Shopify/sarama"
)

type KafkaProducer interface {
	ProduceMessage(data []byte, topic string) error
	Close() error
}

type SaramaProducer struct {
	producer sarama.SyncProducer
}

func NewKafkaProducer() KafkaProducer {
	kafkaBrokers := os.Getenv("KAFKA_BROKER_URL")
	if kafkaBrokers == "" {
		kafkaBrokers = "localhost:9092"
		log.Printf("KAFKA_BROKER_URL is not set. Using default: %s\n", kafkaBrokers)
	}

	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true

	log.Printf("Connecting to Kafka broker(s): %v\n", kafkaBrokers)

	producer, err := sarama.NewSyncProducer(strings.Split(kafkaBrokers, ","), config)
	if err != nil {
		log.Fatalf("Failed to initialize Kafka producer: %v", err)
	}

	log.Println("Kafka producer initialized successfully.")
	return &SaramaProducer{producer: producer}
}

func (p *SaramaProducer) ProduceMessage(data []byte, topic string) error {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(data),
	}

	log.Printf("Sending message to Kafka topic %s...\n", topic)
	partition, offset, err := p.producer.SendMessage(msg)
	if err != nil {
		return fmt.Errorf("failed to send message to Kafka topic %s: %w", topic, err)
	}

	log.Printf("Message sent successfully to topic %s, partition %d, offset %d\n", topic, partition, offset)
	return nil
}

func (p *SaramaProducer) Close() error {
	log.Println("Closing Kafka producer...")
	return p.producer.Close()
}
