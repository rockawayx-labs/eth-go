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

// GetTokenInfo returns an *eth.Token object if it can.
// It can be called at a specific block number (with an archive node), or "latest" if atBlockNum is 0
// It can validate that a specific block hash exists that chain by setting it as ensureBlockHashIsInChain != "".
func (c *Client) GetTokenInfo(tokenAddr eth.Address, atBlockNum uint64) (token *eth.Token, err error) {

	var atExpression interface{}
	if atBlockNum == 0 {
		atExpression = "latest"
	} else {
		atExpression = atBlockNum
	}

	requests := []*RPCRequest{
		{
			Params: []interface{}{CallParams{To: tokenAddr, Data: decimalsCallData}, atExpression},
			Method: "eth_call",
		},
		{
			Params: []interface{}{CallParams{To: tokenAddr, Data: nameCallData}, atExpression},
			Method: "eth_call",
		},
		{
			Params: []interface{}{CallParams{To: tokenAddr, Data: symbolCallData}, atExpression},
			Method: "eth_call",
		},
		{
			Params: []interface{}{CallParams{To: tokenAddr, Data: totalSupplyCallData}, atExpression},
			Method: "eth_call",
		},
	}

	//  possible way to check if we are in the same "fork"
	//	if ensureBlockHashIsInChain != "" {
	//		requests = append(requests, &RPCRequest{
	//			Method: "eth_getUncleCountByBlockHash",
	//			Params: []interface{}{ensureBlockHashIsInChain}},
	//		)
	//	}

	results, err := c.DoRequests(requests)
	if err != nil {
		return nil, err
	}

	for _, result := range results {
		if result.err != nil {
			return nil, result.err
		}
	}

	//  possible way to check if we are in the same "fork"
	//	if ensureBlockHashIsInChain != "" {
	//		if results[len(results)-1].content == "" {
	//			return nil, fmt.Errorf("blockHash %s requested for lookup but not found on chain", ensureBlockHashIsInChain)
	//		}
	//	}

	decimalsResult := results[0].content
	nameResult := results[1].content
	symbolResult := results[2].content
	totalSupplyResult := results[3].content

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
			return nil, fmt.Errorf("decode decimals %q: %w", decimalsResult, err)
		}

		decimals = out[0]
	}

	if !emptyName {
		out, err := nameMethodDef.DecodeOutput(eth.MustNewHex(nameResult))
		if err != nil {
			return nil, fmt.Errorf("decode name %q: %w", nameResult, err)
		}

		name = out[0]
	}

	if !emptySymbol {
		out, err := symbolMethodDef.DecodeOutput(eth.MustNewHex(symbolResult))
		if err != nil {
			return nil, fmt.Errorf("decode symbol %q: %w", symbolResult, err)
		}

		symbol = out[0]
	}

	if !emptyTotalSupply {
		out, err := totalSupplyMethodDef.DecodeOutput(eth.MustNewHex(totalSupplyResult))
		if err != nil {
			return nil, fmt.Errorf("decode total supply %q: %w", totalSupplyResult, err)
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

var isEmptyRegex = regexp.MustCompile("^0x$")

func isEmptyResult(result string) bool {
	return isEmptyRegex.MatchString(result)
}
