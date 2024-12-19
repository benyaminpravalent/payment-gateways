package repositories

import (
	"context"
	"log"
	"payment-gateway/models"

	"github.com/jmoiron/sqlx"
)

type CountryRepository struct {
	db *sqlx.DB
}

func NewCountryRepository(db *sqlx.DB) *CountryRepository {
	return &CountryRepository{db: db}
}

func (r *CountryRepository) GetCountryByCode(ctx context.Context, code string) (*models.Country, error) {
	var country models.Country
	query := "SELECT id, name, code, currency, created_at, updated_at FROM countries WHERE code = $1"
	err := r.db.GetContext(ctx, &country, query, code)
	if err != nil {
		log.Printf("Error fetching country with code %s: %v", code, err)
		return nil, err
	}

	return &country, nil
}
