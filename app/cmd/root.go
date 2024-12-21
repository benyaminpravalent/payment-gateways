package cmd

import (
	"log"
	"payment-gateway/database"
	"payment-gateway/internal/kafka"
	"payment-gateway/internal/repositories"
	"payment-gateway/internal/services"

	"github.com/spf13/cobra"
)

// Declare services and repositories here
var (
	transactionRepository *repositories.TransactionRepository
	kafkaProducer         kafka.KafkaProducer
	transactionService    *services.TransactionService
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
	cobra.OnInitialize(initApp)
}

func initApp() {
	db := database.GetDB()

	kafkaProducer = kafka.NewKafkaProducer()

	transactionRepository = repositories.NewTransactionRepository(db)

	transactionService = services.NewTransactionService(transactionRepository, kafkaProducer)
}
