package config

import (
	"errors"
	"payment-gateway/models"
	"payment-gateway/pkg/constants"
)

func GatewayConfigSelection(gatewayName string) (models.GatewayConfig, error) {
	switch gatewayName {
	case constants.GATEWAY_A:
		return models.GatewayConfig{
			GatewayUrl:        GatewayAUrl,
			GatewayApiKey:     GatewayAApiKey,
			GatewayPrivateKey: GatewayAPrivateKey,
		}, nil
	case constants.GATEWAY_B:
		return models.GatewayConfig{
			GatewayUrl:        GatewayBUrl,
			GatewayApiKey:     GatewayBApiKey,
			GatewayPrivateKey: GatewayBPrivateKey,
		}, nil
	case constants.GATEWAY_C:
		return models.GatewayConfig{
			GatewayUrl:        GatewayCUrl,
			GatewayApiKey:     GatewayCApiKey,
			GatewayPrivateKey: GatewayCPrivateKey,
		}, nil
	default:
		return models.GatewayConfig{}, errors.New("unsupported gateway name")
	}
}
