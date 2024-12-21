package client

import (
	"context"
	"fmt"
	"log"
	"payment-gateway/models"
)

type TransactionClient struct{}

func NewTransactionClient() *TransactionClient {
	return &TransactionClient{}
}

func (c *TransactionClient) SendTransaction(
	ctx context.Context,
	transactionRequest models.BuildExternalTransaction,
	gatewayName string,
	gatewayConfig models.GatewayConfig,
) error {
	isSuccess := true
	// uncomment to test the fault-tolerance
	// randGen := rand.New(rand.NewSource(time.Now().UnixNano()))
	// isSuccess = randGen.Intn(2) == 0

	if !isSuccess {
		return fmt.Errorf("gateway failed to process transaction")
	}

	log.Printf("Transaction processed successfully by gateway=[%s]", gatewayName)

	return nil
}
