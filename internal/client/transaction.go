package client

import (
	"context"
	"fmt"
	"math/rand"
	"payment-gateway/models"
	"time"
)

type TransactionClient struct{}

func NewTransactionClient() *TransactionClient {
	return &TransactionClient{}
}

func (c *TransactionClient) SendTransaction(ctx context.Context, encryptedPayload string, gatewayConfig *models.GatewayConfig) (*models.GatewayResponse, error) {
	randGen := rand.New(rand.NewSource(time.Now().UnixNano()))
	isSuccess := randGen.Intn(2) == 0

	if !isSuccess {
		return nil, fmt.Errorf("gateway failed to process transaction")
	}

	response := &models.GatewayResponse{
		Status:  "success",
		Message: "Transaction processed successfully",
	}

	return response, nil
}
