package models

import (
	"time"

	"github.com/google/uuid"
)

type Transaction struct {
	ID          int       `json:"id" db:"id"`
	ReferenceID uuid.UUID `json:"reference_id" db:"reference_id"`
	Amount      float64   `json:"amount" db:"amount"`
	Currency    string    `json:"currency"`
	Type        string    `json:"type" db:"type"`     // deposit/withdrawal
	Status      string    `json:"status" db:"status"` // pending, completed, failed
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
	GatewayID   int       `json:"gateway_id" db:"gateway_id"`
	CountryID   int       `json:"country_id" db:"country_id"`
	UserID      int       `json:"user_id" db:"user_id"`
}

type SendTransactionRequest struct {
	ReferenceID string  `json:"reference_id"`
	Amount      float64 `json:"amount"`
	UserID      int     `json:"user_id"`
	Currency    string  `json:"currency"`
}

type EncryptedTransactionRequest struct {
	EncryptedData string `json:"encrypted_data"`
}

type DepositRequest struct {
	UserID    int     `json:"user_id" validate:"required"`
	Amount    float64 `json:"amount" validate:"required,gt=0"`
	Currency  string  `json:"currency" validate:"required"`
	CountryID int     `json:"country_id" validate:"required"`
}

type WithdrawalRequest struct {
	UserID    int     `json:"user_id" validate:"required"`
	Amount    float64 `json:"amount" validate:"required,gt=0"`
	Currency  string  `json:"currency" validate:"required"`
	CountryID int     `json:"country_id" validate:"required"`
}

type DepositResponse struct {
	ReferenceID string  `json:"reference_id"`
	UserID      int     `json:"user_id"`
	Status      string  `json:"status"`
	Amount      float64 `json:"amount"`
	Currency    string  `json:"currency"`
	CountryID   int     `json:"country_id"`
}

type BuildExternalTransaction struct {
	Request     string
	ContentType string
}

type TransactionCallbackRequest struct {
	ReferenceID string  `json:"id" xml:"Body>TransactionCallbackRequest>id"`
	Amount      float64 `json:"amount" xml:"Body>TransactionCallbackRequest>amount"`
	Currency    string  `json:"currency" xml:"Body>TransactionCallbackRequest>currency"`
	Status      string  `json:"status" xml:"Body>TransactionCallbackRequest>status"`
}
