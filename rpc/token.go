// Copyright 2021 dfuse Platform Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package rpc

import (
	"fmt"
	"math/big"
	"regexp"

	"github.com/streamingfast/eth-go"
)

var decimalsMethodDef = eth.MustNewMethodDef("decimals() (uint256)")
var nameMethodDef = eth.MustNewMethodDef("name() (string)")
var symbolMethodDef = eth.MustNewMethodDef("symbol() (string)")
var totalSupplyMethodDef = eth.MustNewMethodDef("totalSupply() (uint256)")

var decimalsCallData = decimalsMethodDef.NewCall().MustEncode()
var nameCallData = nameMethodDef.NewCall().MustEncode()
var symbolCallData = symbolMethodDef.NewCall().MustEncode()
var totalSupplyCallData = totalSupplyMethodDef.NewCall().MustEncode()

var b0 = new(big.Int)

func (c *Client) GetTokenInfo(tokenAddr eth.Address) (token *eth.Token, headBlockNum uint64, err error) {

	results, err := c.DoRequests([]*RPCRequest{
		{
			Params: []interface{}{},
			Method: "eth_blockNumber",
		},
		{
			Params: []interface{}{CallParams{To: tokenAddr, Data: decimalsCallData}, "latest"},
			Method: "eth_call",
		},
		{
			Params: []interface{}{CallParams{To: tokenAddr, Data: nameCallData}, "latest"},
			Method: "eth_call",
		},
		{
			Params: []interface{}{CallParams{To: tokenAddr, Data: symbolCallData}, "latest"},
			Method: "eth_call",
		},
		{
			Params: []interface{}{CallParams{To: tokenAddr, Data: totalSupplyCallData}, "latest"},
			Method: "eth_call",
		},
	},
	)
	if err != nil {
		return nil, 0, err
	}

	for _, result := range results {
		if result.err != nil {
			return nil, 0, result.err
		}
	}

	headBlockNum = hex2uint64(results[0].content)
	decimalsResult := results[1].content
	nameResult := results[2].content
	symbolResult := results[3].content
	totalSupplyResult := results[4].content

	emptyDecimal := isEmptyResult(decimalsResult)
	emptyName := isEmptyResult(nameResult)
	emptySymbol := isEmptyResult(symbolResult)
	emptyTotalSupply := isEmptyResult(totalSupplyResult)

	var decimals interface{} = b0
	var symbol interface{} = ""
	var name interface{} = ""
	var totalSupply interface{} = b0

	if !emptyDecimal {
		out, err := decimalsMethodDef.DecodeOutput(eth.MustNewHex(decimalsResult))
		if err != nil {
			return nil, headBlockNum, fmt.Errorf("decode decimals %q: %w", decimalsResult, err)
		}

		decimals = out[0]
	}

	if !emptyName {
		out, err := nameMethodDef.DecodeOutput(eth.MustNewHex(nameResult))
		if err != nil {
			return nil, headBlockNum, fmt.Errorf("decode name %q: %w", nameResult, err)
		}

		name = out[0]
	}

	if !emptySymbol {
		out, err := symbolMethodDef.DecodeOutput(eth.MustNewHex(symbolResult))
		if err != nil {
			return nil, headBlockNum, fmt.Errorf("decode symbol %q: %w", symbolResult, err)
		}

		symbol = out[0]
	}

	if !emptyTotalSupply {
		out, err := totalSupplyMethodDef.DecodeOutput(eth.MustNewHex(totalSupplyResult))
		if err != nil {
			return nil, headBlockNum, fmt.Errorf("decode total supply %q: %w", totalSupplyResult, err)
		}

		totalSupply = out[0]
	}

	return &eth.Token{
		Address:            tokenAddr,
		Name:               name.(string),
		Symbol:             symbol.(string),
		Decimals:           uint(decimals.(*big.Int).Uint64()),
		TotalSupply:        totalSupply.(*big.Int),
		IsEmptyName:        emptyName,
		IsEmptyDecimal:     emptyDecimal,
		IsEmptySymbol:      emptySymbol,
		IsEmptyTotalSupply: emptyTotalSupply,
	}, headBlockNum, nil
}

func methodSignatureBytes(def *eth.MethodDef) []byte {
	encoder := eth.NewEncoder()
	err := encoder.Write("method", def.Signature())
	if err != nil {
		return nil
	}

	return encoder.Buffer()
}

var isEmptyRegex = regexp.MustCompile("^0x$")

func isEmptyResult(result string) bool {
	return isEmptyRegex.MatchString(result)
}
