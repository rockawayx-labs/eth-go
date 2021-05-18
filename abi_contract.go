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

package eth

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

func ParseABI(abiFilePath string) (*ABI, error) {
	file, _ := ioutil.ReadFile(abiFilePath)

	return parseABIFromBytes(file)
}

func parseABIFromBytes(input []byte) (*ABI, error) {
	// Could we "stream" it using `json.NewDecoder` and save on memory for this?
	var declarations []*declaration
	if err := json.Unmarshal(input, &declarations); err != nil {
		return nil, fmt.Errorf("read abi: %w", err)
	}

	abi := &ABI{
		LogEventsMap: map[string]*LogEventDef{},
		MethodsMap:   map[string]*MethodDef{},
	}

	for _, decl := range declarations {
		if decl.Type == "event" {
			logEventDef := decl.toLogEventDef()
			abi.LogEventsMap[string(logEventDef.logID())] = logEventDef
		}

		if decl.Type == "function" {

		}
	}

	return abi, nil
}

type declaration struct {
	Name            string      `json:"name,omitempty"`
	Type            string      `json:"type"`
	Inputs          []*typeInfo `json:"inputs"`
	Outputs         []*typeInfo `json:"outputs,omitempty"`
	Payable         bool        `json:"payable,omitempty"`
	StateMutability string      `json:"stateMutability,omitempty"`
	Anonymous       bool        `json:"anonymous,omitempty"`
	Constant        bool        `json:"constant,omitempty"`
}

func (d *declaration) toFunctionDef() *MethodDef {
	out := &MethodDef{}
	out.Name = d.Name
	out.Payable = d.Payable
	out.ViewOnly = d.StateMutability == "view"

	out.Parameters = make([]*MethodParameter, len(d.Inputs))
	for i, input := range d.Inputs {
		out.Parameters[i] = &MethodParameter{
			Name:     input.Name,
			TypeName: input.Type,
		}
	}

	out.ReturnParameters = make([]*MethodParameter, len(d.Outputs))
	for i, input := range d.Outputs {
		out.Parameters[i] = &MethodParameter{
			Name:     input.Name,
			TypeName: input.Type,
		}
	}

	return out
}

func (d *declaration) toLogEventDef() *LogEventDef {
	out := &LogEventDef{}
	out.Name = d.Name

	out.Parameters = make([]*LogParameter, len(d.Inputs))
	for i, input := range d.Inputs {
		out.Parameters[i] = &LogParameter{
			Name:     input.Name,
			TypeName: input.Type,
			Indexed:  input.Indexed,
		}
	}

	return out
}

type typeInfo struct {
	InternalType string `json:"internalType"`
	Name         string `json:"name"`
	Type         string `json:"type"`
	Indexed      bool   `json:"indexed"`
}
