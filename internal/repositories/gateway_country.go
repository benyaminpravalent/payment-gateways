package repositories

import (
	"context"
	"log"
	"payment-gateway/models"

	"github.com/jmoiron/sqlx"
)

type IGatewayCountryRepository interface {
	GetHealthyGatewayByCountryID(ctx context.Context, countryID int) (*models.GatewayDetail, error)
}

type GatewayCountryRepository struct {
	db *sqlx.DB
}

func NewGatewayCountryRepository(db *sqlx.DB) *GatewayCountryRepository {
	return &GatewayCountryRepository{db: db}
}

func (r *GatewayCountryRepository) GetHealthyGatewayByCountryID(ctx context.Context, countryID int) (*models.GatewayDetail, error) {
	var gatewayDetail models.GatewayDetail
	query := `
		SELECT
			g.id,
			g.name,
			g.data_format_supported,
			g.health_status,
			g.last_checked_at,
			g.created_at,
			g.updated_at,
			gc.priority,
			gc.country_id,
			c.currency
		FROM
			gateway_countries gc
		JOIN 
			gateways g ON gc.gateway_id = g.id
		JOIN
			countries c ON gc.country_id = c.id
		WHERE
			gc.country_id = $1
		AND
			g.health_status = 'healthy'
		ORDER BY 
			gc.priority ASC
		LIMIT 1;
	`

	err := r.db.GetContext(ctx, &gatewayDetail, query, countryID)
	if err != nil {
		log.Printf("Error fetching gateway detail for country_id %d: %v", countryID, err)
		return nil, err
	}

	return &gatewayDetail, nil
}
