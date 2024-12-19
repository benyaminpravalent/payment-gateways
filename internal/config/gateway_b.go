package config

import (
	"log"
	"os"

	_ "github.com/lib/pq"
)

var (
	GatewayBUrl        string
	GatewayBApiKey     string
	GatewayBPrivateKey string
)

func init() {
	GatewayBUrl = os.Getenv("GATEWAY_B_URL")
	if GatewayBUrl == "" {
		log.Fatal("GATEWAY_B_URL environment variable is not set")
	}

	GatewayBApiKey = os.Getenv("GATEWAY_B_API_KEY")
	if GatewayBApiKey == "" {
		log.Fatal("GATEWAY_B_API_KEY environment variable is not set")
	}

	GatewayBPrivateKey = os.Getenv("GATEWAY_B_PRIVATE_KEY")
	if GatewayBPrivateKey == "" {
		log.Fatal("GATEWAY_B_PRIVATE_KEY environment variable is not set")
	}
}
