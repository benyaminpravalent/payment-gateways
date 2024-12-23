package rest

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	mocks "payment-gateway/mocks/services"
	"payment-gateway/models"
	"payment-gateway/pkg/constants"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
)

func TestTransactionRest(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "TransactionRest Suite")
}

type failingReader struct{}

func (f *failingReader) Read(p []byte) (n int, err error) {
	return 0, fmt.Errorf("simulated read error")
}

var _ = ginkgo.Describe("TransactionRest", func() {
	var (
		mockService    *mocks.TransactionService
		controller     *TransactionController
		e              *echo.Echo
		contextTimeout time.Duration
	)

	ginkgo.BeforeEach(func() {
		mockService = new(mocks.TransactionService)
		contextTimeout = 5 * time.Second
		e = echo.New()
		controller = &TransactionController{
			service:        mockService,
			contextTimeout: contextTimeout,
		}
	})

	ginkgo.Describe("Deposit Endpoint", func() {
		ginkgo.It("should return 202 Accepted when Deposit is successful", func() {
			// Mock request payload
			request := models.DepositRequest{
				Amount:    1000,
				Currency:  "USD",
				CountryID: 1,
				UserID:    1,
			}

			// Expected transaction response
			expectedTransaction := models.Transaction{
				ReferenceID: uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
				Amount:      request.Amount,
				Currency:    request.Currency,
				Type:        constants.DEPOSIT,
				Status:      constants.PENDING,
				CountryID:   request.CountryID,
				UserID:      request.UserID,
			}

			// Set up mock service behavior
			mockService.On("Deposit", mock.Anything, request).Return(expectedTransaction, nil)

			// Prepare HTTP request
			requestBody, _ := json.Marshal(request)
			req := httptest.NewRequest(http.MethodPost, "/transaction/deposit", bytes.NewReader(requestBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()

			c := e.NewContext(req, rec)
			c.SetPath("/transaction/deposit")

			// Invoke Deposit handler
			err := controller.Deposit(c)

			// Assertions
			gomega.Expect(err).To(gomega.BeNil())
			gomega.Expect(rec.Code).To(gomega.Equal(http.StatusAccepted))

			var response models.APIResponse
			err = json.Unmarshal(rec.Body.Bytes(), &response)
			gomega.Expect(err).To(gomega.BeNil())
			gomega.Expect(response.StatusCode).To(gomega.Equal(http.StatusAccepted))
			gomega.Expect(response.Message).To(gomega.Equal("Transaction is in process"))

			// Unmarshal response.Data into models.Transaction
			responseDataBytes, _ := json.Marshal(response.Data) // Convert interface{} to JSON
			var actualTransaction models.Transaction
			err = json.Unmarshal(responseDataBytes, &actualTransaction) // Unmarshal JSON to models.Transaction
			gomega.Expect(err).To(gomega.BeNil())

			// Compare the actual transaction with the expected transaction
			gomega.Expect(actualTransaction).To(gomega.Equal(expectedTransaction))

			// Verify mock behavior
			mockService.AssertCalled(ginkgo.GinkgoT(), "Deposit", mock.Anything, request)
		})

		ginkgo.It("should return 400 Bad Request when request payload is invalid", func() {
			// Prepare invalid HTTP request
			req := httptest.NewRequest(http.MethodPost, "/transaction/deposit", bytes.NewReader([]byte("invalid payload")))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()

			c := e.NewContext(req, rec)
			c.SetPath("/transaction/deposit")

			// Invoke Deposit handler
			err := controller.Deposit(c)

			// Assertions
			gomega.Expect(err).To(gomega.BeNil())
			gomega.Expect(rec.Code).To(gomega.Equal(http.StatusBadRequest))

			var response models.APIResponse
			err = json.Unmarshal(rec.Body.Bytes(), &response)
			gomega.Expect(err).To(gomega.BeNil())
			gomega.Expect(response.StatusCode).To(gomega.Equal(http.StatusBadRequest))
			gomega.Expect(response.Message).To(gomega.Equal("Invalid request payload"))
		})

		ginkgo.It("should return 500 Internal Server Error when Deposit fails", func() {
			// Mock request payload
			request := models.DepositRequest{
				Amount:    1000,
				Currency:  "USD",
				CountryID: 1,
				UserID:    1,
			}

			// Set up mock service behavior to simulate an error
			mockService.On("Deposit", mock.Anything, request).Return(models.Transaction{}, errors.New("deposit failed"))

			// Prepare HTTP request
			requestBody, _ := json.Marshal(request)
			req := httptest.NewRequest(http.MethodPost, "/transaction/deposit", bytes.NewReader(requestBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()

			c := e.NewContext(req, rec)
			c.SetPath("/transaction/deposit")

			// Invoke Deposit handler
			err := controller.Deposit(c)

			// Assertions
			gomega.Expect(err).To(gomega.BeNil())
			gomega.Expect(rec.Code).To(gomega.Equal(http.StatusInternalServerError))

			var response models.APIResponse
			err = json.Unmarshal(rec.Body.Bytes(), &response)
			gomega.Expect(err).To(gomega.BeNil())
			gomega.Expect(response.StatusCode).To(gomega.Equal(http.StatusInternalServerError))
			gomega.Expect(response.Message).To(gomega.Equal("Failed to process deposit"))

			// Verify mock behavior
			mockService.AssertCalled(ginkgo.GinkgoT(), "Deposit", mock.Anything, request)
		})
	})

	ginkgo.Describe("Withdraw Endpoint", func() {
		ginkgo.It("should return 202 Accepted when Deposit is successful", func() {
			// Mock request payload
			request := models.WithdrawalRequest{
				Amount:    1000,
				Currency:  "USD",
				CountryID: 1,
				UserID:    1,
			}

			// Expected transaction response
			expectedTransaction := models.Transaction{
				ReferenceID: uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
				Amount:      request.Amount,
				Currency:    request.Currency,
				Type:        constants.WITHDRAWAL,
				Status:      constants.PENDING,
				CountryID:   request.CountryID,
				UserID:      request.UserID,
			}

			// Set up mock service behavior
			mockService.On("Withdraw", mock.Anything, request).Return(expectedTransaction, nil)

			// Prepare HTTP request
			requestBody, _ := json.Marshal(request)
			req := httptest.NewRequest(http.MethodPost, "/transaction/withdraw", bytes.NewReader(requestBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()

			c := e.NewContext(req, rec)
			c.SetPath("/transaction/withdraw")

			// Invoke Deposit handler
			err := controller.Withdraw(c)

			// Assertions
			gomega.Expect(err).To(gomega.BeNil())
			gomega.Expect(rec.Code).To(gomega.Equal(http.StatusAccepted))

			var response models.APIResponse
			err = json.Unmarshal(rec.Body.Bytes(), &response)
			gomega.Expect(err).To(gomega.BeNil())
			gomega.Expect(response.StatusCode).To(gomega.Equal(http.StatusAccepted))
			gomega.Expect(response.Message).To(gomega.Equal("Transaction is in process"))

			// Unmarshal response.Data into models.Transaction
			responseDataBytes, _ := json.Marshal(response.Data) // Convert interface{} to JSON
			var actualTransaction models.Transaction
			err = json.Unmarshal(responseDataBytes, &actualTransaction) // Unmarshal JSON to models.Transaction
			gomega.Expect(err).To(gomega.BeNil())

			// Compare the actual transaction with the expected transaction
			gomega.Expect(actualTransaction).To(gomega.Equal(expectedTransaction))

			// Verify mock behavior
			mockService.AssertCalled(ginkgo.GinkgoT(), "Withdraw", mock.Anything, request)
		})

		ginkgo.It("should return 400 Bad Request when request payload is invalid", func() {
			// Prepare invalid HTTP request
			req := httptest.NewRequest(http.MethodPost, "/transaction/withdraw", bytes.NewReader([]byte("invalid payload")))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()

			c := e.NewContext(req, rec)
			c.SetPath("/transaction/withdraw")

			// Invoke Deposit handler
			err := controller.Withdraw(c)

			// Assertions
			gomega.Expect(err).To(gomega.BeNil())
			gomega.Expect(rec.Code).To(gomega.Equal(http.StatusBadRequest))

			var response models.APIResponse
			err = json.Unmarshal(rec.Body.Bytes(), &response)
			gomega.Expect(err).To(gomega.BeNil())
			gomega.Expect(response.StatusCode).To(gomega.Equal(http.StatusBadRequest))
			gomega.Expect(response.Message).To(gomega.Equal("Invalid request payload"))
		})

		ginkgo.It("should return 500 Internal Server Error when Withdraw fails", func() {
			// Mock request payload
			request := models.WithdrawalRequest{
				Amount:    1000,
				Currency:  "USD",
				CountryID: 1,
				UserID:    1,
			}

			// Set up mock service behavior to simulate an error
			mockService.On("Withdraw", mock.Anything, request).Return(models.Transaction{}, errors.New("withdraw failed"))

			// Prepare HTTP request
			requestBody, _ := json.Marshal(request)
			req := httptest.NewRequest(http.MethodPost, "/transaction/withdraw", bytes.NewReader(requestBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()

			c := e.NewContext(req, rec)
			c.SetPath("/transaction/withdraw")

			// Invoke Deposit handler
			err := controller.Withdraw(c)

			// Assertions
			gomega.Expect(err).To(gomega.BeNil())
			gomega.Expect(rec.Code).To(gomega.Equal(http.StatusInternalServerError))

			var response models.APIResponse
			err = json.Unmarshal(rec.Body.Bytes(), &response)
			gomega.Expect(err).To(gomega.BeNil())
			gomega.Expect(response.StatusCode).To(gomega.Equal(http.StatusInternalServerError))
			gomega.Expect(response.Message).To(gomega.Equal("Failed to process withdrawal"))

			// Verify mock behavior
			mockService.AssertCalled(ginkgo.GinkgoT(), "Withdraw", mock.Anything, request)
		})
	})

	ginkgo.Describe("Withdraw Endpoint", func() {
		ginkgo.It("should return 200 OK when callback is processed successfully", func() {
			// Mock request payload
			request := models.TransactionCallbackRequest{
				ReferenceID: "123e4567-e89b-12d3-a456-426614174000",
				Status:      constants.COMPLETED,
			}

			// Use a WaitGroup to wait for the goroutine to complete
			var wg sync.WaitGroup
			wg.Add(1)

			// Wrap the mock service to call Done on the WaitGroup
			mockService.On("TransactionCallback", mock.Anything, &request).Run(func(args mock.Arguments) {
				log.Println("Mock TransactionCallback called.") // Debug log
				wg.Done()
			}).Return(nil)

			// Prepare HTTP request
			requestBody, _ := json.Marshal(request)
			req := httptest.NewRequest(http.MethodPost, "/transaction/callback", bytes.NewReader(requestBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()

			c := e.NewContext(req, rec)
			c.SetPath("/transaction/callback")

			// Invoke TransactionCallback handler
			err := controller.TransactionCallback(c)

			// Wait for the goroutine to finish
			wg.Wait()

			// Assertions
			gomega.Expect(err).To(gomega.BeNil())
			gomega.Expect(rec.Code).To(gomega.Equal(http.StatusOK))

			var response models.APIResponse
			err = json.Unmarshal(rec.Body.Bytes(), &response)
			gomega.Expect(err).To(gomega.BeNil())
			gomega.Expect(response.StatusCode).To(gomega.Equal(http.StatusOK))
			gomega.Expect(response.Message).To(gomega.Equal("Callback received"))

			// Verify mock behavior
			mockService.AssertCalled(ginkgo.GinkgoT(), "TransactionCallback", mock.Anything, &request)
		})

		ginkgo.It("should return 400 Bad Request when io.ReadAll fails", func() {
			// Create a custom io.Reader that always returns an error
			errorReader := io.NopCloser(&failingReader{})

			// Prepare an HTTP request with the failing reader as the body
			req := httptest.NewRequest(http.MethodPost, "/transaction/callback", errorReader)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()

			c := e.NewContext(req, rec)
			c.SetPath("/transaction/callback")

			// Invoke TransactionCallback handler
			err := controller.TransactionCallback(c)

			// Assertions
			gomega.Expect(err).To(gomega.BeNil())
			gomega.Expect(rec.Code).To(gomega.Equal(http.StatusBadRequest))

			var response models.APIResponse
			err = json.Unmarshal(rec.Body.Bytes(), &response)
			gomega.Expect(err).To(gomega.BeNil())
			gomega.Expect(response.StatusCode).To(gomega.Equal(http.StatusBadRequest))
			gomega.Expect(response.Message).To(gomega.Equal("Invalid request payload"))
		})

	})
})
