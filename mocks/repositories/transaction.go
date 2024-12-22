package mocks

import (
	"context"
	"payment-gateway/models"

	"github.com/stretchr/testify/mock"
)

type TransactionRepository struct {
	mock.Mock
}

func (m *TransactionRepository) InsertTransaction(ctx context.Context, transaction *models.Transaction) error {
	args := m.Called(ctx, transaction)
	return args.Error(0)
}

func (m *TransactionRepository) UpdateTransactionStatusByReferenceID(ctx context.Context, referenceID string, status string) error {
	args := m.Called(ctx, referenceID, status)
	return args.Error(0)
}

func (m *TransactionRepository) UpdateGatewayIDByTransactionID(ctx context.Context, transactionID int, gatewayID int) error {
	args := m.Called(ctx, transactionID, gatewayID)
	return args.Error(0)
}
