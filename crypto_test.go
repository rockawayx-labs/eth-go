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
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPrivateKey_Generate63BytesPrivateKey(t *testing.T) {
	t.Skip("Used to generate a random private key < 32 bytes")

	for {
		privKey, err := NewRandomPrivateKey()
		require.NoError(t, err)

		if len(privKey.inner.D.Bytes()) < 32 {
			fmt.Printf("D value (Hex %x, Text %s, Bytes %s) has less than 32 bytes\n", privKey.inner.D, privKey.inner.D.Text(10), hex.EncodeToString(privKey.inner.D.Bytes()))
			require.NoError(t, errors.New("bytes < 32"))
		}
	}
}

func TestPrivateKey_String(t *testing.T) {
	tests := []struct {
		in          string
		expectedErr error
	}{
		{in: "52e1cc4b9c8b4fc9b202adf06462bdcc248e170c9abd56b2adb84c8d87bee674", expectedErr: nil},
		{in: "2f6e6a9af650c60e4b2c6a0a1b440cec202182bed3fc15bc5b8eec1132b2d6ad", expectedErr: nil},
		{in: "cb3c1ca36610c116e9aa478102faeaf100cc79c462f0d25a631110e36a2868c8", expectedErr: nil},
		{in: "b1c73d7b1103387b725df649a70d8c2c3cca61a60781f19d0a91ced1ea1be35b", expectedErr: nil},
		{in: "72194c8757fde016fa20cbc1ec09b8e64413c4ad7aba23bdb12664311e03f419", expectedErr: nil},

		// Key < 32 bytes (64 characters) left padded here
		{in: "0033752648ef6373b2904568fc7452957b73bbb4f91657735bb53e633513b805", expectedErr: nil},
	}

	for _, test := range tests {
		t.Run(test.in, func(t *testing.T) {
			actual, err := NewPrivateKey(test.in)
			if test.expectedErr == nil {
				require.NoError(t, err)
				require.NotNil(t, actual, "invalid private for input %q", test.in)
				assert.Equal(t, test.in, actual.String())
			} else {
				assert.Equal(t, test.expectedErr, err)
			}
		})
	}
}

func TestPrivateKey_PublicKey(t *testing.T) {
	tests := []struct {
		in          string
		expected    string
		expectedErr error
	}{
		{
			in:          "52e1cc4b9c8b4fc9b202adf06462bdcc248e170c9abd56b2adb84c8d87bee674",
			expected:    "821b55d8abe79bc98f05eb675fdc50dfe796b7ab",
			expectedErr: nil,
		},
		{
			in:          "2f6e6a9af650c60e4b2c6a0a1b440cec202182bed3fc15bc5b8eec1132b2d6ad",
			expected:    "403a09cd493e41d381ebff4ffb01de4c9f2ff1dc",
			expectedErr: nil,
		},
		{
			in:          "cb3c1ca36610c116e9aa478102faeaf100cc79c462f0d25a631110e36a2868c8",
			expected:    "dfe15828305e9968835e7797f063bd8bbec038c0",
			expectedErr: nil,
		},
		{
			in:          "b1c73d7b1103387b725df649a70d8c2c3cca61a60781f19d0a91ced1ea1be35b",
			expected:    "ccd3dd1db00b7c5b493b879f04379287227d200d",
			expectedErr: nil,
		},
		{
			in:          "72194c8757fde016fa20cbc1ec09b8e64413c4ad7aba23bdb12664311e03f419",
			expected:    "e8c4b557fb717f7a366d4cb7f8fe41ef1e108c16",
			expectedErr: nil,
		},
		{
			in:          "72194c8757fde016fa20cbc1ec09b8e64413c4ad7aba23bdb12664311e03f419",
			expected:    "e8c4b557fb717f7a366d4cb7f8fe41ef1e108c16",
			expectedErr: nil,
		},

		// Key < 32 bytes (64 characters) left padded here
		{
			in:          "0033752648ef6373b2904568fc7452957b73bbb4f91657735bb53e633513b805",
			expected:    "43bf67eae59c0d0eca283839e5cd9b22ca89d530",
			expectedErr: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.in, func(t *testing.T) {
			actual, err := NewPrivateKey(test.in)
			if test.expectedErr == nil {
				require.NoError(t, err)
				require.NotNil(t, actual, "invalid private for input %q", test.in)
				assert.Equal(t, test.expected, actual.PublicKey().Address().String())
			} else {
				assert.Equal(t, test.expectedErr, err)
			}
		})
	}
}

func TestPrivateKey_Bytes(t *testing.T) {
	tests := []struct {
		in          string
		expectedErr error
	}{
		{in: "52e1cc4b9c8b4fc9b202adf06462bdcc248e170c9abd56b2adb84c8d87bee674", expectedErr: nil},
		{in: "2f6e6a9af650c60e4b2c6a0a1b440cec202182bed3fc15bc5b8eec1132b2d6ad", expectedErr: nil},
		{in: "cb3c1ca36610c116e9aa478102faeaf100cc79c462f0d25a631110e36a2868c8", expectedErr: nil},
		{in: "b1c73d7b1103387b725df649a70d8c2c3cca61a60781f19d0a91ced1ea1be35b", expectedErr: nil},
		{in: "72194c8757fde016fa20cbc1ec09b8e64413c4ad7aba23bdb12664311e03f419", expectedErr: nil},

		// Key < 32 bytes (64 characters) left padded here
		{in: "0033752648ef6373b2904568fc7452957b73bbb4f91657735bb53e633513b805", expectedErr: nil},
	}

	for _, test := range tests {
		t.Run(test.in, func(t *testing.T) {
			actual, err := NewPrivateKey(test.in)

			if test.expectedErr == nil {
				require.NoError(t, err)
				require.NotNil(t, actual, "invalid private for input %q", test.in)

				actualBytes := hex.EncodeToString(actual.Bytes())
				assert.Equal(t, test.in, actualBytes)
			} else {
				assert.Equal(t, test.expectedErr, err)
			}
		})
	}
}

