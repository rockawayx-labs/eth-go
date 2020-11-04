package eth

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLogDecoder_ReadTypedTopic(t *testing.T) {

	data := `
{
	"address":"0x5C69bEe701ef814a2B6a3EDD4B1652CB9cc5aA6f",
	"topics":[
		"0x0d3648bd0f6ba80134a33ba9275ac585d9d315f0ad8355cddefde31afa28d0e9",
		"0x000000000000000000000000a0b86991c6218b36c1d19d4a2e9eb0ce3606eb48",
		"0x000000000000000000000000f1290473e210b2108a85237fbcd7b6eb42cc654f"
	],
	"data":"0x000000000000000000000000fc2890ffb3069a1a9d3f7b11c7775a1a1ee721c00000000000000000000000000000000000000000000000000000000000002f4d"}`

	var l struct {
		Address string   `json:"address"`
		Topics  []string `json:"topics"`
		Data    string   `json:"data"`
	}
	err := json.Unmarshal([]byte(data), &l)
	require.NoError(t, err)

	var topics [][]byte
	for _, t := range l.Topics {
		topics = append(topics, MustDecodeString(t))
	}

	log := &Log{
		Address:    MustNewAddress(l.Address),
		Topics:     topics,
		Data:       MustNewAddress(l.Data),
		Index:      0,
		BlockIndex: 0,
	}

	decoder := NewLogDecoder(log)
	_, _ = decoder.ReadTopic() //skips topic 0 (kessac(signature)

	t0Addr, err := decoder.ReadTypedTopic("address")
	require.NoError(t, err)
}
