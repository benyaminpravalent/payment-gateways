package services

import (
	"errors"
	"payment-gateway/internal/config"
	"payment-gateway/models"
	"payment-gateway/pkg/constants"
)

type GatewayService struct{}

func NewGatewayService() *GatewayService {
	return &GatewayService{}
}

func (s *GatewayService) GatewayConfigSelection(gatewayName string) (models.GatewayConfig, error) {
	switch gatewayName {
	case constants.GATEWAY_A:
		return models.GatewayConfig{
			GatewayUrl:        config.GatewayAUrl,
			GatewayApiKey:     config.GatewayAApiKey,
			GatewayPrivateKey: config.GatewayAPrivateKey,
		}, nil
	case constants.GATEWAY_B:
		return models.GatewayConfig{
			GatewayUrl:        config.GatewayBUrl,
			GatewayApiKey:     config.GatewayBApiKey,
			GatewayPrivateKey: config.GatewayBPrivateKey,
		}, nil
	case constants.GATEWAY_C:
		return models.GatewayConfig{
			GatewayUrl:        config.GatewayCUrl,
			GatewayApiKey:     config.GatewayCApiKey,
			GatewayPrivateKey: config.GatewayCPrivateKey,
		}, nil
	default:
		return models.GatewayConfig{}, errors.New("unsupported gateway name")
	}
}
