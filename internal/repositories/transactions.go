package repositories

import (
	"context"
	"database/sql"
	"log"

	"payment-gateway/models"

	"github.com/jmoiron/sqlx"
)

// TransactionRepository handles database operations for the transactions table
type TransactionRepository struct {
	db *sqlx.DB
}

// NewTransactionRepository creates a new instance of TransactionRepository
func NewTransactionRepository(db *sqlx.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

// InsertTransaction inserts a new transaction into the database
func (r *TransactionRepository) InsertTransaction(ctx context.Context, transaction *models.Transaction) error {
	query := `
		INSERT INTO transactions (
			reference_id, amount, currency, type, status, created_at, updated_at, country_id, user_id
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9
		) RETURNING id;
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		transaction.ReferenceID,
		transaction.Amount,
		transaction.Currency,
		transaction.Type,
		transaction.Status,
		transaction.CreatedAt,
		transaction.UpdatedAt,
		transaction.CountryID,
		transaction.UserID,
	).Scan(&transaction.ID)
	if err != nil {
		log.Printf("Error inserting transaction: %v", err)
		return err
	}
	return nil
}

// GetTransactionByReferenceID retrieves a transaction by its reference ID
func (r *TransactionRepository) GetTransactionByReferenceID(ctx context.Context, referenceID string) (*models.Transaction, error) {
	var transaction models.Transaction
	query := `
		SELECT 
			id, reference_id, amount, currency, type, status, created_at, updated_at, gateway_id, country_id, user_id
		FROM 
			transactions
		WHERE 
			reference_id = $1;
	`
	err := r.db.GetContext(ctx, &transaction, query, referenceID)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("Transaction with Reference ID %s not found", referenceID)
			return nil, nil
		}
		log.Printf("Error fetching transaction with Reference ID %s: %v", referenceID, err)
		return nil, err
	}
	return &transaction, nil
}

// UpdateTransactionStatusByReferenceID updates the status of a transaction by its reference ID
func (r *TransactionRepository) UpdateTransactionStatusByReferenceID(ctx context.Context, referenceID string, status string) error {
	query := `
		UPDATE transactions
		SET status = $1, updated_at = NOW()
		WHERE reference_id = $2;
	`
	result, err := r.db.ExecContext(ctx, query, status, referenceID)
	if err != nil {
		log.Printf("Error updating status for transaction with Reference ID %s: %v", referenceID, err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Error getting rows affected for transaction with Reference ID %s: %v", referenceID, err)
		return err
	}

	if rowsAffected == 0 {
		log.Printf("No transaction found with Reference ID %s to update", referenceID)
		return sql.ErrNoRows
	}

	return nil
}

func (r *TransactionRepository) UpdateGatewayIDByTransactionID(ctx context.Context, transactionID int, gatewayID int) error {
	query := `
		UPDATE transactions
		SET gateway_id = $1, updated_at = NOW()
		WHERE id = $2;
	`
	result, err := r.db.ExecContext(ctx, query, gatewayID, transactionID)
	if err != nil {
		log.Printf("Error updating gateway_id for transaction ID %d: %v", transactionID, err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Error getting rows affected for transaction ID %d: %v", transactionID, err)
		return err
	}

	if rowsAffected == 0 {
		log.Printf("No transaction found with ID %d to update", transactionID)
		return sql.ErrNoRows
	}

	return nil
}
