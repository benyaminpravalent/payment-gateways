package mocks

import (
	"context"
	"payment-gateway/models"

	"github.com/stretchr/testify/mock"
)

// MockGatewayCountryRepository is a mock implementation of the GatewayCountryRepository
type MockGatewayCountryRepository struct {
	mock.Mock
}

// GetHealthyGatewayByCountryID provides a mock function for fetching a healthy gateway by country ID
func (m *MockGatewayCountryRepository) GetHealthyGatewayByCountryID(ctx context.Context, countryID int) (*models.GatewayDetail, error) {
	args := m.Called(ctx, countryID)

	var r0 *models.GatewayDetail
	if args.Get(0) != nil {
		r0 = args.Get(0).(*models.GatewayDetail)
	}
	r1 := args.Error(1)

	return r0, r1
}
