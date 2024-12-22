package mocks

import (
	"context"
	"payment-gateway/models"

	"github.com/stretchr/testify/mock"
)

// MockTransactionClient is a mock implementation of the TransactionClient
type MockTransactionClient struct {
	mock.Mock
}

// SendTransaction provides a mock function for sending a transaction
func (m *MockTransactionClient) SendTransaction(
	ctx context.Context,
	transactionRequest models.BuildExternalTransaction,
	gatewayName string,
	gatewayConfig models.GatewayConfig,
) error {
	args := m.Called(ctx, transactionRequest, gatewayName, gatewayConfig)
	return args.Error(0)
}
