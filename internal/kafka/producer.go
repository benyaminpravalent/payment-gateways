package kafka

import (
	"fmt"
	"log"
	"os"

	"github.com/Shopify/sarama"
)

// KafkaProducer defines an interface for producing Kafka messages
// This allows for mocking in unit tests.
type KafkaProducer interface {
	ProduceMessage(data []byte, topic string) error
}

// SaramaProducer is an implementation of KafkaProducer using Sarama
type SaramaProducer struct {
	producer sarama.SyncProducer
}

// NewKafkaProducer initializes a new Kafka producer instance
func NewKafkaProducer() (KafkaProducer, error) {
	kafkaURL := os.Getenv("KAFKA_BROKER_URL")
	if kafkaURL == "" {
		kafkaURL = "kafka:9092"
	}

	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll // Wait for all in-sync replicas to acknowledge
	config.Producer.Retry.Max = 5                    // Retry up to 5 times
	config.Producer.Return.Successes = true          // Return successes

	brokers := []string{kafkaURL}
	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Kafka producer: %v", err)
	}

	log.Println("Kafka producer initialized successfully.")
	return &SaramaProducer{producer: producer}, nil
}

// ProduceMessage sends a message to the specified Kafka topic
func (p *SaramaProducer) ProduceMessage(data []byte, topic string) error {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(data),
	}

	partition, offset, err := p.producer.SendMessage(msg)
	if err != nil {
		return fmt.Errorf("failed to send message to Kafka: %v", err)
	}

	log.Printf("Message sent to topic %s, partition %d, offset %d\n", topic, partition, offset)
	return nil
}

// Close closes the underlying Kafka producer
func (p *SaramaProducer) Close() error {
	return p.producer.Close()
}
