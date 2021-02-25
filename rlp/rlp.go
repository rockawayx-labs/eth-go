// Copyrights https://github.com/hyperledger/burrow/tree/master/encoding/rlp
//
// Apache License
// Version 2.0, January 2004
// http://www.apache.org/licenses/

// TERMS AND CONDITIONS FOR USE, REPRODUCTION, AND DISTRIBUTION

// 1. Definitions.

// "License" shall mean the terms and conditions for use, reproduction,
// and distribution as defined by Sections 1 through 9 of this document.

// "Licensor" shall mean the copyright owner or entity authorized by
// the copyright owner that is granting the License.

// "Legal Entity" shall mean the union of the acting entity and all
// other entities that control, are controlled by, or are under common
// control with that entity. For the purposes of this definition,
// "control" means (i) the power, direct or indirect, to cause the
// direction or management of such entity, whether by contract or
// otherwise, or (ii) ownership of fifty percent (50%) or more of the
// outstanding shares, or (iii) beneficial ownership of such entity.

// "You" (or "Your") shall mean an individual or Legal Entity
// exercising permissions granted by this License.

// "Source" form shall mean the preferred form for making modifications,
// including but not limited to software source code, documentation
// source, and configuration files.

// "Object" form shall mean any form resulting from mechanical
// transformation or translation of a Source form, including but
// not limited to compiled object code, generated documentation,
// and conversions to other media types.

// "Work" shall mean the work of authorship, whether in Source or
// Object form, made available under the License, as indicated by a
// copyright notice that is included in or attached to the work
// (an example is provided in the Appendix below).

// "Derivative Works" shall mean any work, whether in Source or Object
// form, that is based on (or derived from) the Work and for which the
// editorial revisions, annotations, elaborations, or other modifications
// represent, as a whole, an original work of authorship. For the purposes
// of this License, Derivative Works shall not include works that remain
// separable from, or merely link (or bind by name) to the interfaces of,
// the Work and Derivative Works thereof.

// "Contribution" shall mean any work of authorship, including
// the original version of the Work and any modifications or additions
// to that Work or Derivative Works thereof, that is intentionally
// submitted to Licensor for inclusion in the Work by the copyright owner
// or by an individual or Legal Entity authorized to submit on behalf of
// the copyright owner. For the purposes of this definition, "submitted"
// means any form of electronic, verbal, or written communication sent
// to the Licensor or its representatives, including but not limited to
// communication on electronic mailing lists, source code control systems,
// and issue tracking systems that are managed by, or on behalf of, the
// Licensor for the purpose of discussing and improving the Work, but
// excluding communication that is conspicuously marked or otherwise
// designated in writing by the copyright owner as "Not a Contribution."

// "Contributor" shall mean Licensor and any individual or Legal Entity
// on behalf of whom a Contribution has been received by Licensor and
// subsequently incorporated within the Work.

// 2. Grant of Copyright License. Subject to the terms and conditions of
// this License, each Contributor hereby grants to You a perpetual,
// worldwide, non-exclusive, no-charge, royalty-free, irrevocable
// copyright license to reproduce, prepare Derivative Works of,
// publicly display, publicly perform, sublicense, and distribute the
// Work and such Derivative Works in Source or Object form.

// 3. Grant of Patent License. Subject to the terms and conditions of
// this License, each Contributor hereby grants to You a perpetual,
// worldwide, non-exclusive, no-charge, royalty-free, irrevocable
// (except as stated in this section) patent license to make, have made,
// use, offer to sell, sell, import, and otherwise transfer the Work,
// where such license applies only to those patent claims licensable
// by such Contributor that are necessarily infringed by their
// Contribution(s) alone or by combination of their Contribution(s)
// with the Work to which such Contribution(s) was submitted. If You
// institute patent litigation against any entity (including a
// cross-claim or counterclaim in a lawsuit) alleging that the Work
// or a Contribution incorporated within the Work constitutes direct
// or contributory patent infringement, then any patent licenses
// granted to You under this License for that Work shall terminate
// as of the date such litigation is filed.

// 4. Redistribution. You may reproduce and distribute copies of the
// Work or Derivative Works thereof in any medium, with or without
// modifications, and in Source or Object form, provided that You
// meet the following conditions:

// (a) You must give any other recipients of the Work or
// Derivative Works a copy of this License; and

// (b) You must cause any modified files to carry prominent notices
// stating that You changed the files; and

