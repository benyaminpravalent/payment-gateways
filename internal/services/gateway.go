package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"payment-gateway/internal/config"
	"payment-gateway/models"
	"payment-gateway/pkg/constants"
	"payment-gateway/pkg/utils"
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

func (s *GatewayService) EncryptTransactionRequest(request *models.TransactionRequest, privateKey string) (string, error) {
	jsonData, err := json.Marshal(request)
	if err != nil {
		return "", fmt.Errorf("failed to serialize request to JSON: %w", err)
	}

	encryptedData, err := utils.EncryptAES(string(jsonData), privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to encrypt request: %w", err)
	}

	return encryptedData, nil
}

func (s *GatewayService) DecryptGatewayResponse(encryptedResponse, privateKey string) (*models.GatewayCallback, error) {
	decryptedData, err := utils.DecryptAES(encryptedResponse, privateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt gateway response: %w", err)
	}

	var callback models.GatewayCallback
	err = json.Unmarshal([]byte(decryptedData), &callback)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal gateway response: %w", err)
	}

	return &callback, nil
}
