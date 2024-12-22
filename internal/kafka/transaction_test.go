package kafka

import (
	"context"
	"encoding/json"
	mocksClient "payment-gateway/mocks/client"
	mockKafka "payment-gateway/mocks/kafka"
	mocksRepository "payment-gateway/mocks/repositories"
	"payment-gateway/models"
	"payment-gateway/pkg/constants"
	"testing"

	"github.com/Shopify/sarama"
	"github.com/google/uuid"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

func TestTransactionService(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "TransactionHandler Suite")
}

var _ = ginkgo.Describe("TransactionHandler", func() {
	var (
		mockTransactionRepo       *mocksRepository.TransactionRepository
		mockKafkaProducer         *mockKafka.MockKafkaProducer
		mockSendTransactionClient *mocksClient.MockTransactionClient
		mockGatewayCountryRepo    *mocksRepository.MockGatewayCountryRepository
		mockGatewayRepo           *mocksRepository.MockGatewayRepository
		transactionHandler        *TransactionHandler
	)

	ginkgo.BeforeEach(func() {
		mockTransactionRepo = new(mocksRepository.TransactionRepository)
		mockKafkaProducer = new(mockKafka.MockKafkaProducer)
		mockSendTransactionClient = new(mocksClient.MockTransactionClient)
		mockGatewayCountryRepo = new(mocksRepository.MockGatewayCountryRepository)
		mockGatewayRepo = new(mocksRepository.MockGatewayRepository)
		transactionHandler = NewTransactionHandler(
			mockTransactionRepo,
			mockKafkaProducer,
			mockSendTransactionClient,
			mockGatewayCountryRepo,
			mockGatewayRepo,
		)
	})

	ginkgo.Describe("HandleTransaction", func() {
		var (
			mockCtx      context.Context
			mockMessage  *sarama.ConsumerMessage
			transaction  *models.Transaction
			messageBytes []byte
		)

		ginkgo.BeforeEach(func() {
			mockCtx = context.Background()
			transaction = &models.Transaction{
				ID:          12345,
				ReferenceID: uuid.New(),
				Amount:      1000,
				Currency:    "USD",
				Type:        constants.DEPOSIT,
				GatewayID:   1,
				CountryID:   1,
				Status:      constants.PENDING,
				UserID:      1,
			}

			messageBytes, _ = json.Marshal(transaction)
			mockMessage = &sarama.ConsumerMessage{
				Value: messageBytes,
			}
		})

		ginkgo.It("should handle unmarshalling error", func() {
			mockMessage.Value = []byte("invalid json")
			err := transactionHandler.HandleTransaction(mockCtx, mockMessage)
			gomega.Expect(err).Should(gomega.HaveOccurred())
		})
	})
})
