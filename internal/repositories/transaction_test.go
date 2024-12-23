package repositories

import (
	"context"
	"database/sql"
	"errors"

	"payment-gateway/models"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("TransactionRepository", func() {
	var (
		mockDB      *sqlx.DB
		sqlMock     sqlmock.Sqlmock
		repo        *TransactionRepository
		ctx         context.Context
		transaction *models.Transaction
	)

	ginkgo.BeforeEach(func() {
		var err error
		sqlDB, mock, err := sqlmock.New()
		gomega.Expect(err).ShouldNot(gomega.HaveOccurred())

		mockDB = sqlx.NewDb(sqlDB, "sqlmock")
		sqlMock = mock
		repo = NewTransactionRepository(mockDB)

		ctx = context.Background()
		transaction = &models.Transaction{
			ID:          1,
			ReferenceID: uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			Amount:      100.0,
			Currency:    "USD",
			Type:        "payment",
			Status:      "pending",
			CountryID:   1,
			UserID:      1,
		}
	})

	ginkgo.AfterEach(func() {
		err := sqlMock.ExpectationsWereMet()
		gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
	})

	ginkgo.Describe("InsertTransaction", func() {
		ginkgo.It("should successfully insert a transaction", func() {
			sqlMock.ExpectQuery(`INSERT INTO transactions`).
				WithArgs(
					transaction.ReferenceID, transaction.Amount, transaction.Currency,
					transaction.Type, transaction.Status, transaction.CreatedAt,
					transaction.UpdatedAt, transaction.CountryID, transaction.UserID,
				).
				WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

			err := repo.InsertTransaction(ctx, transaction)
			gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
			gomega.Expect(transaction.ID).Should(gomega.Equal(1))
		})

		ginkgo.It("should return error when insertion fails", func() {
			dbError := errors.New("database error")
			sqlMock.ExpectQuery(`INSERT INTO transactions`).
				WithArgs(
					transaction.ReferenceID, transaction.Amount, transaction.Currency,
					transaction.Type, transaction.Status, transaction.CreatedAt,
					transaction.UpdatedAt, transaction.CountryID, transaction.UserID,
				).
				WillReturnError(dbError)

			err := repo.InsertTransaction(ctx, transaction)
			gomega.Expect(err).Should(gomega.HaveOccurred())
			gomega.Expect(err.Error()).Should(gomega.ContainSubstring(dbError.Error()))
		})
	})

	ginkgo.Describe("UpdateTransactionStatusByReferenceID", func() {
		ginkgo.It("should successfully update transaction status", func() {
			sqlMock.ExpectExec(`UPDATE transactions`).
				WithArgs("completed", "ref123").
				WillReturnResult(sqlmock.NewResult(1, 1))

			err := repo.UpdateTransactionStatusByReferenceID(ctx, "ref123", "completed")
			gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
		})

		ginkgo.It("should return error when no rows are affected", func() {
			sqlMock.ExpectExec(`UPDATE transactions`).
				WithArgs("completed", "ref123").
				WillReturnResult(sqlmock.NewResult(1, 0))

			err := repo.UpdateTransactionStatusByReferenceID(ctx, "ref123", "completed")
			gomega.Expect(err).Should(gomega.HaveOccurred())
			gomega.Expect(err).Should(gomega.Equal(sql.ErrNoRows))
		})

		ginkgo.It("should return error when update query fails", func() {
			dbError := errors.New("database error")
			sqlMock.ExpectExec(`UPDATE transactions`).
				WithArgs("completed", "ref123").
				WillReturnError(dbError)

			err := repo.UpdateTransactionStatusByReferenceID(ctx, "ref123", "completed")
			gomega.Expect(err).Should(gomega.HaveOccurred())
			gomega.Expect(err.Error()).Should(gomega.ContainSubstring(dbError.Error()))
		})
	})

	ginkgo.Describe("UpdateGatewayIDByTransactionID", func() {
		ginkgo.It("should successfully update gateway ID", func() {
			sqlMock.ExpectExec(`UPDATE transactions`).
				WithArgs(2, 1).
				WillReturnResult(sqlmock.NewResult(1, 1))

			err := repo.UpdateGatewayIDByTransactionID(ctx, 1, 2)
			gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
		})

		ginkgo.It("should return error when no rows are affected", func() {
			sqlMock.ExpectExec(`UPDATE transactions`).
				WithArgs(2, 1).
				WillReturnResult(sqlmock.NewResult(1, 0))

			err := repo.UpdateGatewayIDByTransactionID(ctx, 1, 2)
			gomega.Expect(err).Should(gomega.HaveOccurred())
			gomega.Expect(err).Should(gomega.Equal(sql.ErrNoRows))
		})

		ginkgo.It("should return error when update query fails", func() {
			dbError := errors.New("database error")
			sqlMock.ExpectExec(`UPDATE transactions`).
				WithArgs(2, 1).
				WillReturnError(dbError)

			err := repo.UpdateGatewayIDByTransactionID(ctx, 1, 2)
			gomega.Expect(err).Should(gomega.HaveOccurred())
			gomega.Expect(err.Error()).Should(gomega.ContainSubstring(dbError.Error()))
		})
	})
})
