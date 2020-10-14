package eth

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math/big"
	"reflect"
	"strings"

	"go.uber.org/zap"

	"golang.org/x/crypto/sha3"
)

type buffer []byte

func (b buffer) String() string {
	return hex.EncodeToString([]byte(b))
}

type Encoder struct {
	buffer []byte
}

func NewEncoder() *Encoder {
	return &Encoder{}
}

func (e *Encoder) String() string {
	return hex.EncodeToString(e.buffer)
}

func (e *Encoder) Buffer() []byte {
	return e.buffer
}

func (e *Encoder) WriteMethod(method *MethodCall) error {
	methodSignature := method.MethodDef.Signature()
	err := e.Write("method", methodSignature)
	if err != nil {
		return fmt.Errorf("unable to write method in buffer: %w", err)
	}

	if traceEnabled {
		zlog.Debug("written method name in buffer",
			zap.Stringer("buf", buffer(e.buffer)),
			zap.String("method_name", methodSignature),
		)
	}

	type arrayToInsert struct {
		buffOffset uint64
		typeName   string
		value      interface{}
	}

	slicesToInsert := []arrayToInsert{}
	for idx, param := range method.MethodDef.Parameters {

		isAnArray, _ := isArray(param.TypeName)
		if isAnArray {
			slicesToInsert = append(slicesToInsert, arrayToInsert{
				buffOffset: uint64(len(e.buffer)),
				typeName:   param.TypeName,
				value:      method.Data[idx],
			})

			if err := e.Write("uint64", uint64(0)); err != nil {
				return fmt.Errorf("unable to write slice placeholder: %w", err)
			}

			if traceEnabled {
				zlog.Debug("written slice placeholder in buffer",
					zap.Stringer("buf", buffer(e.buffer)),
					zap.String("input_type", param.TypeName),
					zap.Int("input_idx", idx),
				)
			}

			continue
		}

		if err := e.Write(param.TypeName, method.Data[idx]); err != nil {
			return fmt.Errorf("unable to write input.%d %q in buffer: %w", idx, param.TypeName, err)
		}

		if traceEnabled {
			zlog.Debug("written input data in buffer",
				zap.Stringer("buf", buffer(e.buffer)),
				zap.String("input_type", param.TypeName),
				zap.Int("input_idx", idx),
			)
		}
	}

	for sidx, slc := range slicesToInsert {
		// offset should not include the signatures' bytes
		dataLength := uint64(len(e.buffer)) - 4
		d, err := e.encodeUint64(dataLength)
		if err != nil {
			return fmt.Errorf("unable to encode slice offset: %w", err)
		}

		err = e.override(slc.buffOffset, d)
		if err != nil {
			return fmt.Errorf("unable to insert slice offset in buffer: %w", err)
		}

		if traceEnabled {
			zlog.Debug("inserted slice offset in buffer",
				zap.Stringer("buf", buffer(e.buffer)),
				zap.String("input_type", slc.typeName),
				zap.Int("slice_idx", sidx),
			)
		}

		err = e.Write(slc.typeName, slc.value)
		if err != nil {
			return fmt.Errorf("unable to write slice in buffer: %w", err)
		}

		if traceEnabled {
			zlog.Debug("inserted slice in buffer",
				zap.Stringer("buf", buffer(e.buffer)),
				zap.String("input_type", slc.typeName),
				zap.Int("slice_idx", sidx),
			)
		}
	}
	return nil
}

func (e *Encoder) Write(typeName string, in interface{}) error {
	var isAnArray bool
	isAnArray, resolvedTypeName := isArray(typeName)
	if !isAnArray {
		return e.write(resolvedTypeName, in)
	}

	// byte array optimization
	if typeName == "bytes" {
		data := in.([]byte)
		// TOFIX: is this assumption good?
		err := e.write("uint64", uint64(len(data)))
		if err != nil {
			return fmt.Errorf("cannot write byte array length: %w", err)
		}

		err = e.write("bytes", data)
		if err != nil {
			return fmt.Errorf("cannot write byte array: %w", err)
		}
		return nil
	}

	s := reflect.ValueOf(in)
	switch s.Kind() {
	case reflect.Slice:
		// TOFIX: is this assumption good?
		err := e.write("uint64", uint64(s.Len()))
		if err != nil {
			return fmt.Errorf("cannot write slice %s size: %w", typeName, err)
		}

		for i := 0; i < s.Len(); i++ {
			err := e.write(resolvedTypeName, s.Index(i).Interface())
			if err != nil {
				return fmt.Errorf("cannot write item from slice %s.%d: %w", typeName, i, err)
			}
		}
		return nil
	}
	return fmt.Errorf("type %q is not handled right now", resolvedTypeName)
}

func (e *Encoder) write(typeName string, in interface{}) error {
	var d []byte
	var err error
	switch typeName {
	case "address":
		d, err = e.encodeAddress(in.(Address))
	case "uint64":
		d, err = e.encodeUint64(in.(uint64))
	case "uint112":
		d, err = e.encodeBigInt(in.(*big.Int))
	case "uint256":
		d, err = e.encodeBigInt(in.(*big.Int))
	case "bool":
		d, err = e.encodeBool(in.(bool))
	case "bytes":
		d = in.([]byte)
	case "method":
		d, err = e.encodeMethod(in.(string))
	case "event":
		d, err = e.encodeEvent(in.(string))
	default:
		return fmt.Errorf("type %q is not handled right now", typeName)
	}

	if err != nil {
		return err
	}

	e.buffer = append(e.buffer, d...)
	return nil
}

func (e *Encoder) encodeUint64(input uint64) ([]byte, error) {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, input)
	return pad(buf), nil
}

func (e *Encoder) encodeBigInt(input *big.Int) ([]byte, error) {
	return pad(input.Bytes()), nil
}

func (e *Encoder) encodeBool(input bool) ([]byte, error) {
	var v *big.Int
	if input {
		v = big.NewInt(1)
	} else {
		v = big.NewInt(0)
	}
	return pad(v.Bytes()), nil
}

func (e *Encoder) encodeAddress(input Address) ([]byte, error) {
	return pad(input), nil
}

func (e *Encoder) encodeMethod(input string) ([]byte, error) {
	kec := sha3.NewLegacyKeccak256()
	_, err := kec.Write([]byte(input))
	if err != nil {
		return nil, err
	}
	return kec.Sum(nil)[0:4], nil
}

func (e *Encoder) encodeEvent(input string) ([]byte, error) {
	kec := sha3.NewLegacyKeccak256()
	_, err := kec.Write([]byte(input))
	if err != nil {
		return nil, err
	}
	return kec.Sum(nil), nil
}

func (e *Encoder) override(offset uint64, data []byte) error {
	if uint64(len(e.buffer)) < offset+uint64(len(data)) {
		return fmt.Errorf("insuficient room in buffer with length %d to insert data with length %d at offset %d", len(e.buffer), len(data), offset)
	}

	for i := 0; i < len(data); i++ {
		e.buffer[uint64(i)+offset] = data[i]
	}
	return nil
}

func pad(in []byte) []byte {
	d := make([]byte, 32)
	offset := 32 - len(in)
	for i := 0; i < len(in); i++ {
		d[i+offset] = in[i]
	}
	return d
}

func isArray(typeName string) (bool, string) {
	check := strings.HasSuffix(typeName, "[]")
	if check {
		return true, strings.TrimRight(typeName, "[]")
	}

	if typeName == "bytes" {
		return true, "bytes"
	}
	return false, typeName
}
