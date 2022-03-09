package signer

import (
	"reflect"
	"testing"

	"github.com/streamingfast/eth-go"
)

func TestRecoverPersonalSigner(t *testing.T) {
	type args struct {
		signature eth.Hex
		hash      eth.Hash
	}

	tests := []struct {
		name    string
		args    args
		want    eth.Address
		wantErr bool
	}{
		{
			"standard",
			args{
				signature: eth.MustNewHex("cfc0c9160c1fbe884f02298b76e194611904a4f18b6814dfd22f95b659b2f90b2564e455543316f125c3c51ec72b22917401d4872ddabb9a7a2d0bea1a827db31b"),
				hash:      eth.MustNewHash("0xcf36ac4f97dc10d91fc2cbb20d718e94a8cbfe0f82eaedc6a4aa38946fb797cd"),
			},
			eth.MustNewAddress("0xfffdb7377345371817f2b4dd490319755f5899ec"),
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := RecoverPersonalSigner(tt.args.signature, tt.args.hash)
			if (err != nil) != tt.wantErr {
				t.Errorf("RecoverPersonalSigner() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RecoverPersonalSigner() = %v, want %v", got, tt.want)
			}
		})
	}
}
