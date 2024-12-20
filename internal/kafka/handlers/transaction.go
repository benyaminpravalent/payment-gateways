package handlers

import (
	"context"
	"encoding/json"
	"log"

	"payment-gateway/internal/client"
	"payment-gateway/internal/kafka"
	"payment-gateway/internal/repositories"
	"payment-gateway/models"

	"github.com/Shopify/sarama"
)

// TransactionConsumer defines the interface for handling Kafka messages
type TransactionConsumer interface {
	Consume(ctx context.Context, message *sarama.ConsumerMessage) error
}

// TransactionHandler is the implementation of the TransactionConsumer
type TransactionHandler struct {
	transactionRepo        repositories.TransactionRepository
	kafkaProducer          kafka.KafkaProducer
	sendTransactionClient  client.TransactionClient
	gatewayCountryRepo     repositories.GatewayCountryRepository
}

// NewTransactionHandler initializes a new TransactionHandler
func NewTransactionHandler(
	transactionRepo repositories.TransactionRepository,
	kafkaProducer kafka.KafkaProducer,
	sendTransactionClient client.TransactionClient,
	gatewayCountryRepo repositories.GatewayCountryRepository,
) *TransactionHandler {
	return &TransactionHandler{
		transactionRepo:    transactionRepo,
		kafkaProducer:      kafkaProducer,
		sendTransactionClient: sendTransactionClient,
		gatewayCountryRepo: gatewayCountryRepo,
	}
}

// Consume processes a Kafka message
func (h *TransactionHandler) HandleTransaction(ctx context.Context, message *sarama.ConsumerMessage) error {
	log.Printf("Received message: Topic=%s Partition=%d Offset=%d", message.Topic, message.Partition, message.Offset)

	// Decode the message
	var transaction *models.Transaction
	if err := json.Unmarshal(message.Value, &transaction); err != nil {
		log.Printf("Failed to unmarshal message: %v", err)
		return err
	}

	log.Printf("Processing transaction: %v", transaction)

	// Process the transaction (e.g., interact with a third-party service)
	response, err := h.client.SendTransaction(ctx, transaction)
	if err != nil {
		log.Printf("Failed to send transaction to third party: %v", err)
		// Optionally, publish a failure message to another Kafka topic
		return err
	}

	log.Printf("Third-party response: %v", response)

	// Update transaction status in the database
	transaction.Status = "completed"
	if err := h.transactionRepo.UpdateTransaction(ctx, transaction); err != nil {
		log.Printf("Failed to update transaction in database: %v", err)
		return err
	}

	// Optionally, publish a success message to another Kafka topic
	if err := h.kafkaProducer.ProduceMessage(message.Value, "transaction-completed"); err != nil {
		log.Printf("Failed to publish success message: %v", err)
		return err
	}

	log.Printf("Transaction successfully processed: %v", transaction)
	return nil
}
