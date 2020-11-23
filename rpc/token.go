package rpc

import (
	"fmt"
	"math/big"

	"github.com/dfuse-io/eth-go"
)

var decimalsCallData = eth.MustNewMethodDef("decimals()").NewCall().MustEncode()
var nameCallData = eth.MustNewMethodDef("name()").NewCall().MustEncode()
var symbolCallData = eth.MustNewMethodDef("symbol()").NewCall().MustEncode()

var b0 = new(big.Int)

func (c *Client) GetTokenInfo(tokenAddr eth.Address) (*eth.Token, error) {
	decimalsResult, err := c.Call(CallParams{To: tokenAddr, Data: decimalsCallData})
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve decimals for token %q: %w", tokenAddr, err)
	}

	nameResult, err := c.Call(CallParams{To: tokenAddr, Data: nameCallData})
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve name for token %q: %w", tokenAddr, err)
	}

	symbolResult, err := c.Call(CallParams{To: tokenAddr, Data: symbolCallData})
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve symbol for token %q: %w", tokenAddr, err)
	}

	if decimalsResult == "0x" && nameResult == "0x" && symbolResult == "0x" {
		return nil, fmt.Errorf("not implementing one of ERC20 contract's method 'name()', or 'symbol()' or 'decimals()'")
	}

	var decimals interface{} = b0
	var symbol interface{} = ""
	var name interface{} = ""
	var dec *eth.Decoder

	if decimalsResult != "0x" {
		dec, err = eth.NewDecoderFromString(decimalsResult)
		if err != nil {
			return nil, fmt.Errorf("new decimals decoder %q for %q: %w", decimalsResult, tokenAddr, err)
		}

		decimals, err = dec.Read("uint256")
		if err != nil {
			return nil, fmt.Errorf("decode decimals %q: %w", decimalsResult, err)
		}
	}

	if symbolResult != "0x" {
		dec, err = eth.NewDecoderFromString(symbolResult)
		if err != nil {
			return nil, fmt.Errorf("new symbol decoder: %w", err)
		}

		_, err = dec.Read("uint256") // reading the initial string offset... we can ignore it
		if err != nil {
			return nil, fmt.Errorf("decoding symbol offset %s: %w", symbolResult, err)
		}

		symbol, err = dec.Read("string")
		if err != nil {
			return nil, fmt.Errorf("decoding symbol %s: %w", symbolResult, err)
		}
	}

	if nameResult != "0x" {
		dec, err = eth.NewDecoderFromString(nameResult)
		if err != nil {
			return nil, fmt.Errorf("new decoder: %w", err)
		}

		_, err = dec.Read("uint256") // reading the initial string offset... we can ignore it
		if err != nil {
			return nil, fmt.Errorf("decoding name offset %s: %w", symbolResult, err)
		}

		name, err = dec.Read("string")
		if err != nil {
			return nil, fmt.Errorf("decoding name %q: %w", nameResult, err)
		}
	}

	return &eth.Token{
		Address:  tokenAddr,
		Name:     name.(string),
		Symbol:   symbol.(string),
		Decimals: uint(decimals.(*big.Int).Uint64()),
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
