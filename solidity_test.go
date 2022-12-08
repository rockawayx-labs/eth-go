package eth

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseType(t *testing.T) {
	hasError := func(in string) require.ErrorAssertionFunc {
		return func(tt require.TestingT, err error, i ...interface{}) {
			require.EqualError(tt, err, in, i...)
		}
	}

	type args struct {
		raw string
	}
	tests := []struct {
		name      string
		args      args
		want      SolidityType
		assertion require.ErrorAssertionFunc
	}{
		{"boolean", args{"bool"}, BooleanType{}, require.NoError},
		{"address", args{"address"}, AddressType{}, require.NoError},
		{"bytes", args{"bytes"}, BytesType{}, require.NoError},
		{"string", args{"string"}, StringType{}, require.NoError},
		{"tuple", args{"tuple"}, StructType{}, require.NoError},

		{"int", args{"int"}, SignedIntegerType{BitsSize: 256, ByteSize: 32}, require.NoError},
		{"int0", args{"int0"}, nil, hasError("invalid integer type \"int0\": bits size 0 is lower than 8")},
		{"int7", args{"int7"}, nil, hasError("invalid integer type \"int7\": bits size 7 is lower than 8")},
		{"int8", args{"int8"}, SignedIntegerType{BitsSize: 8, ByteSize: 1}, require.NoError},
		{"int9", args{"int9"}, nil, hasError("invalid integer type \"int9\": bits size 9 must be divisble by 8")},
		{"int256", args{"int256"}, SignedIntegerType{BitsSize: 256, ByteSize: 32}, require.NoError},
		{"int257", args{"int257"}, nil, hasError("invalid integer type \"int257\": bits size 257 is bigger than 256")},
		{"int262", args{"int262"}, nil, hasError("invalid integer type \"int262\": bits size 262 is bigger than 256")},

		{"uint", args{"uint"}, UnsignedIntegerType{BitsSize: 256, ByteSize: 32}, require.NoError},
		{"uint0", args{"uint0"}, nil, hasError("invalid integer type \"uint0\": bits size 0 is lower than 8")},
		{"uint7", args{"uint7"}, nil, hasError("invalid integer type \"uint7\": bits size 7 is lower than 8")},
		{"uint8", args{"uint8"}, UnsignedIntegerType{BitsSize: 8, ByteSize: 1}, require.NoError},
		{"uint9", args{"uint9"}, nil, hasError("invalid integer type \"uint9\": bits size 9 must be divisble by 8")},
		{"uint256", args{"uint256"}, UnsignedIntegerType{BitsSize: 256, ByteSize: 32}, require.NoError},
		{"uint257", args{"uint257"}, nil, hasError("invalid integer type \"uint257\": bits size 257 is bigger than 256")},
		{"uint262", args{"uint262"}, nil, hasError("invalid integer type \"uint262\": bits size 262 is bigger than 256")},

		{"fixed", args{"fixed"}, SignedFixedPointType{BitsSize: 128, ByteSize: 16, Decimals: 18}, require.NoError},
		{"fixed0x18", args{"fixed0x18"}, nil, hasError("invalid fixed point type \"fixed0x18\": bits size 0 is lower than 8")},
		{"fixed7x18", args{"fixed7x18"}, nil, hasError("invalid fixed point type \"fixed7x18\": bits size 7 is lower than 8")},
		{"fixed8x18", args{"fixed8x18"}, SignedFixedPointType{BitsSize: 8, ByteSize: 1, Decimals: 18}, require.NoError},
		{"fixed9x18", args{"fixed9x18"}, nil, hasError("invalid fixed point type \"fixed9x18\": bits size 9 must be divisble by 8")},
		{"fixed256x18", args{"fixed256x18"}, SignedFixedPointType{BitsSize: 256, ByteSize: 32, Decimals: 18}, require.NoError},
		{"fixed257x18", args{"fixed257x18"}, nil, hasError("invalid fixed point type \"fixed257x18\": bits size 257 is bigger than 256")},
		{"fixed262x18", args{"fixed262x18"}, nil, hasError("invalid fixed point type \"fixed262x18\": bits size 262 is bigger than 256")},

		{"ufixed", args{"ufixed"}, UnsignedFixedPointType{BitsSize: 128, ByteSize: 16, Decimals: 18}, require.NoError},
		{"ufixed0x18", args{"ufixed0x18"}, nil, hasError("invalid fixed point type \"ufixed0x18\": bits size 0 is lower than 8")},
		{"ufixed7x18", args{"ufixed7x18"}, nil, hasError("invalid fixed point type \"ufixed7x18\": bits size 7 is lower than 8")},
		{"ufixed8x18", args{"ufixed8x18"}, UnsignedFixedPointType{BitsSize: 8, ByteSize: 1, Decimals: 18}, require.NoError},
		{"ufixed9x18", args{"ufixed9x18"}, nil, hasError("invalid fixed point type \"ufixed9x18\": bits size 9 must be divisble by 8")},
		{"ufixed256x18", args{"ufixed256x18"}, UnsignedFixedPointType{BitsSize: 256, ByteSize: 32, Decimals: 18}, require.NoError},
		{"ufixed257x18", args{"ufixed257x18"}, nil, hasError("invalid fixed point type \"ufixed257x18\": bits size 257 is bigger than 256")},
		{"ufixed262x18", args{"ufixed262x18"}, nil, hasError("invalid fixed point type \"ufixed262x18\": bits size 262 is bigger than 256")},

		{"bytes0", args{"bytes0"}, nil, hasError("invalid fixed bytes type \"bytes0\": bits size 0 is lower than 1")},
		{"bytes1", args{"bytes1"}, FixedSizeBytesType{ByteSize: 1}, require.NoError},
		{"bytes32", args{"bytes32"}, FixedSizeBytesType{ByteSize: 32}, require.NoError},
		{"bytes33", args{"bytes33"}, nil, hasError("invalid fixed bytes type \"bytes33\": bits size 33 is bigger than 32")},

		{"uint8[]", args{"uint8[]"}, ArrayType{ElementType: UnsignedIntegerType{BitsSize: 8, ByteSize: 1}}, require.NoError},
		{"bob[]", args{"bob[]"}, nil, hasError("invalid array type \"bob[]\": element type invalid: type \"bob\" unknown")},
		{"uint8[3]", args{"uint8[3]"}, FixedSizeArrayType{ElementType: UnsignedIntegerType{BitsSize: 8, ByteSize: 1}, Length: 3}, require.NoError},
		{"bob[3]", args{"bob[3]"}, nil, hasError("invalid fixed array type \"bob[3]\": element type invalid: type \"bob\" unknown")},

		{"tuple[]", args{"tuple[]"}, ArrayType{ElementType: StructType{}}, require.NoError},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseType(tt.args.raw)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