// (c) You must retain, in the Source form of any Derivative Works
// that You distribute, all copyright, patent, trademark, and
// attribution notices from the Source form of the Work,
// excluding those notices that do not pertain to any part of
// the Derivative Works; and

// (d) If the Work includes a "NOTICE" text file as part of its
// distribution, then any Derivative Works that You distribute must
// include a readable copy of the attribution notices contained
// within such NOTICE file, excluding those notices that do not
// pertain to any part of the Derivative Works, in at least one
// of the following places: within a NOTICE text file distributed
// as part of the Derivative Works; within the Source form or
// documentation, if provided along with the Derivative Works; or,
// within a display generated by the Derivative Works, if and
// wherever such third-party notices normally appear. The contents
// of the NOTICE file are for informational purposes only and
// do not modify the License. You may add Your own attribution
// notices within Derivative Works that You distribute, alongside
// or as an addendum to the NOTICE text from the Work, provided
// that such additional attribution notices cannot be construed
// as modifying the License.

// You may add Your own copyright statement to Your modifications and
// may provide additional or different license terms and conditions
// for use, reproduction, or distribution of Your modifications, or
// for any such Derivative Works as a whole, provided Your use,
// reproduction, and distribution of the Work otherwise complies with
// the conditions stated in this License.

// 5. Submission of Contributions. Unless You explicitly state otherwise,
// any Contribution intentionally submitted for inclusion in the Work
// by You to the Licensor shall be under the terms and conditions of
// this License, without any additional terms or conditions.
// Notwithstanding the above, nothing herein shall supersede or modify
// the terms of any separate license agreement you may have executed
// with Licensor regarding such Contributions.

// 6. Trademarks. This License does not grant permission to use the trade
// names, trademarks, service marks, or product names of the Licensor,
// except as required for reasonable and customary use in describing the
// origin of the Work and reproducing the content of the NOTICE file.

// 7. Disclaimer of Warranty. Unless required by applicable law or
// agreed to in writing, Licensor provides the Work (and each
// Contributor provides its Contributions) on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
// implied, including, without limitation, any warranties or conditions
// of TITLE, NON-INFRINGEMENT, MERCHANTABILITY, or FITNESS FOR A
// PARTICULAR PURPOSE. You are solely responsible for determining the
// appropriateness of using or redistributing the Work and assume any
// risks associated with Your exercise of permissions under this License.

// 8. Limitation of Liability. In no event and under no legal theory,
// whether in tort (including negligence), contract, or otherwise,
// unless required by applicable law (such as deliberate and grossly
// negligent acts) or agreed to in writing, shall any Contributor be
// liable to You for damages, including any direct, indirect, special,
// incidental, or consequential damages of any character arising as a
// result of this License or out of the use or inability to use the
// Work (including but not limited to damages for loss of goodwill,
// work stoppage, computer failure or malfunction, or any and all
// other commercial damages or losses), even if such Contributor
// has been advised of the possibility of such damages.

// 9. Accepting Warranty or Additional Liability. While redistributing
// the Work or Derivative Works thereof, You may choose to offer,
// and charge a fee for, acceptance of support, warranty, indemnity,
// or other liability obligations and/or rights consistent with this
// License. However, in accepting such obligations, You may act only
// on Your own behalf and on Your sole responsibility, not on behalf
// of any other Contributor, and only if You agree to indemnify,
// defend, and hold each Contributor harmless for any liability
// incurred by, or claims asserted against, such Contributor by reason
// of your accepting any such warranty or additional liability.

// END OF TERMS AND CONDITIONS

// Package rlp is a copy from https://github.com/hyperledger/burrow/tree/a934952b79e4174427385dbe7d8a9f307c05328e/encoding/rlp
// slightly modified to avoid any other Burrow dependency.
//
// We update it from the source file from times to times.
package rlp

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math/big"
	"math/bits"
	"reflect"
)

type magicOffset uint8

const (
	ShortLength              = 55
	StringOffset magicOffset = 0x80 // 128 - if string length is less than or equal to 55 [inclusive]
	SliceOffset  magicOffset = 0xC0 // 192 - if slice length is less than or equal to 55 [inclusive]
	SmallByte                = 0x7f // 247 - value less than or equal is itself [inclusive
)

type Code uint32

const (
	ErrUnknown Code = iota
	ErrNoInput
	ErrInvalid
)

var bigIntType = reflect.TypeOf(&big.Int{})

func (c Code) Error() string {
	switch c {
	case ErrNoInput:
		return "no input"
	case ErrInvalid:
		return "input not valid RLP encoding"
	default:
		return "unknown error"
	}
}

