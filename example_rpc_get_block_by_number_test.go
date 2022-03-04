package eth_test

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/streamingfast/eth-go/rpc"
)

func ExampleRPC_GetBlockByNumber() {
	client := rpc.NewClient(getRPCURL())
	blockNumber := uint64(10000000)
	if os.Getenv("ETH_GO_RPC_BLOCK_NUMBER") != "" {
		var err error
		blockNumber, err = strconv.ParseUint(os.Getenv("ETH_GO_RPC_BLOCK_NUMBER"), 0, 64)
		if err != nil {
			panic(fmt.Errorf("parse custom block number info: %w", err))
		}
	}

	result, err := client.GetBlockByNumber(context.Background(), blockNumber)
	if err != nil {
		panic(fmt.Errorf("get block by number %d: %w", blockNumber, err))
	}

	bytes, err := json.Marshal(result)
	if err != nil {
		panic(fmt.Errorf("json marshal response: %w", err))
	}

	fmt.Println(string(bytes))
}
