package config

import (
	"log"
	"os"

	_ "github.com/lib/pq"
)

var (
	GatewayCUrl        string
	GatewayCApiKey     string
	GatewayCPrivateKey string
)

func init() {
	GatewayCUrl = os.Getenv("GATEWAY_C_URL")
	if GatewayCUrl == "" {
		log.Fatal("GATEWAY_C_URL environment variable is not set")
	}

	GatewayCApiKey = os.Getenv("GATEWAY_C_API_KEY")
	if GatewayCApiKey == "" {
		log.Fatal("GATEWAY_C_API_KEY environment variable is not set")
	}

	GatewayCPrivateKey = os.Getenv("GATEWAY_C_PRIVATE_KEY")
	if GatewayCPrivateKey == "" {
		log.Fatal("GATEWAY_C_PRIVATE_KEY environment variable is not set")
	}
}
