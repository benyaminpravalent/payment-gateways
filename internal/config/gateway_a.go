package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var (
	GatewayAUrl        string
	GatewayAApiKey     string
	GatewayAPrivateKey string
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(".env is not available")
	}

	GatewayAUrl = os.Getenv("GATEWAY_A_URL")
	if GatewayAUrl == "" {
		log.Fatal("GATEWAY_A_URL environment variable is not set")
	}

	GatewayAApiKey = os.Getenv("GATEWAY_A_API_KEY")
	if GatewayAApiKey == "" {
		log.Fatal("GATEWAY_A_API_KEY environment variable is not set")
	}

	GatewayAPrivateKey = os.Getenv("GATEWAY_A_PRIVATE_KEY")
	if GatewayAPrivateKey == "" {
		log.Fatal("GATEWAY_A_PRIVATE_KEY environment variable is not set")
	}
}
