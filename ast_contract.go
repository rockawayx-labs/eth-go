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
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"go.uber.org/zap"
)

func ParseAST(astFilepath string) *ABI {
	file, _ := ioutil.ReadFile(astFilepath)
	ast := map[string]interface{}{}
	_ = json.Unmarshal([]byte(file), &ast)

	abi := &ABI{
		MethodsMap: make(map[string]*MethodDef),
	}

	for _, node := range ast["nodes"].([]interface{}) {
		if n, ok := node.(map[string]interface{}); ok {
			abi = convertJsonToASTNode(abi, n)
		}
	}

	return abi
}

func convertJsonToASTNode(abi *ABI, node map[string]interface{}) *ABI {
	nodeType := node["nodeType"].(string)
	switch nodeType {
	case "ContractDefinition":
		abi = createContractDefinition(abi, node)
	case "FunctionDefinition":
		f, err := createFunctionDefinition(node)
		if err != nil {
			//zlog.Warn("error creating function", zap.Error(err))
		} else {
			abi.MethodsMap[string(f.MethodID())] = f

		}
	default:
		// zlog.Info("unhandled node type", zap.String("node_type", nodeType))
	}
	return abi
}

func createFunctionDefinition(node map[string]interface{}) (*MethodDef, error) {
	if _, ok := node["kind"].(string); !ok {
		zlog.Fatal("expected 'kind' to be a string!")
	}
	nodeKind := node["kind"].(string)
	switch nodeKind {
	//case "constructor":
	case "function":
		f, err := createFunctionDefinitionFunc(node)
		if err != nil {
			return nil, fmt.Errorf("error decoding function definition: %w", err)
		}
		return f, nil
	//case "fallback":
	//case "receive":
	default:
		return nil, fmt.Errorf("Expected 'kind' to be one of [constructor, function, fallback, receive]: %q", nodeKind)
	}

}

func createContractDefinition(abi *ABI, node map[string]interface{}) *ABI {
	if _, ok := node["name"].(string); !ok {
		zlog.Fatal("expected 'name' to be a string")
	}

	for _, node := range node["nodes"].([]interface{}) {
		if n, ok := node.(map[string]interface{}); ok {
			abi = convertJsonToASTNode(abi, n)
		}
	}
	return abi
}

func createFunctionDefinitionFunc(node map[string]interface{}) (*MethodDef, error) {
	if _, ok := node["name"].(string); !ok {
		return nil, fmt.Errorf("expected 'name' to be a string")
	}

	f := &MethodDef{
		Name:             node["name"].(string),
		Parameters:       []*MethodParameter{},
		ReturnParameters: []*MethodParameter{},
	}

	if parameterList, ok := node["parameters"].(map[string]interface{}); ok {
		if traceEnabled {
			zlog.Debug("parsing function parameters", zap.String("name", f.Name))
		}
		f.Parameters = getFunctionParameters(parameterList["parameters"].([]interface{}))
		if traceEnabled {
			zlog.Debug("function found", zap.String("name", f.Name), zap.Reflect("parameters", f.Parameters))
		}
	} else {
		zlog.Warn("no parameter list for function", zap.String("function_name", f.Name))
	}

	if parameterList, ok := node["returnParameters"].(map[string]interface{}); ok {
		if traceEnabled {
			zlog.Debug("parsing function parameters", zap.String("name", f.Name))
		}
		f.ReturnParameters = getFunctionParameters(parameterList["parameters"].([]interface{}))
		if traceEnabled {
			zlog.Debug("function found", zap.String("name", f.Name), zap.Reflect("parameters", f.Parameters))
		}

	} else {
		zlog.Warn("no return parameter list for function", zap.String("function_name", f.Name))
	}

	return f, nil
}

func getFunctionParameters(parameters []interface{}) (out []*MethodParameter) {
	for _, parameter := range parameters {
		if param, ok := parameter.(map[string]interface{}); ok {
			p := &MethodParameter{Name: param["name"].(string)}
			if varType, ok := param["typeName"].(map[string]interface{}); ok {
				p.TypeName = varType["name"].(string)
				if str, ok := varType["stateMutability"].(string); ok {
					p.TypeMutability = str
				}
				out = append(out, p)
			} else {
				zlog.Warn("expected a 'parameter.typeName' to be a map[string]interface{}")
			}
			continue
		} else {
			zlog.Warn("expected a 'parameter' to be a map[string]interface{}")
		}
	}
	return
}

func B(input string) []byte {
	data, err := hex.DecodeString(SanitizeHex(input))
	if err != nil {
		panic(err)
	}

	return data
}
