package eth

import "go.uber.org/zap"

// ABI is our custom internal definition of a contract's ABI that bridges the information
// between the two ABI like formats for contract's, i.e. `.abi` file and AST file as output
// by `solc` compiler.
type ABI struct {
	LogEventsMap map[string]*LogEventDef
	MethodsMap   map[string]*MethodDef
}

func (a *ABI) FindLog(topic []byte) *LogEventDef {
	zlog.Info("looking for log event def by topic", zap.Stringer("topic", Hash(topic)))

	return a.LogEventsMap[string(topic)]
}

func (a *ABI) FindMethod(methodHash []byte) *MethodDef {
	zlog.Info("looking for function by method hash", zap.Stringer("method_hash", Hash(methodHash)))

	return a.MethodsMap[string(methodHash)]
}
