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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogEventDef_Signature(t *testing.T) {
	type fields struct {
		Name       string
		Parameters []*LogParameter
	}
	tests := []struct {
		want   string
		fields fields
	}{
		{
			"EventAddressIdxString(address,string)",
			fields{Name: "EventAddressIdxString", Parameters: []*LogParameter{
				{Name: "any", TypeName: "address", Indexed: true},
				{Name: "any", TypeName: "string", Indexed: false},
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			l := &LogEventDef{
				Name:       tt.fields.Name,
				Parameters: tt.fields.Parameters,
			}
			assert.Equal(t, tt.want, l.Signature())
		})
	}
}
