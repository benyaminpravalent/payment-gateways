package models

import "time"

type GatewayDetail struct {
	ID                  int       `db:"id"`
	Name                string    `db:"name"`
	DataFormatSupported string    `db:"data_format_supported"`
	HealthStatus        string    `db:"health_status"`
	LastCheckedAt       time.Time `db:"last_checked_at"`
	CreatedAt           time.Time `db:"created_at"`
	UpdatedAt           time.Time `db:"updated_at"`
	Priority            int       `db:"priority"`
	CountryID           int       `db:"country_id"`
	Currency            string    `db:"currency"`
}

type GatewayConfig struct {
	GatewayUrl        string
	GatewayApiKey     string
	GatewayPrivateKey string
}

type GatewayCallback struct {
	ReferenceID     string    `json:"reference_id"`     // Transaction reference ID
	Status          string    `json:"status"`           // Transaction status
	Amount          float64   `json:"amount"`           // Transaction amount
	Currency        string    `json:"currency"`         // Currency code
	GatewayResponse string    `json:"gateway_response"` // Response message from the gateway
	Timestamp       time.Time `json:"timestamp"`        // Callback timestamp
}
