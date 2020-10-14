package eth

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"github.com/dfuse-io/eth-go/constants"
	"go.uber.org/zap"
)

type Decoder struct {
	buffer []byte
	offset uint64
	total  uint64
}

func NewDecoderFromString(input string) (*Decoder, error) {
	data, err := hex.DecodeString(SanitizeHex(input))
	if err != nil {
		return nil, fmt.Errorf("unable to decode hex input %q: %w", input, err)
	}

	return NewDecoder(data), nil
}

func NewDecoder(input []byte) *Decoder {
	return &Decoder{
		buffer: input,
		offset: 0,
		total:  uint64(len(input)),
	}
}

func (d *Decoder) String() string {
	return fmt.Sprintf("offset %d, total: %d", d.offset, d.total)
}

func (d *Decoder) SetBytes(input []byte) *Decoder {
	d.buffer = input
	d.offset = 0
	d.total = uint64(len(input))

	return d
}

func (d *Decoder) ReadWithMethodCall() (*MethodCall, error) {
	methodSignature, err := d.ReadMethod()
	if err != nil {
		return nil, err
	}

	methodDef, err := NewMethodDef(methodSignature)
	if err != nil {
		return nil, err
	}

	methodCall := methodDef.NewCall()

	for _, param := range methodCall.methodDef.Parameters {
		var currentOffset uint64

		isAnArray, _ := isArray(param.TypeName)
		if isAnArray {
			currentOffset = d.offset
			jumpToOffset, err := d.read("uint256")
			if err != nil {
				return nil, fmt.Errorf("unable to array lenght %w", err)
			}
			// 4 bytes here is to take into account the method name
			d.offset = (jumpToOffset.(*big.Int).Uint64() + 4)
		}

		out, err := d.Read(param.TypeName)
		if err != nil {
			return nil, fmt.Errorf("unable to decode method input: %w", err)
		}

		if isAnArray {
			d.offset = (currentOffset + 32)
		}

		methodCall.AppendArg(out)
	}
	return methodCall, nil
}

func (d *Decoder) Read(typeName string) (interface{}, error) {
	var isAnArray bool
	isAnArray, resolvedTypeName := isArray(typeName)
	if !isAnArray {
		return d.read(resolvedTypeName)
	}

	length, err := d.read("uint256")
	if err != nil {
		return nil, fmt.Errorf("cannot write slice %s size: %w", typeName, err)
	}

	size := length.(*big.Int).Uint64()
	// little optimization for byte array
	if resolvedTypeName == "bytes" {
		return d.readBytes(size)
	}

	arr, err := newArray(resolvedTypeName, size)
	if err != nil {
		return nil, fmt.Errorf("cannot setup new array: %w", err)
	}

	for i := uint64(0); i < size; i++ {
		out, err := d.read(resolvedTypeName)
		if err != nil {
			return nil, fmt.Errorf("cannot write item from slice %s.%d: %w", typeName, i, err)
		}
		arr.At(i, out)
	}
	return arr, nil
}

func (d *Decoder) read(typeName string) (out interface{}, err error) {
	switch typeName {
	case "address":
		return d.ReadAddress()
	case "uint112":
		return d.ReadUint112()
	case "uint256":
		return d.ReadUint256()
	case "string":
		return d.ReadString()
	case "method":
		return d.ReadMethod()
	}

	return nil, fmt.Errorf("type %q is not handled right now", typeName)
}

func (d *Decoder) ReadMethod() (out string, err error) {
	data, err := d.readBytes(4)
	if err != nil {
		return out, err
	}
	idx := hex.EncodeToString(data)
	out, ok := constants.Signatures[idx]
	if !ok {
		return "", fmt.Errorf("method signature not found for %s", idx)
	}
	return out, nil
}

func (d *Decoder) ReadString() (out string, err error) {
	prevLocation := d.offset
	offset, err := d.ReadUint256()
	if err != nil {
		return out, err
	}
	d.offset = prevLocation + offset.Uint64()
	size, err := d.ReadUint256()
	if err != nil {
		return out, err
	}

	remaining := 32 - (size.Uint64() % 32)
	data, err := d.readBytes(size.Uint64())
	if err != nil {
		return out, err
	}
	out = string(data)
	d.offset += remaining

	return
}

func (d *Decoder) ReadAddress() (out Address, err error) {
	data, err := d.readBytes(32)
	if err != nil {
		return out, err
	}

	address := Address(data[12:])
	if traceEnabled {
		zlog.Debug("read address", zap.Stringer("value", address))
	}

	return address, nil
}

func (d *Decoder) ReadUint112() (out *big.Int, err error) {
	return d.readBigInt(112)
}

func (d *Decoder) ReadUint256() (out *big.Int, err error) {
	return d.readBigInt(256)
}

func (d *Decoder) readBigInt(bits int) (out *big.Int, err error) {
	data, err := d.readBytes(32)
	if err != nil {
		return out, err
	}

	out = new(big.Int).SetBytes(data[:])
	if traceEnabled {
		zlog.Debug("read uint"+strconv.FormatInt(int64(bits), 10), zap.String("value", out.Text(16)))
	}

	return
}
func (d *Decoder) readBytes(byteCount uint64) ([]byte, error) {
	if traceEnabled {
		zlog.Debug("trying to read bytes", zap.Uint64("byte_count", byteCount), zap.Uint64("remaining", d.total-d.offset))
	}

	if d.total-d.offset < byteCount {
		return nil, fmt.Errorf("not enough bytes to read %d bytes, only %d remaining", byteCount, d.total-d.offset)
	}

	out := d.buffer[d.offset : d.offset+byteCount]
	if traceEnabled {
		zlog.Debug("read bytes", zap.Uint64("byte_count", byteCount), zap.String("data", hex.EncodeToString(out)))
	}

	d.offset += byteCount

	return out, nil
}

type decodedArray interface {
	At(index uint64, value interface{})
}

func newArray(typeName string, count uint64) (decodedArray, error) {
	switch typeName {
	case "address":
		return addressArray(make([]Address, count)), nil
	case "uint112":
		return bigIntArray(make([]*big.Int, count)), nil
	case "uint256":
		return bigIntArray(make([]*big.Int, count)), nil
	case "string":
		return stringArray(make([]string, count)), nil
	}
	return nil, fmt.Errorf("type %q is not handled right now", typeName)
}

type stringArray []string

func (a stringArray) At(index uint64, value interface{}) {
	([]string)(a)[index] = value.(string)
}

type addressArray []Address

func (a addressArray) At(index uint64, value interface{}) {
	([]Address)(a)[index] = value.(Address)
}

type bigIntArray []*big.Int

func (a bigIntArray) At(index uint64, value interface{}) {
	([]*big.Int)(a)[index] = value.(*big.Int)
}

type biteArray []byte

func (a biteArray) At(index uint64, value interface{}) {
	([]byte)(a)[index] = value.(byte)
}

func MustHexDecode(input string) []byte {
	input = strings.TrimPrefix(input, "0x")
	value, err := hex.DecodeString(input)
	if err != nil {
		panic(fmt.Errorf("should have been possible to transform decode %q as hex: %s", input, err))
	}

	return value
}
