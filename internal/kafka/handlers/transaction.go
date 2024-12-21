package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"payment-gateway/internal/client"
	"payment-gateway/internal/config"
	"payment-gateway/internal/kafka"
	"payment-gateway/internal/repositories"
	"payment-gateway/models"
	"payment-gateway/pkg/constants"
	"payment-gateway/pkg/utils"

	"github.com/Shopify/sarama"
)

const (
	maxRetries = 3
)

// TransactionConsumer defines the interface for handling Kafka messages
type TransactionConsumer interface {
	Consume(ctx context.Context, message *sarama.ConsumerMessage) error
}

// TransactionHandler is the implementation of the TransactionConsumer
type TransactionHandler struct {
	transactionRepo       repositories.TransactionRepository
	kafkaProducer         kafka.KafkaProducer
	sendTransactionClient client.TransactionClient
	gatewayCountryRepo    repositories.GatewayCountryRepository
	gatewayRepo           repositories.GatewayRepository
}

// NewTransactionHandler initializes a new TransactionHandler
func NewTransactionHandler(
	transactionRepo repositories.TransactionRepository,
	kafkaProducer kafka.KafkaProducer,
	sendTransactionClient client.TransactionClient,
	gatewayCountryRepo repositories.GatewayCountryRepository,
	gatewayRepo repositories.GatewayRepository,
) *TransactionHandler {
	return &TransactionHandler{
		transactionRepo:       transactionRepo,
		kafkaProducer:         kafkaProducer,
		sendTransactionClient: sendTransactionClient,
		gatewayCountryRepo:    gatewayCountryRepo,
		gatewayRepo:           gatewayRepo,
	}
}

func (h *TransactionHandler) HandleTransaction(ctx context.Context, message *sarama.ConsumerMessage) {
	var transaction *models.Transaction
	if err := json.Unmarshal(message.Value, &transaction); err != nil {
		log.Printf("Failed to unmarshal message: %v", err)
		return
	}

	log.Printf("Processing transaction: %v", transaction)

	err := h.TransactionProcessor(ctx, transaction)
	if err != nil {
		if err.Error() == "error while SendTransaction" {
			log.Printf("Republish transactionID=%d to be retried, fallback to another gateway", transaction.ID)

			go h.kafkaProducer.ProduceMessage(message.Value, kafka.SendTransactionKafkaTopic)
			return
		}

		err = h.transactionRepo.UpdateTransactionStatusByReferenceID(ctx, transaction.ReferenceID.String(), constants.RETRY)
		if err != nil {
			log.Printf("Failed to UpdateTransactionStatusByReferenceID: %v", err)
		}

		return
	}

	log.Printf("Transaction successfully processed: %v", transaction)
}

func (h *TransactionHandler) TransactionProcessor(ctx context.Context, transaction *models.Transaction) error {
	gateway, err := h.gatewayCountryRepo.GetHealthyGatewayByCountryID(ctx, transaction.CountryID)
	if err != nil {
		log.Printf("Failed to GetHealthyGatewayByCountryID: %v", err)
		return err
	}

	gatewayConfig, err := config.GatewayConfigSelection(gateway.Name)
	if err != nil {
		log.Printf("Failed while GatewayConfigSelection: %v", err)
		return err
	}

	transactionRequest := models.SendTransactionRequest{
		ReferenceID: transaction.ReferenceID.String(),
		Amount:      transaction.Amount,
		UserID:      transaction.UserID,
		Currency:    transaction.Currency,
	}
	jsonData, err := json.Marshal(transactionRequest)
	if err != nil {
		return fmt.Errorf("failed to serialize request to JSON: %w", err)
	}

	encryptedPayload, err := utils.EncryptAES(string(jsonData), gatewayConfig.GatewayPrivateKey)
	if err != nil {
		log.Printf("Failed while EncryptAES: %v", err)
		return err
	}

	builtExternalTransaction, err := utils.BuildExternalTransactionRequest(gateway.DataFormatSupported, encryptedPayload)
	if err != nil {
		log.Printf("Failed while BuildExternalTransactionRequest: %v", err)
		return err
	}

	// retry until 3 times, if its still failed, fallback to another available gateway
	err = utils.RetryOperation(func() error {
		return h.sendTransactionClient.SendTransaction(ctx, builtExternalTransaction, gateway.Name, gatewayConfig)
	}, maxRetries)
	if err != nil {
		log.Printf("Failed while SendTransaction: %v", err)
		updateErr := h.gatewayRepo.UpdateHealthStatus(ctx, gateway.ID, constants.UNHEALTHY)
		if updateErr != nil {
			log.Printf("Failed to update health status for gateway ID %d: %v", gateway.ID, updateErr)
			return updateErr
		}
		return errors.New("error while SendTransaction")
	}

	err = h.transactionRepo.UpdateGatewayIDByTransactionID(ctx, transaction.ID, gateway.ID)
	if err != nil {
		log.Printf("Failed while UpdateGatewayIDByTransactionID: %v", err)
		return err
	}

	log.Printf("Transaction successfully processed: %v", transaction)
	return nil
}