func Encode(input interface{}) ([]byte, error) {
	val := reflect.ValueOf(input)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	return encode(val)
}

func Decode(src []byte, dst interface{}) error {
	fields, err := decode(src)
	if err != nil {
		return err
	}

	val := reflect.ValueOf(dst)
	typ := reflect.TypeOf(dst)

	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	switch val.Kind() {
	case reflect.Slice:
		switch typ.Elem().Kind() {
		case reflect.Uint8:
			out, ok := dst.([]byte)
			if !ok {
				return fmt.Errorf("cannot decode into type %s", val.Type())
			}
			found := bytes.Join(fields, []byte(""))
			if len(out) < len(found) {
				return fmt.Errorf("cannot decode %d bytes into slice of size %d", len(found), len(out))
			}
			for i, b := range found {
				out[i] = b
			}
		default:
			for i := 0; i < val.Len(); i++ {
				elem := val.Index(i)
				err = decodeField(elem, fields[i])
				if err != nil {
					return err
				}
			}
		}
	case reflect.Struct:
		if val.NumField() != len(fields) {
			return fmt.Errorf("wrong number of fields; have %d, want %d", len(fields), val.NumField())
		}
		for i := 0; i < val.NumField(); i++ {
			err := decodeField(val.Field(i), fields[i])
			if err != nil {
				return err
			}
		}
	default:
		return fmt.Errorf("cannot decode into unsupported type %v", reflect.TypeOf(dst))
	}
	return nil
}

func encodeUint8(input uint8) ([]byte, error) {
	if input == 0 {
		// yes this makes no sense, but it does seem to be what everyone else does, apparently 'no leading zeroes'.
		// It means we cannot store []byte{0} because that is indistinguishable from byte{}
		return []byte{uint8(StringOffset)}, nil
	} else if input <= SmallByte {
		return []byte{input}, nil
	} else if input >= uint8(StringOffset) {
		return []byte{0x81, input}, nil
	}
	return []byte{uint8(StringOffset)}, nil
}

func encodeUint64(i uint64) ([]byte, error) {
	size := bits.Len64(i)/8 + 1
	if size == 1 {
		return encodeUint8(uint8(i))
	}
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(i))
	return encodeString(b[8-size:])
}

func encodeBigInt(b *big.Int) ([]byte, error) {
	if b.Sign() == -1 {
		return nil, fmt.Errorf("cannot RLP encode negative number")
	}
	if b.IsUint64() {
		return encodeUint64(b.Uint64())
	}
	bs := b.Bytes()
	length := encodeLength(len(bs), StringOffset)
	return append(length, bs...), nil
}

func encodeLength(n int, offset magicOffset) []byte {
	// > if a string is 0-55 bytes long, the RLP encoding consists of a single byte with value 0x80 plus
	// > the length of the string followed by the string.
	if n <= ShortLength {
		return []uint8{uint8(offset) + uint8(n)}
	}

	i := uint64(n)
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, i)
	byteLengthOfLength := bits.Len64(i)/8 + 1
	// > If a string is more than 55 bytes long, the RLP encoding consists of a single byte with value 0xb7
	// > plus the length in bytes of the length of the string in binary form, followed by the length of the string,
	// > followed by the string
	return append([]byte{uint8(offset) + ShortLength + uint8(byteLengthOfLength)}, b[8-byteLengthOfLength:]...)
}

func encodeString(input []byte) ([]byte, error) {
	if len(input) == 1 && input[0] <= SmallByte {
		return encodeUint8(input[0])
	} else {
		return append(encodeLength(len(input), StringOffset), input...), nil
	}
}

func encodeList(val reflect.Value) ([]byte, error) {
	if val.Len() == 0 {
		return []byte{uint8(SliceOffset)}, nil
	}

	out := make([][]byte, 0)
	for i := 0; i < val.Len(); i++ {
		data, err := encode(val.Index(i))
		if err != nil {
			return nil, err
		}
		out = append(out, data)
	}

	sum := bytes.Join(out, []byte{})
	return append(encodeLength(len(sum), SliceOffset), sum...), nil
}

func encodeStruct(val reflect.Value) ([]byte, error) {
	out := make([][]byte, 0)

	for i := 0; i < val.NumField(); i++ {
		data, err := encode(val.Field(i))
		if err != nil {
			return nil, err
		}
		out = append(out, data)
	}
	sum := bytes.Join(out, []byte{})
	length := encodeLength(len(sum), SliceOffset)
	return append(length, sum...), nil
}

