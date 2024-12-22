package main

import (
	"log"
	"payment-gateway/app/cmd"
	"payment-gateway/database"
	"payment-gateway/internal/config"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(".env is not available")
	}
	database.InitDB()

	config.InitGatewayA()
	config.InitGatewayB()
	config.InitGatewayC()

	cmd.Execute()
}
