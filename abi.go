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
