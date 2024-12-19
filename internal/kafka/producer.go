package kafka

import (
	"fmt"
	"log"
	"os"

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
	kafkaURL := os.Getenv("KAFKA_BROKER_URL")
	if kafkaURL == "" {
		kafkaURL = "localhost:9092"
		log.Printf("KAFKA_BROKER_URL is not set. Using default: %s\n", kafkaURL)
	}

	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Errors = true
	config.Producer.Return.Successes = true

	brokers := []string{kafkaURL}
	log.Printf("Connecting to Kafka broker(s): %v\n", brokers)

	producer, err := sarama.NewSyncProducer(brokers, config)
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
