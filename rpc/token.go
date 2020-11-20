package rpc

import (
	"fmt"
	"math/big"

	"github.com/dfuse-io/eth-go"
)

var decimalsCallData = eth.MustNewMethodDef("decimals()").NewCall().MustEncode()
var nameCallData = eth.MustNewMethodDef("name()").NewCall().MustEncode()
var symbolCallData = eth.MustNewMethodDef("symbol()").NewCall().MustEncode()

func (c *Client) GetTokenInfo(tokenAddr eth.Address) (*eth.Token, error) {
	decimals, err := c.Call(CallParams{To: tokenAddr, Data: decimalsCallData})
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve decimals for token %q: %w", tokenAddr, err)
	}

	name, err := c.Call(CallParams{To: tokenAddr, Data: nameCallData})
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve name for token %q: %w", tokenAddr, err)
	}

	symbol, err := c.Call(CallParams{To: tokenAddr, Data: symbolCallData})
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve symbol for token %q: %w", tokenAddr, err)
	}

	dec, err := eth.NewDecoderFromString(decimals)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve name for symbol %q: %w", tokenAddr, err)
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

func methodSignatureBytes(def *eth.MethodDef) []byte {
	encoder := eth.NewEncoder()
	err := encoder.Write("method", def.Signature())
	if err != nil {
		return nil
	}

	return encoder.Buffer()
}