func TestPrivateKey_ToECDSA(t *testing.T) {
	tests := []struct {
		in          string
		expectedD   string
		expectedX   string
		expectedY   string
		expectedErr error
	}{
		{
			in:          "52e1cc4b9c8b4fc9b202adf06462bdcc248e170c9abd56b2adb84c8d87bee674",
			expectedX:   "2c8f6d4764c3aca75696e18aeef683932a2bfa0be1603adb54f30dfad8e5cf23",
			expectedY:   "72a9d6eeb0e5caffba1fca22e12878c450e6ef09434888f04c6a97b6f50c75d4",
			expectedErr: nil,
		},
		{
			in:          "2f6e6a9af650c60e4b2c6a0a1b440cec202182bed3fc15bc5b8eec1132b2d6ad",
			expectedX:   "73523c72e55da219c8650a9fe47f83477958a6185b2322b2c948b88e1dcea92",
			expectedY:   "c7e872d816c129b6cf51b05b9c5bc22ddd34fc7d7e35da17fcdbe5598fee56a1",
			expectedErr: nil,
		},
		{
			in:          "cb3c1ca36610c116e9aa478102faeaf100cc79c462f0d25a631110e36a2868c8",
			expectedX:   "b235c1450002e43d49f25de3b479aa7ce98c8b369ebf6ba9dcfee1e8e1ff1a1a",
			expectedY:   "d5f52b3d4347659c35d1d9a44e8bf0b1a70f52d99b85d716e36b74dfe602da26",
			expectedErr: nil,
		},
		{
			in:          "b1c73d7b1103387b725df649a70d8c2c3cca61a60781f19d0a91ced1ea1be35b",
			expectedX:   "ab23cd020e638a977d57907668bf4923108797a8f5860e4c412598bdc2f484b",
			expectedY:   "cda179d557bbcddc83174f47ec45b1e49bd9dbd6bcea26932fc3c25ade0c72b9",
			expectedErr: nil,
		},
		{
			in:          "72194c8757fde016fa20cbc1ec09b8e64413c4ad7aba23bdb12664311e03f419",
			expectedX:   "70203b56d3f40f1823f2de383bd90891d75c00b151e65f4b01d70202f38a9f8",
			expectedY:   "e7d116e7ad511340119b69208c9a78776a8b72ab87359e7a365f0d7f68e313bb",
			expectedErr: nil,
		},

		// Key < 32 bytes (64 characters) left padded here
		{
			in:          "0033752648ef6373b2904568fc7452957b73bbb4f91657735bb53e633513b805",
			expectedD:   "33752648ef6373b2904568fc7452957b73bbb4f91657735bb53e633513b805",
			expectedX:   "685ed2e19b95b5aee50709c579300f93960e3750a106859fd02e75ad2a9d481f",
			expectedY:   "f114535858b7b4856a6e62332599c94a5a751cb31aaea58fad6765c4cf8e786c",
			expectedErr: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.in, func(t *testing.T) {
			actual, err := NewPrivateKey(test.in)
			if test.expectedErr == nil {
				require.NoError(t, err)
				require.NotNil(t, actual, "invalid private for input %q", test.in)

				ecdsaPk := actual.ToECDSA()

				expectedD := test.expectedD
				if expectedD == "" {
					expectedD = test.in
				}

				assert.Equal(t, expectedD, ecdsaPk.D.Text(16), "private coordinate D")
				assert.Equal(t, test.expectedX, ecdsaPk.X.Text(16), "public point X")
				assert.Equal(t, test.expectedY, ecdsaPk.Y.Text(16), "public point Y")
			} else {
				assert.Equal(t, test.expectedErr, err)
			}
		})
	}
}

func TestKeccak256(t *testing.T) {
	tests := []struct {
		in        string
		expectOut string
	}{
		{
			in:        "Pregnant(address,uint256,uint256,uint256)",
			expectOut: "241ea03ca20251805084d27d4440371c34a0b85ff108f6bb5611248f73818b80",
		},
		{
			in:        "Transfer(address,address,uint256)",
			expectOut: "ddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef",
		},
		{
			in:        "Approval(address,address,uint256)",
			expectOut: "8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925",
		},
		{
			in:        "Birth(address,uint256,uint256,uint256,uint256)",
			expectOut: "0a5311bd2a6608f08a180df2ee7c5946819a649b204b554bb8e39825b2c50ad5",
		},
	}

	for _, test := range tests {
		t.Run(test.in, func(t *testing.T) {
			out := Keccak256([]byte(test.in))
			assert.Equal(t, test.expectOut, hex.EncodeToString(out))
		})
	}
}
