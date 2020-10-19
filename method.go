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
		MethodDef: f,
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
	MethodDef *MethodDef
	Data      []interface{}

	err []error
}

func (f *MethodCall) AppendArgFromString(v string) {
	i := len(f.Data)
	if i >= len(f.MethodDef.Parameters) {
		f.err = append(f.err, fmt.Errorf("args exceeds method definition parameter count %d", len(f.MethodDef.Parameters)))
		return
	}
	param := f.MethodDef.Parameters[i]
	var out interface{}
	switch param.TypeName {
	case "bytes":
		data, err := hex.DecodeString(SanitizeHex(v))
		if err != nil {
			f.err = append(f.err, fmt.Errorf("unable to convert %q to bytes: %w", v, err))
			return
		}
		out = data
	case "address":
		var addr Address
		err := json.Unmarshal([]byte(fmt.Sprintf("%q", v)), &addr)
		if err != nil {
			f.err = append(f.err, fmt.Errorf("unable to convert %q to address: %w", v, err))
			return
		}
		out = addr
	case "uint64":
		v, err := strconv.ParseUint(v, 10, 64)
		if err != nil {
			f.err = append(f.err, fmt.Errorf("unable to convert %q to uint64: %w", v, err))
			return
		}
		out = v
	case "uint112", "uint256":
		var ok bool
		out, ok = new(big.Int).SetString(v, 10)
		if !ok {
			f.err = append(f.err, fmt.Errorf("unable to convert %q to %s ", v, param.TypeName))
			return
		}
	case "bool":
		out = v == "true"
	}
	f.Data = append(f.Data, out)
}

func (f *MethodCall) AppendArg(v interface{}) {
	f.Data = append(f.Data, v)
}

func (f *MethodCall) Encode() ([]byte, error) {
	if len(f.err) > 0 {
		return nil, fmt.Errorf("%s", f.err)

	}
	enc := NewEncoder()
	err := enc.WriteMethodCall(f)
	if err != nil {
		return nil, err
	}
	return enc.Buffer(), nil
}

var methodRE = regexp.MustCompile(`(.*)\(`)
var methodInputsRE = regexp.MustCompile(`\((.*?)\)`)

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
		return []string{}, nil
	}
	return strings.Split(s, ","), nil
}
