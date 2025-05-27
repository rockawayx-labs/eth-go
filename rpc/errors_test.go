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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsDeterministicError(t *testing.T) {
	tests := []struct {
		name     string
		err      *ErrResponse
		expected bool
	}{
		// Generic JSON-RPC errors
		{
			name:     "invalid request error",
			err:      &ErrResponse{Code: JSON_RPC_INVALID_REQUEST_ERROR, Message: "Invalid request"},
			expected: true,
		},
		{
			name:     "invalid argument error",
			err:      &ErrResponse{Code: JSON_RPC_INVALID_ARGUMENT_ERROR, Message: "Invalid argument"},
			expected: true,
		},
		// Geth deterministic errors
		{
			name:     "geth revert error",
			err:      &ErrResponse{Code: -32000, Message: "execution reverted"},
			expected: true,
		},
		{
			name:     "geth invalid jump destination",
			err:      &ErrResponse{Code: -32000, Message: "invalid jump destination"},
			expected: true,
		},
		{
			name:     "geth invalid opcode",
			err:      &ErrResponse{Code: -32000, Message: "invalid opcode"},
			expected: true,
		},
		{
			name:     "geth stack limit reached",
			err:      &ErrResponse{Code: -32000, Message: "stack limit reached 1024"},
			expected: true,
		},
		{
			name:     "geth stack underflow",
			err:      &ErrResponse{Code: -32000, Message: "stack underflow (0 <-> 1)"},
			expected: true,
		},
		{
			name:     "geth gas uint64 overflow",
			err:      &ErrResponse{Code: -32000, Message: "gas uint64 overflow"},
			expected: true,
		},
		{
			name:     "geth out of gas",
			err:      &ErrResponse{Code: -32000, Message: "out of gas"},
			expected: true,
		},
		// Parity deterministic errors
		{
			name:     "parity vm execution error",
			err:      &ErrResponse{Code: PARITY_VM_EXECUTION_ERROR, Message: "VM execution error"},
			expected: true,
		},
		{
			name:     "parity bad instruction fe",
			err:      &ErrResponse{Code: -32000, Message: PARITY_BAD_INSTRUCTION_FE},
			expected: true,
		},
		{
			name:     "parity bad instruction fd",
			err:      &ErrResponse{Code: -32000, Message: PARITY_BAD_INSTRUCTION_FD},
			expected: true,
		},
		{
			name:     "parity bad jump",
			err:      &ErrResponse{Code: -32000, Message: "Bad jump to invalid destination"},
			expected: true,
		},
		{
			name:     "parity stack limit",
			err:      &ErrResponse{Code: -32000, Message: "Out of stack space"},
			expected: true,
		},
		{
			name:     "parity out of gas",
			err:      &ErrResponse{Code: -32000, Message: PARITY_OUT_OF_GAS},
			expected: true,
		},
		{
			name:     "parity out of bounds",
			err:      &ErrResponse{Code: -32000, Message: PARITY_OUT_OF_BOUND},
			expected: true,
		},
		{
			name:     "parity revert with data",
			err:      &ErrResponse{Code: -32000, Message: "Reverted 0x08c379a00000000000000000000000000000000000000000000000000000000000000020"},
			expected: true,
		},
		// Ganache deterministic errors
		{
			name:     "ganache vm exception revert",
			err:      &ErrResponse{Code: GANACHE_VM_EXECUTION_ERROR, Message: GANACHE_REVERT_MESSAGE},
			expected: true,
		},
		{
			name:     "ganache vm exception with reason",
			err:      &ErrResponse{Code: GANACHE_VM_EXECUTION_ERROR, Message: "VM Exception while processing transaction: revert Custom error message"},
			expected: true,
		},
		// Non-deterministic errors
		{
			name:     "network timeout",
			err:      &ErrResponse{Code: -32603, Message: "network timeout"},
			expected: false,
		},
		{
			name:     "internal error",
			err:      &ErrResponse{Code: -32603, Message: "internal error"},
			expected: false,
		},
		{
			name:     "connection refused",
			err:      &ErrResponse{Code: -32000, Message: "connection refused"},
			expected: false,
		},
		{
			name:     "unknown error",
			err:      &ErrResponse{Code: -1, Message: "unknown error"},
			expected: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := IsDeterministicError(test.err)
			assert.Equal(t, test.expected, result, "Expected IsDeterministicError to return %v for error: %+v", test.expected, test.err)
		})
	}
}

func TestIsGenericDeterministicError(t *testing.T) {
	tests := []struct {
		name     string
		err      *ErrResponse
		expected bool
	}{
		{
			name:     "invalid request error",
			err:      &ErrResponse{Code: JSON_RPC_INVALID_REQUEST_ERROR, Message: "Invalid request"},
			expected: true,
		},
		{
			name:     "invalid argument error",
			err:      &ErrResponse{Code: JSON_RPC_INVALID_ARGUMENT_ERROR, Message: "Invalid argument"},
			expected: true,
		},
		{
			name:     "other error code",
			err:      &ErrResponse{Code: -32000, Message: "Other error"},
			expected: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := IsGenericDeterministicError(test.err)
			assert.Equal(t, test.expected, result)
		})
	}
}

