package repositories

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type IGatewayRepository interface {
	UpdateHealthStatus(ctx context.Context, gatewayID int, healthStatus string) error
}

type GatewayRepository struct {
	db *sqlx.DB
}

func NewGatewayRepository(db *sqlx.DB) *GatewayRepository {
	return &GatewayRepository{
		db: db,
	}
}

// UpdateHealthStatus updates the health_status of a gateway by its ID
func (r *GatewayRepository) UpdateHealthStatus(ctx context.Context, gatewayID int, healthStatus string) error {
	query := `
		UPDATE gateways
		SET health_status = $1,
		    last_checked_at = CURRENT_TIMESTAMP,
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = $2;
	`

	result, err := r.db.ExecContext(ctx, query, healthStatus, gatewayID)
	if err != nil {
		return fmt.Errorf("failed to update health_status for gatewayID %d: %w", gatewayID, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to retrieve rows affected for gatewayID %d: %w", gatewayID, err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no gateway found with id %d", gatewayID)
	}

	return nil
}
