package repositories

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

func TestGatewayRepository(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "GatewayRepository Suite")
}

var _ = ginkgo.Describe("GatewayRepository", func() {
	var (
		mockDB    *sqlx.DB
		sqlMock   sqlmock.Sqlmock
		repo      *GatewayRepository
		ctx       context.Context
		gatewayID int
		status    string
	)

	ginkgo.BeforeEach(func() {
		var err error
		sqlDB, mock, err := sqlmock.New()
		gomega.Expect(err).ShouldNot(gomega.HaveOccurred())

		mockDB = sqlx.NewDb(sqlDB, "sqlmock")
		sqlMock = mock
		repo = NewGatewayRepository(mockDB)

		ctx = context.Background()
		gatewayID = 1
		status = "healthy"
	})

	ginkgo.AfterEach(func() {
		err := sqlMock.ExpectationsWereMet()
		gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
	})

	ginkgo.Describe("UpdateHealthStatus", func() {
		ginkgo.It("should successfully update health status", func() {
			sqlMock.ExpectExec(`UPDATE gateways`).
				WithArgs(status, gatewayID).
				WillReturnResult(sqlmock.NewResult(1, 1)) // 1 row affected

			err := repo.UpdateHealthStatus(ctx, gatewayID, status)
			gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
		})

		ginkgo.It("should return error when no rows are affected", func() {
			sqlMock.ExpectExec(`UPDATE gateways`).
				WithArgs(status, gatewayID).
				WillReturnResult(sqlmock.NewResult(1, 0)) // 0 rows affected

			err := repo.UpdateHealthStatus(ctx, gatewayID, status)
			gomega.Expect(err).Should(gomega.HaveOccurred())
			gomega.Expect(err.Error()).Should(gomega.ContainSubstring("no gateway found with id"))
		})

		ginkgo.It("should return error when database query fails", func() {
			dbError := errors.New("database error")

			sqlMock.ExpectExec(`UPDATE gateways`).
				WithArgs(status, gatewayID).
				WillReturnError(dbError)

			err := repo.UpdateHealthStatus(ctx, gatewayID, status)
			gomega.Expect(err).Should(gomega.HaveOccurred())
			gomega.Expect(err.Error()).Should(gomega.ContainSubstring(dbError.Error()))
		})

		ginkgo.It("should return error when RowsAffected fails", func() {
			sqlMock.ExpectExec(`UPDATE gateways`).
				WithArgs(status, gatewayID).
				WillReturnResult(sqlmock.NewErrorResult(fmt.Errorf("rows affected error")))

			err := repo.UpdateHealthStatus(ctx, gatewayID, status)
			gomega.Expect(err).Should(gomega.HaveOccurred())
			gomega.Expect(err.Error()).Should(gomega.ContainSubstring("failed to retrieve rows affected"))
		})
	})
})
