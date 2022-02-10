package rpc

import (
	"reflect"
	"testing"

	"github.com/streamingfast/eth-go"
)

func TestNewTopicFilter(t *testing.T) {
	type args struct {
		exprs    []interface{}
		appender func(f *TopicFilter)
	}
	tests := []struct {
		name    string
		args    args
		wantOut *TopicFilter
	}{
		{
			"topic0",
			args{exprs: []interface{}{"0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"}},
			&TopicFilter{
				topics: []TopicFilterExpr{{exact: topic("0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef")}},
			},
		},
		{
			"topic0 with append null",
			args{
				exprs: []interface{}{"0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"},
				appender: func(f *TopicFilter) {
					f.Append(nil)
					f.Append(eth.MustNewAddress("0xFffDB7377345371817F2b4dD490319755F5899eC"))
				},
			},
			&TopicFilter{
				topics: []TopicFilterExpr{
					{exact: topic("0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef")},
					{exact: nil},
					{exact: topic("0x000000000000000000000000FffDB7377345371817F2b4dD490319755F5899eC")},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := NewTopicFilter(tt.args.exprs...)
			if tt.args.appender != nil {
				tt.args.appender(f)
			}

			if gotOut := f; !reflect.DeepEqual(gotOut, tt.wantOut) {
				t.Errorf("NewTopicFilter() = %v, want %v", gotOut, tt.wantOut)
			}
		})
	}
}

func topic(in string) (out *eth.Topic) {
	var bytes [32]byte
	copy(bytes[:], eth.MustNewHash(in))
	return (*eth.Topic)(&bytes)
}
