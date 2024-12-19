package rest

import (
	"context"
	"net/http"
	"time"

	"payment-gateway/models"

	"github.com/labstack/echo/v4"
)

type ITransactionService interface {
	Deposit(ctx context.Context, request models.DepositRequest) (models.Transaction, error)
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

	transactionGroup.POST("/deposit", controller.deposit)
}

func (controller *TransactionController) deposit(c echo.Context) error {
	ctx, cancel := context.WithTimeout(c.Request().Context(), controller.contextTimeout)
	defer cancel()

	var request models.DepositRequest
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request payload",
		})
	}

	result, err := controller.service.Deposit(ctx, request)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to process deposit",
		})
	}

	response := models.APIResponse{
		StatusCode: http.StatusAccepted,
		Message:    "Transaction is in process",
		Data:       result,
	}
	
	return c.JSON(http.StatusAccepted, response)
}
