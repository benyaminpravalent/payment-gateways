package services

import (
	"context"
	"database/sql"
	"errors"
	"sync"
	"testing"
	"time"

	"payment-gateway/internal/kafka"
	mockKafka "payment-gateway/mocks/kafka"
	mocksRepository "payment-gateway/mocks/repositories"
	"payment-gateway/pkg/constants"

	"payment-gateway/models"

	"github.com/google/uuid"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
)

func TestTransactionService(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "Service Suite")
}

var _ = ginkgo.Describe("TransactionService", func() {
	var (
		mockRepo           *mocksRepository.TransactionRepository
		mockKafkaProducer  *mockKafka.MockKafkaProducer
		transactionService *TransactionService
	)

	ginkgo.BeforeEach(func() {
		mockRepo = new(mocksRepository.TransactionRepository)
		mockKafkaProducer = new(mockKafka.MockKafkaProducer)
		transactionService = NewTransactionService(mockRepo, mockKafkaProducer)
	})

	ginkgo.Describe("Deposit", func() {
		ginkgo.It("should successfully process a deposit transaction", func() {
			request := models.DepositRequest{
				Amount:    1000,
				Currency:  "USD",
				CountryID: 1,
				UserID:    123,
			}

			mockRepo.On("InsertTransaction", mock.Anything, mock.MatchedBy(func(tx *models.Transaction) bool {
				gomega.Expect(tx.Amount).To(gomega.Equal(request.Amount))
				gomega.Expect(tx.Currency).To(gomega.Equal(request.Currency))
				gomega.Expect(tx.Type).To(gomega.Equal(constants.DEPOSIT))
				gomega.Expect(tx.Status).To(gomega.Equal(constants.PENDING))
				gomega.Expect(tx.CountryID).To(gomega.Equal(request.CountryID))
				gomega.Expect(tx.UserID).To(gomega.Equal(request.UserID))
				gomega.Expect(tx.ReferenceID).NotTo(gomega.BeNil())
				gomega.Expect(tx.CreatedAt).NotTo(gomega.BeZero())
				gomega.Expect(tx.UpdatedAt).NotTo(gomega.BeZero())
				return true
			})).Return(nil)

			mockKafkaProducer.On("ProduceMessage", mock.Anything, kafka.SendTransactionKafkaTopic).Return(nil)

			result, err := transactionService.Deposit(context.Background(), request)

			gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
			gomega.Expect(result.Amount).Should(gomega.Equal(request.Amount))
			gomega.Expect(result.Currency).Should(gomega.Equal(request.Currency))
			gomega.Expect(result.Type).Should(gomega.Equal(constants.DEPOSIT))
			mockRepo.AssertCalled(ginkgo.GinkgoT(), "InsertTransaction", mock.Anything, mock.AnythingOfType("*models.Transaction"))
		})

		ginkgo.It("should return error when InsertTransaction fails", func() {
			request := models.DepositRequest{
				Amount:    1000,
				Currency:  "USD",
				CountryID: 1,
				UserID:    123,
			}

			mockRepo.On("InsertTransaction", mock.Anything, mock.MatchedBy(func(tx *models.Transaction) bool {
				gomega.Expect(tx.Amount).To(gomega.Equal(request.Amount))
				gomega.Expect(tx.Currency).To(gomega.Equal(request.Currency))
				return true
			})).Return(errors.New("insert error"))

			result, err := transactionService.Deposit(context.Background(), request)

			gomega.Expect(err).Should(gomega.HaveOccurred())
			gomega.Expect(result).Should(gomega.Equal(models.Transaction{}))
			gomega.Expect(err.Error()).Should(gomega.ContainSubstring("insert error"))
			mockRepo.AssertCalled(ginkgo.GinkgoT(), "InsertTransaction", mock.Anything, mock.AnythingOfType("*models.Transaction"))
		})

		ginkgo.It("should invoke Kafka producer even if it fails", func() {
			request := models.DepositRequest{
				Amount:    1000,
				Currency:  "USD",
				CountryID: 1,
				UserID:    123,
			}

			transaction := models.Transaction{
				ReferenceID: uuid.New(),
				Amount:      request.Amount,
				Currency:    request.Currency,
				Type:        constants.DEPOSIT,
				Status:      constants.PENDING,
				CountryID:   request.CountryID,
				UserID:      request.UserID,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			}

			var wg sync.WaitGroup
			wg.Add(1)

			mockRepo.On("InsertTransaction", mock.Anything, mock.MatchedBy(func(tx *models.Transaction) bool {
				tx.ReferenceID = transaction.ReferenceID
				return true
			})).Return(nil)

			mockKafkaProducer.On("ProduceMessage", mock.Anything, kafka.SendTransactionKafkaTopic).Run(func(args mock.Arguments) {
				wg.Done()
			}).Return(errors.New("kafka error"))

			result, err := transactionService.Deposit(context.Background(), request)
			gomega.Expect(err).ShouldNot(gomega.HaveOccurred())

			wg.Wait()

			gomega.Expect(result.Amount).Should(gomega.Equal(request.Amount))
			mockRepo.AssertCalled(ginkgo.GinkgoT(), "InsertTransaction", mock.Anything, mock.AnythingOfType("*models.Transaction"))
			mockKafkaProducer.AssertCalled(ginkgo.GinkgoT(), "ProduceMessage", mock.Anything, kafka.SendTransactionKafkaTopic)
		})
	})

	ginkgo.Describe("Withdraw", func() {
		ginkgo.It("should successfully process a deposit transaction", func() {
			request := models.WithdrawalRequest{
				Amount:    1000,
				Currency:  "USD",
				CountryID: 1,
				UserID:    123,
			}

			mockRepo.On("InsertTransaction", mock.Anything, mock.MatchedBy(func(tx *models.Transaction) bool {
				gomega.Expect(tx.Amount).To(gomega.Equal(request.Amount))
				gomega.Expect(tx.Currency).To(gomega.Equal(request.Currency))
				gomega.Expect(tx.Type).To(gomega.Equal(constants.WITHDRAWAL))
				gomega.Expect(tx.Status).To(gomega.Equal(constants.PENDING))
				gomega.Expect(tx.CountryID).To(gomega.Equal(request.CountryID))
				gomega.Expect(tx.UserID).To(gomega.Equal(request.UserID))
				gomega.Expect(tx.ReferenceID).NotTo(gomega.BeNil())
				gomega.Expect(tx.CreatedAt).NotTo(gomega.BeZero())
				gomega.Expect(tx.UpdatedAt).NotTo(gomega.BeZero())
				return true
			})).Return(nil)

			mockKafkaProducer.On("ProduceMessage", mock.Anything, kafka.SendTransactionKafkaTopic).Return(nil)

			result, err := transactionService.Withdraw(context.Background(), request)

			gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
			gomega.Expect(result.Amount).Should(gomega.Equal(request.Amount))
			gomega.Expect(result.Currency).Should(gomega.Equal(request.Currency))
			gomega.Expect(result.Type).Should(gomega.Equal(constants.WITHDRAWAL))
			mockRepo.AssertCalled(ginkgo.GinkgoT(), "InsertTransaction", mock.Anything, mock.AnythingOfType("*models.Transaction"))
		})

		ginkgo.It("should return error when InsertTransaction fails", func() {
			request := models.WithdrawalRequest{
				Amount:    1000,
				Currency:  "USD",
				CountryID: 1,
				UserID:    123,
			}

			mockRepo.On("InsertTransaction", mock.Anything, mock.MatchedBy(func(tx *models.Transaction) bool {
				gomega.Expect(tx.Amount).To(gomega.Equal(request.Amount))
				gomega.Expect(tx.Currency).To(gomega.Equal(request.Currency))
				return true
			})).Return(errors.New("insert error"))

			result, err := transactionService.Withdraw(context.Background(), request)

			gomega.Expect(err).Should(gomega.HaveOccurred())
			gomega.Expect(result).Should(gomega.Equal(models.Transaction{}))
			gomega.Expect(err.Error()).Should(gomega.ContainSubstring("insert error"))
			mockRepo.AssertCalled(ginkgo.GinkgoT(), "InsertTransaction", mock.Anything, mock.AnythingOfType("*models.Transaction"))
		})

		ginkgo.It("should invoke Kafka producer even if it fails", func() {
			request := models.WithdrawalRequest{
				Amount:    1000,
				Currency:  "USD",
				CountryID: 1,
				UserID:    123,
			}

			transaction := models.Transaction{
				ReferenceID: uuid.New(),
				Amount:      request.Amount,
				Currency:    request.Currency,
				Type:        constants.WITHDRAWAL,
				Status:      constants.PENDING,
				CountryID:   request.CountryID,
				UserID:      request.UserID,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			}

			// Use a WaitGroup to wait for the goroutine
			var wg sync.WaitGroup
			wg.Add(1)

			// Mock repository behavior
			mockRepo.On("InsertTransaction", mock.Anything, mock.MatchedBy(func(tx *models.Transaction) bool {
				// Ensure the transaction contains expected values
				tx.ReferenceID = transaction.ReferenceID
				return true
			})).Return(nil)

			mockKafkaProducer.On("ProduceMessage", mock.Anything, kafka.SendTransactionKafkaTopic).Run(func(args mock.Arguments) {
				wg.Done()
			}).Return(errors.New("kafka error"))

			result, err := transactionService.Withdraw(context.Background(), request)

			wg.Wait()

			gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
			gomega.Expect(result.Amount).Should(gomega.Equal(request.Amount))
			mockRepo.AssertCalled(ginkgo.GinkgoT(), "InsertTransaction", mock.Anything, mock.AnythingOfType("*models.Transaction"))
			mockKafkaProducer.AssertCalled(ginkgo.GinkgoT(), "ProduceMessage", mock.Anything, kafka.SendTransactionKafkaTopic)
		})
	})

	ginkgo.Describe("TransactionCallback", func() {
		ginkgo.It("should successfully update transaction status", func() {
			request := &models.TransactionCallbackRequest{
				ReferenceID: "test-reference-id",
				Status:      "success",
			}

			mockRepo.On("UpdateTransactionStatusByReferenceID", mock.Anything, request.ReferenceID, request.Status).Return(nil)

			err := transactionService.TransactionCallback(context.Background(), request)

			gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
			mockRepo.AssertCalled(ginkgo.GinkgoT(), "UpdateTransactionStatusByReferenceID", mock.Anything, request.ReferenceID, request.Status)
		})

		ginkgo.It("should return error when UpdateTransactionStatusByReferenceID fails", func() {
			request := &models.TransactionCallbackRequest{
				ReferenceID: "test-reference-id",
				Status:      "failed",
			}

			mockRepo.On("UpdateTransactionStatusByReferenceID", mock.Anything, request.ReferenceID, request.Status).Return(errors.New("update error"))

			err := transactionService.TransactionCallback(context.Background(), request)

			gomega.Expect(err).Should(gomega.HaveOccurred())
			gomega.Expect(err.Error()).Should(gomega.ContainSubstring("update error"))
			mockRepo.AssertCalled(ginkgo.GinkgoT(), "UpdateTransactionStatusByReferenceID", mock.Anything, request.ReferenceID, request.Status)
		})

		ginkgo.It("should return error when no rows are affected", func() {
			request := &models.TransactionCallbackRequest{
				ReferenceID: "nonexistent-reference-id",
				Status:      "success",
			}

			mockRepo.On("UpdateTransactionStatusByReferenceID", mock.Anything, request.ReferenceID, request.Status).Return(sql.ErrNoRows)

			err := transactionService.TransactionCallback(context.Background(), request)

			gomega.Expect(err).Should(gomega.HaveOccurred())
			gomega.Expect(err).Should(gomega.Equal(sql.ErrNoRows))
			mockRepo.AssertCalled(ginkgo.GinkgoT(), "UpdateTransactionStatusByReferenceID", mock.Anything, request.ReferenceID, request.Status)
		})

	})
})