func TestIsGethDeterministicError(t *testing.T) {
	tests := []struct {
		name     string
		err      *ErrResponse
		expected bool
	}{
		{
			name:     "revert error lowercase",
			err:      &ErrResponse{Code: -32000, Message: "execution reverted"},
			expected: true,
		},
		{
			name:     "revert error uppercase",
			err:      &ErrResponse{Code: -32000, Message: "EXECUTION REVERTED"},
			expected: true,
		},
		{
			name:     "revert error mixed case",
			err:      &ErrResponse{Code: -32000, Message: "Execution Reverted"},
			expected: true,
		},
		{
			name:     "invalid jump destination",
			err:      &ErrResponse{Code: -32000, Message: "invalid jump destination"},
			expected: true,
		},
		{
			name:     "invalid opcode",
			err:      &ErrResponse{Code: -32000, Message: "invalid opcode"},
			expected: true,
		},
		{
			name:     "stack limit reached",
			err:      &ErrResponse{Code: -32000, Message: "stack limit reached 1024"},
			expected: true,
		},
		{
			name:     "stack underflow",
			err:      &ErrResponse{Code: -32000, Message: "stack underflow (0 <-> 1)"},
			expected: true,
		},
		{
			name:     "gas uint64 overflow",
			err:      &ErrResponse{Code: -32000, Message: "gas uint64 overflow"},
			expected: true,
		},
		{
			name:     "out of gas",
			err:      &ErrResponse{Code: -32000, Message: "out of gas"},
			expected: true,
		},
		{
			name:     "non-geth error",
			err:      &ErrResponse{Code: -32000, Message: "network timeout"},
			expected: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := IsGethDeterministicError(test.err)
			assert.Equal(t, test.expected, result)
		})
	}
}

func TestIsParityDeterministicError(t *testing.T) {
	tests := []struct {
		name     string
		err      *ErrResponse
		expected bool
	}{
		{
			name:     "parity vm execution error code",
			err:      &ErrResponse{Code: PARITY_VM_EXECUTION_ERROR, Message: "Any message"},
			expected: true,
		},
		{
			name:     "parity out of gas message",
			err:      &ErrResponse{Code: -32000, Message: PARITY_OUT_OF_GAS},
			expected: true,
		},
		{
			name:     "parity out of bounds with correct code",
			err:      &ErrResponse{Code: -32000, Message: PARITY_OUT_OF_BOUND},
			expected: true,
		},
		{
			name:     "parity out of bounds with wrong code",
			err:      &ErrResponse{Code: -32001, Message: PARITY_OUT_OF_BOUND},
			expected: false,
		},
		{
			name:     "parity bad instruction fd",
			err:      &ErrResponse{Code: -32000, Message: PARITY_BAD_INSTRUCTION_FD},
			expected: true,
		},
		{
			name:     "parity bad instruction fe",
			err:      &ErrResponse{Code: -32000, Message: PARITY_BAD_INSTRUCTION_FE},
			expected: true,
		},
		{
			name:     "parity bad jump prefix",
			err:      &ErrResponse{Code: -32000, Message: "Bad jump to invalid destination"},
			expected: true,
		},
		{
			name:     "parity stack limit prefix",
			err:      &ErrResponse{Code: -32000, Message: "Out of stack space available"},
			expected: true,
		},
		{
			name:     "non-parity error",
			err:      &ErrResponse{Code: -32603, Message: "internal error"},
			expected: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := IsParityDeterministicError(test.err)
			assert.Equal(t, test.expected, result)
		})
	}
}

func TestIsGanacheDeterministicError(t *testing.T) {
	tests := []struct {
		name     string
		err      *ErrResponse
		expected bool
	}{
		{
			name:     "ganache vm exception revert",
			err:      &ErrResponse{Code: GANACHE_VM_EXECUTION_ERROR, Message: GANACHE_REVERT_MESSAGE},
			expected: true,
		},
		{
			name:     "ganache vm exception with custom message",
			err:      &ErrResponse{Code: GANACHE_VM_EXECUTION_ERROR, Message: "VM Exception while processing transaction: revert Custom error"},
			expected: true,
		},
		{
			name:     "ganache wrong code",
			err:      &ErrResponse{Code: -32001, Message: GANACHE_REVERT_MESSAGE},
			expected: false,
		},
		{
			name:     "ganache wrong message prefix",
			err:      &ErrResponse{Code: GANACHE_VM_EXECUTION_ERROR, Message: "Different error message"},
			expected: false,
		},
		{
			name:     "non-ganache error",
			err:      &ErrResponse{Code: -32603, Message: "internal error"},
			expected: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := IsGanacheDeterministicError(test.err)
			assert.Equal(t, test.expected, result)
		})
	}
}

func TestErrResponse_Error(t *testing.T) {
	tests := []struct {
		name     string
		err      *ErrResponse
		expected string
	}{
		{
			name:     "basic error",
			err:      &ErrResponse{Code: -32000, Message: "execution reverted"},
			expected: "rpc error (code -32000): execution reverted",
		},
		{
			name:     "error with data",
			err:      &ErrResponse{Code: -32602, Message: "Invalid params", Data: "extra info"},
			expected: "rpc error (code -32602): Invalid params",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := test.err.Error()
			assert.Equal(t, test.expected, result)
		})
	}
}

func TestErrorCode_MarshalJSONRPC(t *testing.T) {
	tests := []struct {
		name     string
		code     ErrorCode
		expected string
	}{
		{
			name:     "negative error code",
			code:     ErrorCode(-32000),
			expected: "-32000",
		},
		{
			name:     "positive error code",
			code:     ErrorCode(1000),
			expected: "1000",
		},
		{
			name:     "zero error code",
			code:     ErrorCode(0),
			expected: "0",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := test.code.MarshalJSONRPC()
			assert.NoError(t, err)
			assert.Equal(t, test.expected, string(result))
		})
	}
}
