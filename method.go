package eth

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"regexp"
	"strconv"
	"strings"
)

type MethodParameter struct {
	Name           string
	TypeName       string
	TypeMutability string
	Payable        bool
}

type MethodDef struct {
	Name             string
	Parameters       []*MethodParameter
	ReturnParameters []*MethodParameter
	Payable          bool
	ViewOnly         bool
}

func NewMethodDef(signature string) (*MethodDef, error) {
	methodName := extractMethodNameFromSignature(signature)
	if methodName == "" {
		return nil, fmt.Errorf("invalid signature: %s", signature)
	}
	params, err := extractInputsFromSignature(signature)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve inputs %q: %w", signature, err)
	}

	return &MethodDef{
		Name:       methodName,
		Parameters: params,
	}, nil
}
func (f *MethodDef) NewCall() *MethodCall {
	return &MethodCall{
		methodDef: f,
	}
}

func (f *MethodDef) methodID() []byte {
	return Keccak256([]byte(f.Signature()))[0:4]
}

func (f *MethodDef) Signature() string {
	var args []string
	for _, parameter := range f.Parameters {
		args = append(args, parameter.TypeName)
	}

	return fmt.Sprintf("%s(%s)", f.Name, strings.Join(args, ","))
}

func (f *MethodDef) String() string {
	var args []string
	for _, parameter := range f.Parameters {
		args = append(args, fmt.Sprintf("%s %s", parameter.TypeName, parameter.Name))
	}

	return fmt.Sprintf("%s(%s)", f.Name, strings.Join(args, ", "))
}

type MethodCall struct {
	methodDef *MethodDef
	data      []interface{}
}

func (f *MethodCall) AppendArgFromString(v string) error {
	i := len(f.data)
	if i >= len(f.methodDef.Parameters) {
		return fmt.Errorf("args exceeds method definition parameter count %d", len(f.methodDef.Parameters))
	}
	param := f.methodDef.Parameters[i]
	var out interface{}
	switch param.TypeName {
	case "bytes":
		data, err := hex.DecodeString(SanitizeHex(v))
		if err != nil {
			return err
		}
		out = data
	case "address":
		var addr Address
		err := json.Unmarshal([]byte(fmt.Sprintf("%q", v)), &addr)
		if err != nil {
			return err
		}
		out = addr
	case "uint64":
		v, err := strconv.ParseUint(v, 10, 64)
		if err != nil {
			return err
		}
		out = v
	case "uint112", "uint256":
		var ok bool
		out, ok = new(big.Int).SetString(v, 10)
		if !ok {
			fmt.Errorf("unable to convert %q to %s ", v, param.TypeName)
		}
	case "bool":
		out = v == "true"
	}
	f.data = append(f.data, out)
	return nil
}

func (f *MethodCall) AppendArg(v interface{}) {

}

var methodRE = regexp.MustCompile(`(.*)\(`)
var methodInputsRE = regexp.MustCompile(`\((.*?)\)`)

//
//type MethodParameters struct {
//	Type string `json:"type"`
//	Name string `json:"name"`
//}
//
//type MethodSignature struct {
//	Raw        string `json:"raw"`
//	Name       string `json:"name"`
//	Parameters []*MethodParameters
//}
//
//func NewMethodSignature(signature string) (*MethodSignature, error) {
//	methodName := extractMethodNameFromSignature(signature)
//	if methodName == "" {
//		return nil, fmt.Errorf("invalid signature %s", signature)
//	}
//
//	params, err := extractInputsFromSignature(signature)
//	if err != nil {
//		return nil, fmt.Errorf("unable to retrieve inputs %q: %w", signature, err)
//	}
//
//	return &MethodSignature{
//		Raw:        signature,
//		Name:       methodName,
//		Parameters: params,
//	}, nil
//}
//
//func (m *MethodSignature) NewCall() *MethodCall {
//
//	return &MethodCall{
//		Signature: MethodSignature{},
//		Arguments: nil,
//	}
//}
//
//type MethodCall struct {
//	Signature MethodSignature `json:"signature"`
//	Arguments []interface{}   `json:"args"`
//}
//
////
////type Input struct {
////	Type string `json:"type"`
////	// TODO: Method is a struct to model the input data or a Ethereum CALL functions, should it also contain the potentially value on the decoding side?
////	Value interface{} `json:"value"`
////}
//
////func NewMethodFromJSON(cnt []byte) (*Method, error) {
////	var m *Method
////	err := json.Unmarshal(cnt, &m)
////	if err != nil {
////		return nil, err
////	}
////
////	for _, input := range m.Inputs {
////		switch input.Type {
////		case "address":
////			addr, err := NewAddress(SanitizeHex(input.Value.(string)))
////			if err != nil {
////				return nil, fmt.Errorf("unable to unmarshal address: %w", err)
////			}
////			input.Value = addr
////		case "uint64", "uint112", "uint256":
////			input.Value, _ = new(big.Int).SetString(SanitizeHex(input.Value.(string)), 10)
////		case "bool":
////			input.Value = input.Value.(bool)
////		case "bytes":
////			d, err := hex.DecodeString(SanitizeHex(input.Value.(string)))
////			if err != nil {
////				return nil, fmt.Errorf("unable to unmarshal address: %w", err)
////			}
////			input.Value = d
////		}
////	}
////
////	return m, err
////}
//

func extractMethodNameFromSignature(signature string) string {
	methodName := methodRE.FindString(signature)
	methodName = strings.TrimRight(methodName, "(")
	return strings.Replace(methodName, " ", "", -1) // this should not do anything
}

func extractInputsFromSignature(signature string) (out []*MethodParameter, err error) {
	types, err := extractTypesFromSignature(signature)
	if err != nil {
		return nil, err
	}
	for _, t := range types {
		out = append(out, &MethodParameter{
			TypeName: t,
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
