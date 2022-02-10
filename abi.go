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

import "go.uber.org/zap"

// ABI is our custom internal definition of a contract's ABI that bridges the information
// between the two ABI like formats for contract's, i.e. `.abi` file and AST file as output
// by `solc` compiler.
type ABI struct {
	LogEventsMap    map[string]*LogEventDef
	FunctionsMap    map[string]*MethodDef
	ConstructorsMap map[string]*MethodDef

	LogEventsByNameMap    map[string]*LogEventDef
	FunctionsByNameMap    map[string]*MethodDef
	ConstructorsByNameMap map[string]*MethodDef
}

func (a *ABI) FindLog(topic []byte) *LogEventDef {
	zlog.Info("looking for log event by topic", zap.Stringer("topic", Hash(topic)))

	return a.LogEventsMap[string(topic)]
}

func (a *ABI) FindFunction(functionHash []byte) *MethodDef {
	zlog.Info("looking for function by hash", zap.Stringer("method_hash", Hash(functionHash)))

	return a.FunctionsMap[string(functionHash)]
}

func (a *ABI) FindFunctionByName(name string) *MethodDef {
	zlog.Info("looking for function by name", zap.Stringer("method_name", Hash(name)))

	return a.FunctionsByNameMap[name]
}
