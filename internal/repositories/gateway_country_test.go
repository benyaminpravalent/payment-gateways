package repositories

import (
	"context"
	"errors"

	"payment-gateway/models"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("GatewayCountryRepository", func() {
	var (
		mockDB       *sqlx.DB
		sqlMock      sqlmock.Sqlmock
		repo         *GatewayCountryRepository
		ctx          context.Context
		countryID    int
		expectedData *models.GatewayDetail
	)

	ginkgo.BeforeEach(func() {
		var err error
		sqlDB, mock, err := sqlmock.New()
		gomega.Expect(err).ShouldNot(gomega.HaveOccurred())

		mockDB = sqlx.NewDb(sqlDB, "sqlmock")
		sqlMock = mock
		repo = NewGatewayCountryRepository(mockDB)

		ctx = context.Background()
		countryID = 1
		expectedData = &models.GatewayDetail{
			ID:                  1,
			Name:                "Test Gateway",
			DataFormatSupported: "JSON",
			HealthStatus:        "healthy",
			Priority:            1,
			CountryID:           countryID,
			Currency:            "USD",
		}
	})

	ginkgo.AfterEach(func() {
		err := sqlMock.ExpectationsWereMet()
		gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
	})

	ginkgo.Describe("GetHealthyGatewayByCountryID", func() {
		ginkgo.It("should successfully fetch healthy gateway detail", func() {
			rows := sqlmock.NewRows([]string{
				"id", "name", "data_format_supported", "health_status", "last_checked_at",
				"created_at", "updated_at", "priority", "country_id", "currency",
			}).AddRow(
				expectedData.ID, expectedData.Name, expectedData.DataFormatSupported,
				expectedData.HealthStatus, expectedData.LastCheckedAt, expectedData.CreatedAt,
				expectedData.UpdatedAt, expectedData.Priority, expectedData.CountryID,
				expectedData.Currency,
			)

			sqlMock.ExpectQuery(`SELECT g.id, g.name, g.data_format_supported, g.health_status, g.last_checked_at, g.created_at, g.updated_at, gc.priority, gc.country_id, c.currency`).
				WithArgs(countryID).
				WillReturnRows(rows)

			result, err := repo.GetHealthyGatewayByCountryID(ctx, countryID)
			gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
			gomega.Expect(result).Should(gomega.Equal(expectedData))
		})

		ginkgo.It("should return error when no rows match the query", func() {
			sqlMock.ExpectQuery(`SELECT g.id, g.name, g.data_format_supported, g.health_status, g.last_checked_at, g.created_at, g.updated_at, gc.priority, gc.country_id, c.currency`).
				WithArgs(countryID).
				WillReturnRows(sqlmock.NewRows(nil)) // Empty result

			result, err := repo.GetHealthyGatewayByCountryID(ctx, countryID)
			gomega.Expect(err).Should(gomega.HaveOccurred())
			gomega.Expect(err.Error()).Should(gomega.ContainSubstring("sql: no rows in result set"))
			gomega.Expect(result).Should(gomega.BeNil())
		})

		ginkgo.It("should return error when database query fails", func() {
			dbError := errors.New("database error")

			sqlMock.ExpectQuery(`SELECT g.id, g.name, g.data_format_supported, g.health_status, g.last_checked_at, g.created_at, g.updated_at, gc.priority, gc.country_id, c.currency`).
				WithArgs(countryID).
				WillReturnError(dbError)

			result, err := repo.GetHealthyGatewayByCountryID(ctx, countryID)
			gomega.Expect(err).Should(gomega.HaveOccurred())
			gomega.Expect(err.Error()).Should(gomega.ContainSubstring(dbError.Error()))
			gomega.Expect(result).Should(gomega.BeNil())
		})
	})
})
