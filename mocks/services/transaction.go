package mocks

import (
	"context"
	"payment-gateway/models"

	"github.com/stretchr/testify/mock"
)

type TransactionService struct {
	mock.Mock
}

func (m *TransactionService) Deposit(ctx context.Context, request models.DepositRequest) (models.Transaction, error) {
	args := m.Called(ctx, request)
	return args.Get(0).(models.Transaction), args.Error(1)
}

func (m *TransactionService) Withdraw(ctx context.Context, request models.WithdrawalRequest) (models.Transaction, error) {
	args := m.Called(ctx, request)
	return args.Get(0).(models.Transaction), args.Error(1)
}

func (m *TransactionService) TransactionCallback(ctx context.Context, request *models.TransactionCallbackRequest) error {
	args := m.Called(ctx, request)
	return args.Error(0)
}
