package services

import (
	"context"
	"errors"
	mocks "payment-gateway/mocks/repositories"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("GatewayService", func() {
	var (
		mockRepo       *mocks.MockGatewayRepository
		gatewayService *GatewayService
	)

	ginkgo.BeforeEach(func() {
		mockRepo = new(mocks.MockGatewayRepository)
		gatewayService = NewGatewayService(mockRepo)
	})

	ginkgo.Describe("UpdateGatewayHealthStatusByID", func() {
		ginkgo.It("should successfully update health status", func() {
			gatewayID := 1
			status := "healthy"

			mockRepo.On("UpdateHealthStatus", context.Background(), gatewayID, status).Return(nil)

			err := gatewayService.UpdateGatewayHealthStatusByID(context.Background(), gatewayID, status)

			gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
			mockRepo.AssertCalled(ginkgo.GinkgoT(), "UpdateHealthStatus", context.Background(), gatewayID, status)
		})

		ginkgo.It("should return error when no rows are affected", func() {
			gatewayID := 2
			status := "unhealthy"

			mockRepo.On("UpdateHealthStatus", context.Background(), gatewayID, status).Return(errors.New("no gateway found with id"))

			err := gatewayService.UpdateGatewayHealthStatusByID(context.Background(), gatewayID, status)

			gomega.Expect(err).Should(gomega.HaveOccurred())
			gomega.Expect(err.Error()).Should(gomega.ContainSubstring("no gateway found with id"))
			mockRepo.AssertCalled(ginkgo.GinkgoT(), "UpdateHealthStatus", context.Background(), gatewayID, status)
		})

		ginkgo.It("should return error when database query fails", func() {
			gatewayID := 3
			status := "healthy"
			dbError := errors.New("database error")

			mockRepo.On("UpdateHealthStatus", context.Background(), gatewayID, status).Return(dbError)

			err := gatewayService.UpdateGatewayHealthStatusByID(context.Background(), gatewayID, status)

			gomega.Expect(err).Should(gomega.HaveOccurred())
			gomega.Expect(err.Error()).Should(gomega.ContainSubstring(dbError.Error()))
			mockRepo.AssertCalled(ginkgo.GinkgoT(), "UpdateHealthStatus", context.Background(), gatewayID, status)
		})
	})
})
