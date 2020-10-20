package rpc

import (
	"fmt"
	"math/big"

	"github.com/dfuse-io/eth-go"
)

func (c *Client) GetTokenInfo(tokenAddr eth.Address) (*eth.Token, error) {
	rpcTokenAddr := tokenAddr.Pretty()
	decimals, err := c.Call(map[string]string{
		"data": encodeMethod("decimals()"),
		"to":   rpcTokenAddr,
	})
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve decimals for token %q: %w", rpcTokenAddr, err)
	}

	name, err := c.Call(map[string]string{
		"data": encodeMethod("name()"),
		"to":   rpcTokenAddr,
	})
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve name for token %q: %w", rpcTokenAddr, err)
	}

	symbol, err := c.Call(map[string]string{
		"data": encodeMethod("symbol()"),
		"to":   rpcTokenAddr,
	})
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve name for symbol %q: %w", rpcTokenAddr, err)
	}

	dec, err := eth.NewDecoderFromString(decimals)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve name for symbol %q: %w", rpcTokenAddr, err)
	}
	decodedDecimals, err := dec.Read("uint256")
	if err != nil {
		return nil, fmt.Errorf("unable to decoding decimals %q: %w", decimals, err)

	}

	dec, err = eth.NewDecoderFromString(symbol)
	if err != nil {
		return nil, fmt.Errorf("unable to new decoder: %w", err)
	}
	_, err = dec.Read("uint256") // reading the initial string offset... we can ignore it
	if err != nil {
		return nil, fmt.Errorf("decoding symbol offset %s: %w", symbol, err)
	}

	decodedSymbol, err := dec.Read("string")
	if err != nil {
		return nil, fmt.Errorf("decoding symbol %s: %w", symbol, err)
	}

	dec, err = eth.NewDecoderFromString(name)
	if err != nil {
		return nil, fmt.Errorf("new decoder: %w", err)
	}

	_, err = dec.Read("uint256") // reading the initial string offset... we can ignore it
	if err != nil {
		return nil, fmt.Errorf("decoding name offset %s: %w", symbol, err)
	}

	decodedName, err := dec.Read("string")
	if err != nil {
		return nil, fmt.Errorf("decoding name %s: %w", name, err)
	}

	return &eth.Token{
		Name:     decodedName.(string),
		Symbol:   decodedSymbol.(string),
		Address:  tokenAddr,
		Decimals: uint(decodedDecimals.(*big.Int).Uint64()),
	}, nil
}

func encodeMethod(methodStr string) string {
	encoder := eth.NewEncoder()
	err := encoder.Write("method", methodStr)
	if err != nil {
		return ""
	}
	return "0x" + encoder.String()
}
