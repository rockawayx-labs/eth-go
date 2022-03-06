package rpc

import (
	"strconv"
	"strings"

	"github.com/streamingfast/eth-go"
)

var LatestBlock = &BlockRef{tag: "latest"}
var PendingBlock = &BlockRef{tag: "pending"}
var EarliestBlock = &BlockRef{tag: "earliest"}

type BlockRef struct {
	tag   string
	value uint64
}

func BlockNumber(number uint64) *BlockRef {
	return &BlockRef{tag: "", value: number}
}

func (b *BlockRef) UnmarshalText(text []byte) error {
	lowerTextString := strings.ToLower(string(text))
	switch lowerTextString {
	case "latest":
		*b = *LatestBlock
	case "pending":
		*b = *PendingBlock
	case "earliest":
		*b = *EarliestBlock
	}

	var value eth.Uint64
	if err := value.UnmarshalText(text); err != nil {
		return err
	}

	b.tag = ""
	b.value = uint64(value)
	return nil
}

func (b BlockRef) MarshalJSONRPC() ([]byte, error) {
	if b.tag != "" {
		return MarshalJSONRPC(b.tag)
	}

	return MarshalJSONRPC(b.value)
}

func (b BlockRef) String() string {
	if b.tag != "" {
		return strings.ToUpper(string(b.tag[0])) + b.tag[1:]
	}

	return "#" + strconv.FormatUint(b.value, 10)
}

type LogEntry struct {
	Address          eth.Address `json:"address"`
	Topics           []eth.Hash  `json:"topics"`
	Data             eth.Hex     `json:"data"`
	BlockNumber      eth.Uint64  `json:"blockNumber"`
	TransactionHash  eth.Hash    `json:"transactionHash"`
	TransactionIndex eth.Uint64  `json:"transactionIndex"`
	BlockHash        eth.Hash    `json:"blockHash"`
	LogIndex         eth.Uint64  `json:"logIndex"`
	Removed          bool        `json:"removed"`
}

func (e *LogEntry) ToLog() (out eth.Log) {
	out.Address = e.Address
	if len(e.Topics) > 0 {
		out.Topics = make([][]byte, len(e.Topics))
		for i, topic := range e.Topics {
			out.Topics[i] = topic
		}
	}
	out.Data = e.Data
	out.BlockIndex = uint32(e.LogIndex)

	return
}

type TransactionReceipt struct {
	// TransactionHash is the hash of the transaction.
	TransactionHash eth.Hash `json:"transactionHash"`
	// TransactionIndex is the transactions index position in the block.
	TransactionIndex eth.Uint64 `json:"transactionIndex"`
	// BlockHash is the hash of the block where this transaction was in.
	BlockHash eth.Hash `json:"blockHash"`
	// BlockNumber is the block number where this transaction was in.
	BlockNumber eth.Uint64 `json:"blockNumber"`
	// From is the address of the sender.
	From eth.Address `json:"from"`
	// To is the address of the receiver, `null` when the transaction is a contract creation transaction.
	To *eth.Address `json:"to,omitempty"`
	// CumulativeGasUsed is the the total amount of gas used when this transaction was executed in the block.
	CumulativeGasUsed eth.Uint64 `json:"cumulativeGasUsed"`
	// GasUsed is the the amount of gas used by this specific transaction alone.
	GasUsed eth.Uint64 `json:"gasUsed"`
	// ContractAddress is the the contract address created, if the transaction was a contract creation, otherwise - null.
	ContractAddress *eth.Address `json:"contractAddress,omitempty"`
	// Logs is the Array of log objects, which this transaction generated.
	Logs []*LogEntry `json:"logs"`
	// LogsBloom is the Bloom filter for light clients to quickly retrieve related logs.
	LogsBloom eth.Hex `json:"logsBloom"`
}

type Block struct {
	Number           eth.Uint64    `json:"number"`
	Hash             eth.Hash      `json:"hash"`
	ParentHash       eth.Hash      `json:"parentHash"`
	Timestamp        eth.Timestamp `json:"timestamp"`
	StateRoot        eth.Hash      `json:"stateRoot"`
	TransactionsRoot eth.Hash      `json:"transactionsRoot"`
	ReceiptsRoot     eth.Hash      `json:"receiptsRoot"`
	MixHash          eth.Hash      `json:"mixHash"`
	GasLimit         eth.Uint64    `json:"gasLimit"`
	GasUsed          eth.Uint64    `json:"gasUsed"`
	// Left out for now because we need to actually pull in https://github.com/holiman/uint256
	// and get a eth.Uint256 wrapper in that is able to properly read those value from a JSON
	// string like "0x8792c6f47f70f".
	// Difficulty     *big.Int    `json:"difficulty"`
	// TotalDifficult *big.Int    `json:"totalDifficulty"`
	Miner         eth.Address `json:"miner"`
	Nonce         eth.Hex     `json:"nonce,omitempty"`
	LogsBloom     eth.Hex     `json:"logsBloom"`
	ExtraData     eth.Hex     `json:"extraData"`
	BaseFeePerGas eth.Uint64  `json:"baseFeePerGas,omitempty"`
	BlockSize     eth.Uint64  `json:"size,omitempty"`
	// Left out for now because this is a dynamic type. It's a []eth.Hash value when hydrating
	// params is sets to `false` and it's a `rpc.TransactionReceipt` when sets to `true`.
	// Transactions  interface{} `json:"transactions,omitempty"`
	UnclesSHA3 eth.Hash   `json:"sha3Uncles,omitempty"`
	Uncles     []eth.Hash `json:"uncles,omitempty"`
}
