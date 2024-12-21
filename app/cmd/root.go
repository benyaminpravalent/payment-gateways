package cmd

import (
	"log"
	"payment-gateway/database"
	"payment-gateway/internal/client"
	"payment-gateway/internal/kafka"
	"payment-gateway/internal/repositories"
	"payment-gateway/internal/services"

	"github.com/spf13/cobra"
)

// Declare services and repositories here
var (
	TransactionRepository *repositories.TransactionRepository
	KafkaProducer         kafka.KafkaProducer
	TransactionService    *services.TransactionService
	SendTransactionClient *client.TransactionClient
	GatewayCountryRepo    *repositories.GatewayCountryRepository
	GatewayRepo           *repositories.GatewayRepository
)

var (
	envFilePath string
	rootCmd     = &cobra.Command{
		Use:   "payment-gateway",
		Short: "Backend service for PaymentGateway",
	}
)

func Execute() {
	rootCmd.PersistentFlags().StringVarP(&envFilePath, "env", "e", ".env", ".env file to read from")
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func init() {
	cobra.OnInitialize(initApp, initConsumer)
}

func initApp() {
	db := database.GetDB()

	KafkaProducer = kafka.NewKafkaProducer()

	SendTransactionClient = client.NewTransactionClient()

	GatewayCountryRepo = repositories.NewGatewayCountryRepository(db)
	TransactionRepository = repositories.NewTransactionRepository(db)
	GatewayRepo = repositories.NewGatewayRepository(db)

	TransactionService = services.NewTransactionService(TransactionRepository, KafkaProducer)
}

func initConsumer() {
	go kafka.InitializeKafkaConsumer()
}
