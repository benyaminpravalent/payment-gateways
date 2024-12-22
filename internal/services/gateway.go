package services

import (
	"context"
	"payment-gateway/internal/repositories"
)

type GatewayService struct {
	gatewayRepository repositories.IGatewayRepository
}

func NewGatewayService(
	gatewayRepository repositories.IGatewayRepository,
) *GatewayService {
	return &GatewayService{
		gatewayRepository: gatewayRepository,
	}
}

func (g *GatewayService) UpdateGatewayHealthStatusByID(ctx context.Context, gatewayID int, status string) error {
	return g.gatewayRepository.UpdateHealthStatus(ctx, gatewayID, status)
}
