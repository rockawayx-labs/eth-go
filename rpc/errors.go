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
	"strings"
)

type ErrResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func (e *ErrResponse) Error() string {
	return fmt.Sprintf("rpc error (code %d): %s", e.Code, e.Message)
}

// Try to check if the call was reverted. The JSON-RPC response for reverts is
// not standardized, so we have ad-hoc checks for each of Geth, Parity and
// Ganache.
// These come from https://github.com/graphprotocol/graph-node/blob/master/chain/ethereum/src/ethereum_adapter.rs

var GETH_DETERMINISTIC_ERRORS = []string{
	"execution reverted",
	"invalid jump destination",
	"invalid opcode",
	"stack limit reached 1024",
	"stack underflow (",
	"gas uint64 overflow",
}

const PARITY_BAD_INSTRUCTION_FE = "Bad instruction fe"
const PARITY_BAD_INSTRUCTION_FD = "Bad instruction fd"
const PARITY_BAD_JUMP_PREFIX = "Bad jump"
const PARITY_STACK_LIMIT_PREFIX = "Out of stack"
const PARITY_VM_EXECUTION_ERROR = -32015
const PARITY_REVERT_PREFIX = "Reverted 0x"

const GANACHE_VM_EXECUTION_ERROR = -32000
const GANACHE_REVERT_MESSAGE = "VM Exception while processing transaction: revert"

func IsDeterministicError(err *ErrResponse) bool {
	return IsGethDeterministicError(err) ||
		IsParityDeterministicError(err) ||
		IsGanacheDeterministicError(err)
}

func IsGethDeterministicError(err *ErrResponse) bool {
	msg := strings.ToLower(err.Message)
	for _, e := range GETH_DETERMINISTIC_ERRORS {
		if strings.Contains(msg, e) {
			return true
		}
	}
	return false
}

func IsParityDeterministicError(err *ErrResponse) bool {
	if err.Code == PARITY_VM_EXECUTION_ERROR {
		return true
	}

	if err.Message == PARITY_BAD_INSTRUCTION_FD ||
		err.Message == PARITY_BAD_INSTRUCTION_FE {
		return true
	}

	if strings.HasPrefix(err.Message, PARITY_BAD_JUMP_PREFIX) ||
		strings.HasPrefix(err.Message, PARITY_STACK_LIMIT_PREFIX) {
		return true
	}

	return false
}

func IsGanacheDeterministicError(err *ErrResponse) bool {
	return err.Code == GANACHE_VM_EXECUTION_ERROR && strings.HasPrefix(err.Message, GANACHE_REVERT_MESSAGE)
}
