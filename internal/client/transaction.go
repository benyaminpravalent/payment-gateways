package client

import (
	"context"
	"fmt"
	"log"
	"payment-gateway/models"
	"time"

	"golang.org/x/exp/rand"
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
	randGen := rand.New(rand.NewSource(uint64(time.Now().UnixNano())))
	isSuccess := randGen.Intn(2) == 0

	/* ====to test the fault-tolerance==== */
	// isSuccess = false
	// if gatewayName == constants.GATEWAY_B {
	// 	isSuccess = true
	// }

	if !isSuccess {
		return fmt.Errorf("gateway failed to process transaction")
	}

	log.Printf("Transaction is processed by gateway=[%s]", gatewayName)

	return nil
}