func encode(val reflect.Value) ([]byte, error) {
	if val.Kind() == reflect.Interface {
		val = val.Elem()
	}

	switch val.Kind() {
	case reflect.Ptr:
		if !val.Type().AssignableTo(bigIntType) {
			return nil, fmt.Errorf("cannot encode pointer type %v", val.Type())
		}
		return encodeBigInt(val.Interface().(*big.Int))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i := val.Int()
		if i < 0 {
			return nil, fmt.Errorf("cannot rlp encode negative integer")
		}
		return encodeUint64(uint64(i))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return encodeUint64(val.Uint())
	case reflect.Bool:
		if val.Bool() {
			return []byte{0x01}, nil
		}
		return []byte{uint8(StringOffset)}, nil
	case reflect.String:
		return encodeString([]byte(val.String()))
	case reflect.Slice:
		switch val.Type().Elem().Kind() {
		case reflect.Uint8:
			i, err := encodeString(val.Bytes())
			return i, err
		default:
			return encodeList(val)
		}
	case reflect.Struct:
		return encodeStruct(val)
	default:
		return []byte{uint8(StringOffset)}, nil
	}
}

// Split into RLP fields by reading length prefixes and consuming chunks
func decode(in []byte) ([][]byte, error) {
	if len(in) == 0 {
		return nil, nil
	}

	offset, length, typ := decodeLength(in)
	end := offset + length

	if end > uint64(len(in)) {
		return nil, fmt.Errorf("read length prefix of %d but there is only %d bytes of unconsumed input",
			length, uint64(len(in))-offset)
	}

	suffix, err := decode(in[end:])
	if err != nil {
		return nil, err
	}
	switch typ {
	case reflect.String:
		return append([][]byte{in[offset:end]}, suffix...), nil
	case reflect.Slice:
		prefix, err := decode(in[offset:end])
		if err != nil {
			return nil, err
		}
		return append(prefix, suffix...), nil
	}

	return suffix, nil
}

func decodeLength(input []byte) (uint64, uint64, reflect.Kind) {
	magicByte := magicOffset(input[0])

	switch {
	case magicByte <= SmallByte:
		// small byte: sufficiently small single byte
		return 0, 1, reflect.String

	case magicByte <= StringOffset+ShortLength:
		// short string: length less than or equal to 55 bytes
		length := uint64(magicByte - StringOffset)
		return 1, length, reflect.String

	case magicByte < SliceOffset:
		// long string: length described by magic = 0xb7 + <byte length of length of string>
		byteLengthOfLength := magicByte - StringOffset - ShortLength
		length := getUint64(input[1:byteLengthOfLength])
		offset := uint64(byteLengthOfLength + 1)
		return offset, length, reflect.String

	case magicByte <= SliceOffset+ShortLength:
		// short slice: length less than or equal to 55 bytes
		length := uint64(magicByte - SliceOffset)
		return 1, length, reflect.Slice

	// Note this takes us all the way up to <= 255 so this switch is exhaustive
	default:
		// long string: length described by magic = 0xf7 + <byte length of length of string>
		byteLengthOfLength := magicByte - SliceOffset - ShortLength
		length := getUint64(input[1:byteLengthOfLength])
		offset := uint64(byteLengthOfLength + 1)
		return offset, length, reflect.Slice
	}
}

func getUint64(bs []byte) uint64 {
	bs = leftPadBytes(bs, 8)
	return binary.BigEndian.Uint64(bs)
}

func decodeField(val reflect.Value, field []byte) error {
	typ := val.Type()

	switch val.Kind() {
	case reflect.Ptr:
		if !typ.AssignableTo(bigIntType) {
			return fmt.Errorf("cannot decode into pointer type %v", typ)
		}
		bi := new(big.Int).SetBytes(field)
		val.Set(reflect.ValueOf(bi))

	case reflect.String:
		val.SetString(string(field))
	case reflect.Uint64:
		out := make([]byte, 8)
		for j := range field {
			out[len(out)-(len(field)-j)] = field[j]
		}
		val.SetUint(binary.BigEndian.Uint64(out))
	case reflect.Slice:
		if typ.Elem().Kind() != reflect.Uint8 {
			// skip
			return nil
		}
		out := make([]byte, len(field))
		for i, b := range field {
			out[i] = b
		}
		val.SetBytes(out)
	}
	return nil
}

func leftPadBytes(slice []byte, l int) []byte {
	if l < len(slice) {
		return slice
	}
	padded := make([]byte, l)
	copy(padded[l-len(slice):], slice)
	return padded
}
