package rest

import (
	"bytes"
	"context"
	"io"
	"log"
	"net/http"
	"time"

	"payment-gateway/models"
	"payment-gateway/pkg/utils"

	"github.com/labstack/echo/v4"
)

type ITransactionService interface {
	Deposit(ctx context.Context, request models.DepositRequest) (models.Transaction, error)
	Withdraw(ctx context.Context, request models.WithdrawalRequest) (models.Transaction, error)
	TransactionCallback(ctx context.Context, request *models.TransactionCallbackRequest) error
}

type TransactionController struct {
	service        ITransactionService
	contextTimeout time.Duration
}

func InstallTransactionController(e *echo.Echo, s ITransactionService, contextTimeout time.Duration) {
	controller := &TransactionController{
		service:        s,
		contextTimeout: contextTimeout,
	}

	transactionGroup := e.Group("/transaction")

	transactionGroup.POST("/deposit", controller.Deposit)
	transactionGroup.POST("/withdraw", controller.Withdraw)
	transactionGroup.POST("/callback", controller.TransactionCallback)
}

func (controller *TransactionController) Deposit(c echo.Context) error {
	ctx, cancel := context.WithTimeout(c.Request().Context(), controller.contextTimeout)
	defer cancel()

	var request models.DepositRequest
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, models.APIResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid request payload",
			Data:       nil,
		})
	}

	result, err := controller.service.Deposit(ctx, request)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.APIResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to process deposit",
			Data:       nil,
		})
	}

	response := models.APIResponse{
		StatusCode: http.StatusAccepted,
		Message:    "Transaction is in process",
		Data:       result,
	}

	return c.JSON(http.StatusAccepted, response)
}

func (controller *TransactionController) Withdraw(c echo.Context) error {
	ctx, cancel := context.WithTimeout(c.Request().Context(), controller.contextTimeout)
	defer cancel()

	var request models.WithdrawalRequest
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, models.APIResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid request payload",
			Data:       nil,
		})
	}

	result, err := controller.service.Withdraw(ctx, request)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.APIResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to process withdrawal",
			Data:       nil,
		})
	}

	response := models.APIResponse{
		StatusCode: http.StatusAccepted,
		Message:    "Transaction is in process",
		Data:       result,
	}

	return c.JSON(http.StatusAccepted, response)
}

func (controller *TransactionController) TransactionCallback(c echo.Context) error {
	// Read and copy the request body
	bodyBytes, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.APIResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid request payload",
		})
	}

	// Reset the body so it can be reused
	c.Request().Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	go func(body []byte) {
		ctx, cancel := context.WithTimeout(context.Background(), controller.contextTimeout)
		defer cancel()

		var request models.TransactionCallbackRequest

		// Create a new request object with the copied body
		req := &http.Request{
			Body:   io.NopCloser(bytes.NewBuffer(body)),
			Header: c.Request().Header,
		}

		// Decode the request using utils
		err := utils.DecodeRequest(req, &request)
		if err != nil {
			log.Printf("Invalid request payload: %v", err)
			return
		}

		if err := controller.service.TransactionCallback(ctx, &request); err != nil {
			log.Printf("Failed to process callback: %v", err)
		}
	}(bodyBytes)

	return c.JSON(http.StatusOK, models.APIResponse{
		StatusCode: http.StatusOK,
		Message:    "Callback received",
	})
}
