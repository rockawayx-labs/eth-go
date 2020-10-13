package eth

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"regexp"
	"strings"
)

var methodRE = regexp.MustCompile(`(.*)\(`)
var methodInputsRE = regexp.MustCompile(`\((.*?)\)`)

type Method struct {
	Signature string   `json:"signature"`
	Inputs    []*Input `json:"inputs"`
}

type Input struct {
	Type string `json:"type"`
	// TODO: Method is a struct to model the input data or a Ethereum CALL functions, should it also contain the potentially value on the decoding side?
	Value interface{} `json:"value"`
}

func NewMethodFromJSON(cnt []byte) (*Method, error) {
	var m *Method
	err := json.Unmarshal(cnt, &m)
	if err != nil {
		return nil, err
	}

	for _, input := range m.Inputs {
		switch input.Type {
		case "address":
			addr, err := NewAddress(SanitizeHex(input.Value.(string)))
			if err != nil {
				return nil, fmt.Errorf("unable to unmarshal address: %w", err)
			}
			input.Value = addr
		case "uint64", "uint112", "uint256":
			input.Value, _ = new(big.Int).SetString(SanitizeHex(input.Value.(string)), 10)
		case "bool":
			input.Value = input.Value.(bool)
		case "bytes":
			d, err := hex.DecodeString(SanitizeHex(input.Value.(string)))
			if err != nil {
				return nil, fmt.Errorf("unable to unmarshal address: %w", err)
			}
			input.Value = d
		}
	}

	return m, err
}

func NewMethodFromSignature(signature string) (*Method, error) {
	methodName := extractMethodNameFromSignature(signature)
	if methodName == "" {
		return nil, fmt.Errorf("invalid signature %s", signature)
	}

	inputs, err := extractInputsFromSignature(signature)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve inputs %q: %w", signature, err)
	}

	return &Method{
		Signature: signature,
		Inputs:    inputs,
	}, nil
}

func extractMethodNameFromSignature(signature string) string {
	methodName := methodRE.FindString(signature)
	methodName = strings.TrimRight(methodName, "(")
	return strings.Replace(methodName, " ", "", -1) // this should not do anything
}

func extractInputsFromSignature(signature string) (out []*Input, err error) {
	types, err := extractTypesFromSignature(signature)
	if err != nil {
		return nil, err
	}
	for _, t := range types {
		out = append(out, &Input{
			Type: t,
		})
	}
	return out, nil
}

func extractTypesFromSignature(method string) ([]string, error) {
	s := methodInputsRE.FindString(method)
	s = strings.TrimLeft(s, "(")
	s = strings.TrimRight(s, ")")
	s = strings.Replace(s, " ", "", -1)
	if s == "" {
		return nil, fmt.Errorf("invalid method %s", method)
	}

	return strings.Split(s, ","), nil
}
