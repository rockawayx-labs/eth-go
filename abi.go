package eth

import (
	"fmt"
	"go.uber.org/zap"
	"strings"
)

// ABI is our custom internal definition of a contract's ABI that bridges the information
// between the two ABI like formats for contract's, i.e. `.abi` file and AST file as output
// by `solc` compiler.
type ABI struct {
	LogEventsMap map[string]*LogEventDef
	FunctionsMap map[string]*FunctionDef
}


func (a *ABI) FindLog(topic []byte) *LogEventDef {
	zlog.Info("looking for log event def by topic", zap.Stringer("topic", Hash(topic)))

	return a.LogEventsMap[string(topic)]
}

func (a *ABI) FindFunction(methodHash []byte) *FunctionDef {
	zlog.Info("looking for function by method hash", zap.Stringer("method_hash", Hash(methodHash)))

	return a.FunctionsMap[string(methodHash)]
}

type FunctionDef struct {
	Name             string
	Parameters       []*FunctionParameter
	ReturnParameters []*FunctionParameter
	Payable          bool
	ViewOnly         bool
}

type FunctionParameter struct {
	Name           string
	TypeName       string
	TypeMutability string
	Payable        bool
}

func (f *FunctionDef) methodID() []byte {
	return Keccak256([]byte(f.Signature()))[0:4]
}

func (f *FunctionDef) Signature() string {
	var args []string
	for _, parameter := range f.Parameters {
		args = append(args, parameter.TypeName)
	}

	return fmt.Sprintf("%s(%s)", f.Name, strings.Join(args, ", "))
}

func (f *FunctionDef) String() string {
	var args []string
	for _, parameter := range f.Parameters {
		args = append(args, fmt.Sprintf("%s %s", parameter.TypeName, parameter.Name))
	}

	return fmt.Sprintf("%s(%s)", f.Name, strings.Join(args, ", "))
}

type LogEventDef struct {
	Name       string
	Parameters []*LogParameter
}

type LogParameter struct {
	Name     string
	TypeName string
	Indexed  bool
}

func (p *LogParameter) GetName(index int) string {
	name := p.Name
	if name == "" {
		name = fmt.Sprintf("unamed%d", index+1)
	}

	return name
}

func (l *LogEventDef) logID() []byte {
	return Keccak256([]byte(l.Signature()))
}

func (l *LogEventDef) Signature() string {
	var args []string
	for _, parameter := range l.Parameters {
		args = append(args, parameter.TypeName)
	}

	return fmt.Sprintf("%s(%s)", l.Name, strings.Join(args, ", "))
}

func (l *LogEventDef) String() string {
	var args []string
	for i, parameter := range l.Parameters {
		indexed := ""
		if parameter.Indexed {
			indexed = " indexed"
		}

		args = append(args, fmt.Sprintf("%s%s %s", parameter.TypeName, indexed, parameter.GetName(i)))
	}

	return fmt.Sprintf("%s(%s)", l.Name, strings.Join(args, ", "))
}
