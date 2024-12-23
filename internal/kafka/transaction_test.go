package kafka

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"payment-gateway/internal/config"
	mocksClient "payment-gateway/mocks/client"
	mockKafka "payment-gateway/mocks/kafka"
	mocksRepository "payment-gateway/mocks/repositories"
	"payment-gateway/models"
	"payment-gateway/pkg/constants"
	"sync"
	"testing"

	"github.com/Shopify/sarama"
	"github.com/google/uuid"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
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
		os.Setenv("GATEWAY_A_URL", "A")
		os.Setenv("GATEWAY_A_API_KEY", "api-key")
		os.Setenv("GATEWAY_A_PRIVATE_KEY", "12345678901234567890123456789012")
		config.InitGatewayA()

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

	ginkgo.AfterEach(func() {
		os.Unsetenv("GATEWAY_A_URL")
		os.Unsetenv("GATEWAY_A_API_KEY")
		os.Unsetenv("GATEWAY_A_PRIVATE_KEY")
	})

	ginkgo.Describe("HandleTransaction", func() {
		var (
			mockCtx      context.Context
			mockMessage  *sarama.ConsumerMessage
			transaction  *models.Transaction
			messageBytes []byte
		)

		gateway := &models.GatewayDetail{
			ID:                  1,
			Name:                "A",
			DataFormatSupported: "json",
			HealthStatus:        "healthy",
			Priority:            1,
			CountryID:           1,
			Currency:            "USD",
		}
		gatewayConfig := models.GatewayConfig{
			GatewayUrl:        "A",
			GatewayApiKey:     "api-key",
			GatewayPrivateKey: "12345678901234567890123456789012",
		}

		ginkgo.BeforeEach(func() {
			mockCtx = context.Background()
			transaction = &models.Transaction{
				ID:          12345,
				ReferenceID: uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
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

		ginkgo.It("should handle error from TransactionProcessor", func() {
			mockGatewayCountryRepo.
				On("GetHealthyGatewayByCountryID", mock.Anything, transaction.CountryID).
				Return(nil, errors.New("error")).
				Once()
			mockTransactionRepo.
				On("UpdateTransactionStatusByReferenceID", mock.Anything, transaction.ReferenceID.String(), constants.RETRY).
				Return(nil).
				Once()

			err := transactionHandler.HandleTransaction(mockCtx, mockMessage)
			gomega.Expect(err).Should(gomega.HaveOccurred())

			mockGatewayCountryRepo.AssertCalled(ginkgo.GinkgoT(), "GetHealthyGatewayByCountryID", mock.Anything, transaction.CountryID)
			mockTransactionRepo.AssertCalled(ginkgo.GinkgoT(), "UpdateTransactionStatusByReferenceID", mock.Anything, mock.AnythingOfType("string"), constants.RETRY)
		})

		ginkgo.It("should handle error while SendTransaction from TransactionProcessor", func() {
			mockGatewayCountryRepo.
				On("GetHealthyGatewayByCountryID", mock.Anything, transaction.CountryID).
				Return(nil, errors.New("error while SendTransaction")).
				Once()

			var wg sync.WaitGroup
			wg.Add(1)

			mockKafkaProducer.On("ProduceMessage", mockMessage.Value, SendTransactionKafkaTopic).Run(func(args mock.Arguments) {
				wg.Done()
			}).Return(nil)

			err := transactionHandler.HandleTransaction(mockCtx, mockMessage)
			gomega.Expect(err).Should(gomega.HaveOccurred())

			wg.Wait()

			mockGatewayCountryRepo.AssertCalled(ginkgo.GinkgoT(), "GetHealthyGatewayByCountryID", mock.Anything, transaction.CountryID)
			mockKafkaProducer.AssertCalled(ginkgo.GinkgoT(), "ProduceMessage", mockMessage.Value, SendTransactionKafkaTopic)
		})

		ginkgo.It("should successfully process the transaction", func() {
			mockGatewayCountryRepo.
				On("GetHealthyGatewayByCountryID", mock.Anything, transaction.CountryID).
				Return(gateway, nil).
				Once()

			mockSendTransactionClient.
				On("SendTransaction", mockCtx, mock.Anything, gateway.Name, gatewayConfig).
				Return(nil).
				Once()

			mockTransactionRepo.
				On("UpdateGatewayIDByTransactionID", mock.Anything, transaction.ID, gateway.ID).
				Return(nil).
				Once()

			err := transactionHandler.HandleTransaction(mockCtx, mockMessage)
			gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
		})
	})

	ginkgo.Describe("TransactionProcessor", func() {
		var (
			mockCtx     context.Context
			transaction *models.Transaction
		)

		gateway := &models.GatewayDetail{
			ID:                  1,
			Name:                "A",
			DataFormatSupported: "json",
			HealthStatus:        "healthy",
			Priority:            1,
			CountryID:           1,
			Currency:            "USD",
		}
		gatewayConfig := models.GatewayConfig{
			GatewayUrl:        "A",
			GatewayApiKey:     "api-key",
			GatewayPrivateKey: "12345678901234567890123456789012",
		}

		ginkgo.BeforeEach(func() {
			mockCtx = context.Background()
			transaction = &models.Transaction{
				ID:          12345,
				ReferenceID: uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
				Amount:      1000,
				Currency:    "USD",
				Type:        constants.DEPOSIT,
				GatewayID:   1,
				CountryID:   1,
				Status:      constants.PENDING,
				UserID:      1,
			}
		})

		ginkgo.It("should handle when there's no healthy gateway", func() {
			mockGatewayCountryRepo.
				On("GetHealthyGatewayByCountryID", mock.Anything, 1).
				Return(nil, errors.New("error")).
				Once()

			err := transactionHandler.TransactionProcessor(mockCtx, transaction)
			gomega.Expect(err).Should(gomega.HaveOccurred())
		})

		ginkgo.It("should handle when the gateway is not configured", func() {
			gateway := &models.GatewayDetail{
				ID:                  1,
				Name:                "B",
				DataFormatSupported: "json",
				HealthStatus:        "healthy",
				Priority:            1,
				CountryID:           1,
				Currency:            "USD",
			}

			mockGatewayCountryRepo.
				On("GetHealthyGatewayByCountryID", mock.Anything, transaction.CountryID).
				Return(gateway, nil).
				Once()

			err := transactionHandler.TransactionProcessor(mockCtx, transaction)
			gomega.Expect(err).Should(gomega.HaveOccurred())
		})

		ginkgo.It("should handle error when ecrypting the request", func() {
			gateway := &models.GatewayDetail{
				ID:                  1,
				Name:                "B",
				DataFormatSupported: "json",
				HealthStatus:        "healthy",
				Priority:            1,
				CountryID:           1,
				Currency:            "USD",
			}

			mockGatewayCountryRepo.
				On("GetHealthyGatewayByCountryID", mock.Anything, transaction.CountryID).
				Return(gateway, nil).
				Once()

			err := transactionHandler.TransactionProcessor(mockCtx, transaction)
			gomega.Expect(err).Should(gomega.HaveOccurred())
		})

		ginkgo.It("should handle error when build external request", func() {
			gateway := &models.GatewayDetail{
				ID:                  1,
				Name:                "A",
				DataFormatSupported: "unsupported-data-format",
				HealthStatus:        "healthy",
				Priority:            1,
				CountryID:           1,
				Currency:            "USD",
			}

			mockGatewayCountryRepo.
				On("GetHealthyGatewayByCountryID", mock.Anything, transaction.CountryID).
				Return(gateway, nil).
				Once()

			err := transactionHandler.TransactionProcessor(mockCtx, transaction)
			gomega.Expect(err).Should(gomega.HaveOccurred())
		})

		ginkgo.It("should handle error when build external request", func() {
			mockGatewayCountryRepo.
				On("GetHealthyGatewayByCountryID", mock.Anything, transaction.CountryID).
				Return(gateway, nil).
				Once()

			mockSendTransactionClient.
				On("SendTransaction", mockCtx, mock.AnythingOfType("models.BuildExternalTransaction"), gateway.Name, gatewayConfig).
				Return(errors.New("error")).
				Times(3)

			mockGatewayRepo.
				On("UpdateHealthStatus", mockCtx, gateway.ID, constants.UNHEALTHY).
				Return(nil).
				Once()

			err := transactionHandler.TransactionProcessor(mockCtx, transaction)
			gomega.Expect(err).Should(gomega.HaveOccurred())
		})

		ginkgo.It("should handle error when update gateway health status", func() {
			mockGatewayCountryRepo.
				On("GetHealthyGatewayByCountryID", mock.Anything, transaction.CountryID).
				Return(gateway, nil).
				Once()

			mockSendTransactionClient.
				On("SendTransaction", mockCtx, mock.AnythingOfType("models.BuildExternalTransaction"), gateway.Name, gatewayConfig).
				Return(errors.New("error")).
				Times(3)

			mockGatewayRepo.
				On("UpdateHealthStatus", mockCtx, gateway.ID, constants.UNHEALTHY).
				Return(errors.New("error")).
				Once()

			err := transactionHandler.TransactionProcessor(mockCtx, transaction)
			gomega.Expect(err).Should(gomega.HaveOccurred())
		})

		ginkgo.It("should handle error when build external request", func() {
			mockGatewayCountryRepo.
				On("GetHealthyGatewayByCountryID", mock.Anything, transaction.CountryID).
				Return(gateway, nil).
				Once()

			mockSendTransactionClient.
				On("SendTransaction", mockCtx, mock.AnythingOfType("models.BuildExternalTransaction"), gateway.Name, gatewayConfig).
				Return(nil).
				Once()

			mockGatewayRepo.
				On("UpdateHealthStatus", mockCtx, gateway.ID, constants.UNHEALTHY).
				Return(nil).
				Once()

			mockTransactionRepo.
				On("UpdateGatewayIDByTransactionID", mockCtx, transaction.ID, gateway.ID).
				Return(errors.New("error")).
				Once()

			err := transactionHandler.TransactionProcessor(mockCtx, transaction)
			gomega.Expect(err).Should(gomega.HaveOccurred())
		})

		ginkgo.It("should process the transction successfully", func() {
			mockGatewayCountryRepo.
				On("GetHealthyGatewayByCountryID", mock.Anything, transaction.CountryID).
				Return(gateway, nil).
				Once()

			mockSendTransactionClient.
				On("SendTransaction", mockCtx, mock.AnythingOfType("models.BuildExternalTransaction"), gateway.Name, gatewayConfig).
				Return(nil).
				Once()

			mockGatewayRepo.
				On("UpdateHealthStatus", mockCtx, gateway.ID, constants.UNHEALTHY).
				Return(nil).
				Once()

			mockTransactionRepo.
				On("UpdateGatewayIDByTransactionID", mockCtx, transaction.ID, gateway.ID).
				Return(nil).
				Once()

			err := transactionHandler.TransactionProcessor(mockCtx, transaction)
			gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
		})
	})
})
