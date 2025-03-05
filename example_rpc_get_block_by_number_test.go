package eth_test

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/streamingfast/eth-go/rpc"
)

func Example_rpc_getBlockByNumber() {
	client := rpc.NewClient(getRPCURL())
	blockRef := rpc.LatestBlock
	if os.Getenv("ETH_GO_RPC_BLOCK_NUMBER") != "" {
		var err error
		blockNumber, err := strconv.ParseUint(os.Getenv("ETH_GO_RPC_BLOCK_NUMBER"), 0, 64)
		if err != nil {
			panic(fmt.Errorf("parse custom block number info: %w", err))
		}

		blockRef = rpc.BlockNumber(blockNumber)
	}

	result, err := client.GetBlockByNumber(context.Background(), blockRef, rpc.WithGetBlockFullTransaction())
	if err != nil {
		panic(fmt.Errorf("get block by number %s: %w", blockRef, err))
	}

	bytes, err := rpc.MarshalJSONRPC(result)
	if err != nil {
		panic(fmt.Errorf("json marshal response: %w", err))
	}

	fmt.Println(string(bytes))
}
