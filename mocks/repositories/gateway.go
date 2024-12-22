package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
)

// MockGatewayRepository is a mock implementation of the GatewayRepository
type MockGatewayRepository struct {
	mock.Mock
}

// UpdateHealthStatus provides a mock function for updating the health status of a gateway
func (m *MockGatewayRepository) UpdateHealthStatus(ctx context.Context, gatewayID int, healthStatus string) error {
	args := m.Called(ctx, gatewayID, healthStatus)
	return args.Error(0)
}
