package eth_test

import (
	"fmt"
	"os"

	"github.com/streamingfast/logging"
)

func init() {
	logging.InstantiateLoggers()
}

func getRPCURL() string {
	rpcURL := os.Getenv("ETH_GO_RPC_URL")
	if rpcURL != "" {
		return rpcURL
	}

	fmt.Println("Using 'http://localhost:8545' for testing purposes which is unlikely to work without extra setup")
	fmt.Println("There is free RPC provider(s) as long as you register like Infura, Alchemy or QuickNode")
	fmt.Println()
	fmt.Println("Once you have your personal RPC endpoint, specify environment variable ETH_GO_RPC_URL to provide one.")
	fmt.Println()

	return "http://localhost:8545"
}
