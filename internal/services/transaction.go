package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"payment-gateway/internal/kafka"
	"payment-gateway/internal/repositories"
	"payment-gateway/models"
	"payment-gateway/pkg/constants"
	"time"

	"github.com/google/uuid"
)

type TransactionService struct {
	TransactionRepository *repositories.TransactionRepository
	KafkaProducer         kafka.KafkaProducer
}

func NewTransactionService(
	transactionRepository *repositories.TransactionRepository,
	kafkaProducer kafka.KafkaProducer,
) *TransactionService {
	return &TransactionService{
		TransactionRepository: transactionRepository,
		KafkaProducer:         kafkaProducer,
	}
}

func (s *TransactionService) Deposit(ctx context.Context, request models.DepositRequest) (models.Transaction, error) {
	transaction := &models.Transaction{
		ReferenceID: uuid.New(),
		Amount:      request.Amount,
		Currency:    request.Currency,
		Type:        constants.DEPOSIT,
		Status:      constants.PENDING,
		CountryID:   request.CountryID,
		UserID:      request.UserID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err := s.TransactionRepository.InsertTransaction(ctx, transaction)
	if err != nil {
		return models.Transaction{}, fmt.Errorf("[service-Deposit] Error while InsertTransaction = %v", err)
	}

	messageBytes, err := json.Marshal(transaction)
	if err != nil {
		log.Printf("Failed to marshal Kafka message: %v", err)
		return models.Transaction{}, fmt.Errorf("[service-Deposit] Error while Marshal = %v", err)
	}

	go func(messageBytes []byte) {

		s.KafkaProducer.ProduceMessage(messageBytes, kafka.SendTransactionKafkaTopic)
	}(messageBytes)

	return *transaction, nil
}

func (s *TransactionService) Withdraw(ctx context.Context, request models.WithdrawalRequest) (models.Transaction, error) {
	transaction := &models.Transaction{
		ReferenceID: uuid.New(),
		Amount:      request.Amount,
		Currency:    request.Currency,
		Type:        constants.WITHDRAWAL,
		Status:      constants.PENDING,
		CountryID:   request.CountryID,
		UserID:      request.UserID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err := s.TransactionRepository.InsertTransaction(ctx, transaction)
	if err != nil {
		return models.Transaction{}, fmt.Errorf("[service-Withdraw] Error while InsertTransaction = %v", err)
	}

	messageBytes, err := json.Marshal(transaction)
	if err != nil {
		log.Printf("Failed to marshal Kafka message: %v", err)
		return models.Transaction{}, fmt.Errorf("[service-Withdraw] Error while Marshal = %v", err)
	}

	go func(messageBytes []byte) {

		s.KafkaProducer.ProduceMessage(messageBytes, kafka.SendTransactionKafkaTopic)
	}(messageBytes)

	return *transaction, nil
}

func (s *TransactionService) TransactionCallback(ctx context.Context, request *models.TransactionCallbackRequest) error {
	return s.TransactionRepository.UpdateTransactionStatusByReferenceID(ctx, request.ReferenceID, request.Status)
}
