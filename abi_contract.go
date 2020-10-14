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
