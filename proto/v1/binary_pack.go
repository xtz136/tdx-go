// Copyright 2017 Roman Kachanovsky. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
	Package binary_pack performs conversions between some Go values represented as byte slices.
	This can be used in handling binary data stored in files or from network connections,
	among other sources. It uses format slices of strings as compact descriptions of the layout
	of the Go structs.

	Format characters (some characters like H have been reserved for future implementation of unsigned numbers):
		? - bool, packed size 1 byte
		h, H - int, packed size 2 bytes (in future it will support pack/unpack of int8, uint8 values)
		i, I, l, L - int, packed size 4 bytes (in future it will support pack/unpack of int16, uint16, int32, uint32 values)
		q, Q - int, packed size 8 bytes (in future it will support pack/unpack of int64, uint64 values)
		f - float32, packed size 4 bytes
		d - float64, packed size 8 bytes
		Ns - string, packed size N bytes, N is a number of runes to pack/unpack

*/
package v1

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
)

type BinaryPack struct{}

// Return a byte slice containing the values of msg slice packed according to the given format.
// The items of msg slice must match the values required by the format exactly.
func (bp *BinaryPack) Pack(format []string, msg []interface{}) ([]byte, error) {
	if len(format) > len(msg) {
		return nil, errors.New("Format is longer than values to pack")
	}

	res := []byte{}

	for i, f := range format {
		switch f {
		case "?":
			casted_value, ok := msg[i].(bool)
			if !ok {
				return nil, errors.New("Type of passed value doesn't match to expected '" + f + "' (bool)")
			}
			res = append(res, boolToBytes(casted_value)...)
		case "h", "H":
			casted_value, ok := msg[i].(int)
			if !ok {
				return nil, errors.New("Type of passed value doesn't match to expected '" + f + "' (int, 2 bytes)")
			}
			res = append(res, intToBytes(casted_value, 2)...)
		case "i", "I", "l", "L":
			casted_value, ok := msg[i].(int)
			if !ok {
				return nil, errors.New("Type of passed value doesn't match to expected '" + f + "' (int, 4 bytes)")
			}
			res = append(res, intToBytes(casted_value, 4)...)
		case "q", "Q":
			casted_value, ok := msg[i].(int)
			if !ok {
				return nil, errors.New("Type of passed value doesn't match to expected '" + f + "' (int, 8 bytes)")
			}
			res = append(res, intToBytes(casted_value, 8)...)
		case "f":
			casted_value, ok := msg[i].(float32)
			if !ok {
				return nil, errors.New("Type of passed value doesn't match to expected '" + f + "' (float32)")
			}
			res = append(res, float32ToBytes(casted_value, 4)...)
		case "d":
			casted_value, ok := msg[i].(float64)
			if !ok {
				return nil, errors.New("Type of passed value doesn't match to expected '" + f + "' (float64)")
			}
			res = append(res, float64ToBytes(casted_value, 8)...)
		default:
			if strings.Contains(f, "s") {
				casted_value, ok := msg[i].(string)
				if !ok {
					return nil, errors.New("Type of passed value doesn't match to expected '" + f + "' (string)")
				}
				n, _ := strconv.Atoi(strings.TrimRight(f, "s"))
				res = append(res, []byte(fmt.Sprintf("%s%s",
					casted_value, strings.Repeat("\x00", n-len(casted_value))))...)
			} else {
				return nil, errors.New("Unexpected format token: '" + f + "'")
			}
		}
	}

	return res, nil
}

// Unpack the byte slice (presumably packed by Pack(format, msg)) according to the given format.
// The result is a []interface{} slice even if it contains exactly one item.
// The byte slice must contain not less the amount of data required by the format
// (len(msg) must more or equal CalcSize(format)).
func (bp *BinaryPack) UnPack(format []string, msg []byte) ([]byte, []interface{}, error) {
	expected_size, err := bp.CalcSize(format)

	if err != nil {
		return nil, nil, err
	}

	if expected_size > len(msg) {
		return nil, nil, errors.New("Expected size is bigger than actual size of message")
	}

	res := []interface{}{}

	for _, f := range format {
		switch f {
		case "?":
			res = append(res, bytesToBool(msg[:1]))
			msg = msg[1:]
		case "h", "H":
			res = append(res, bytesToInt(msg[:2]))
			msg = msg[2:]
		case "b", "B":
			res = append(res, bytesToInt(msg[:1]))
			msg = msg[1:]
		case "i", "I", "l", "L":
			res = append(res, bytesToInt(msg[:4]))
			msg = msg[4:]
		case "q", "Q":
			res = append(res, bytesToInt(msg[:8]))
			msg = msg[8:]
		case "f":
			res = append(res, bytesToFloat32(msg[:4]))
			msg = msg[4:]
		case "d":
			res = append(res, bytesToFloat64(msg[:8]))
			msg = msg[8:]
		default:
			if strings.Contains(f, "s") {
				n, _ := strconv.Atoi(strings.TrimRight(f, "s"))
				res = append(res, string(msg[:n]))
				msg = msg[n:]
			} else {
				return nil, nil, errors.New("Unexpected format token: '" + f + "'")
			}
		}
	}

	return msg, res, nil
}

func (bp *BinaryPack) UnPackPrice(msg []byte) ([]byte, float64) {
	posbyte := 6
	pos := 0
	bdata := msg[pos]
	intdata := int(bdata & 0x3f)
	sign := (bdata & 0x40) != 0

	if (bdata & 0x80) != 0 {
		for {
			pos += 1
			bdata = msg[pos]
			intdata += (int(bdata&0x7f) << posbyte)
			posbyte += 7
			if (bdata & 0x80) == 0 {
				break
			}
		}
	}
	pos += 1
	if sign {
		intdata = -intdata
	}

	return msg[pos:], float64(intdata)
}

func (bp *BinaryPack) UnPackAmount(msg []byte) ([]byte, float64, error) {
	msg, values, err := bp.UnPack([]string{"I"}, msg)
	if err != nil {
		return nil, 0, err
	}

	ivol := values[0].(int)
	logpoint := ivol >> (8 * 3)
	hleax := (ivol >> (8 * 2)) & 0xff
	lheax := (ivol >> 8) & 0xff
	lleax := ivol & 0xff

	dwEcx := logpoint*2 - 0x7f
	dwEdx := logpoint*2 - 0x86
	dwEsi := logpoint*2 - 0x8e
	dwEax := logpoint*2 - 0x96

	tmpEax := 0
	if dwEcx < 0 {
		tmpEax = -dwEcx
	} else {
		tmpEax = dwEcx
	}

	dblXmm6 := math.Pow(2.0, float64(tmpEax))
	if dwEax < 0 {
		dblXmm6 = 1.0 / dblXmm6
	}

	var dblXmm4 float64
	if hleax > 0x80 {
		dwtmpeax := dwEdx + 1
		tmpdblXmm3 := math.Pow(2.0, float64(dwtmpeax))
		dblXmm0 := math.Pow(2.0, float64(dwEdx)) * 128.0
		dblXmm0 += float64(hleax&0x7f) * tmpdblXmm3
		dblXmm4 = dblXmm0
	} else {
		dblXmm0 := 0.0
		if dwEdx >= 0 {
			dblXmm0 = math.Pow(2.0, float64(dwEdx)) * float64(hleax)
		} else {
			dblXmm0 = (1 / math.Pow(2.0, float64(dwEdx))) * float64(hleax)
		}
		dblXmm4 = dblXmm0
	}

	dblXmm3 := math.Pow(2.0, float64(dwEsi)) * float64(lheax)
	dblXmm1 := math.Pow(2.0, float64(dwEax)) * float64(lleax)
	if hleax&0x80 > 0 {
		dblXmm3 *= 2.0
		dblXmm1 *= 2.0
	}

	ret := dblXmm6 + dblXmm4 + dblXmm3 + dblXmm1
	return msg, ret, nil

}

// Return the size of the struct (and hence of the byte slice) corresponding to the given format.
func (bp *BinaryPack) CalcSize(format []string) (int, error) {
	var size int

	for _, f := range format {
		switch f {
		case "?":
			size = size + 1
		case "h", "H":
			size = size + 2
		case "b", "B":
			size = size + 1
		case "i", "I", "l", "L", "f":
			size = size + 4
		case "q", "Q", "d":
			size = size + 8
		default:
			if strings.Contains(f, "s") {
				n, _ := strconv.Atoi(strings.TrimRight(f, "s"))
				size = size + n
			} else {
				return 0, errors.New("Unexpected format token: '" + f + "'")
			}
		}
	}

	return size, nil
}

func boolToBytes(x bool) []byte {
	if x {
		return intToBytes(1, 1)
	}
	return intToBytes(0, 1)
}

func bytesToBool(b []byte) bool {
	return bytesToInt(b) > 0
}

func intToBytes(n int, size int) []byte {
	buf := bytes.NewBuffer([]byte{})
	binary.Write(buf, binary.LittleEndian, int64(n))
	return buf.Bytes()[0:size]
}

func bytesToInt(b []byte) int {
	buf := bytes.NewBuffer(b)

	switch len(b) {
	case 1:
		var x int8
		binary.Read(buf, binary.LittleEndian, &x)
		return int(x)
	case 2:
		var x int16
		binary.Read(buf, binary.LittleEndian, &x)
		return int(x)
	case 4:
		var x int32
		binary.Read(buf, binary.LittleEndian, &x)
		return int(x)
	default:
		var x int64
		binary.Read(buf, binary.LittleEndian, &x)
		return int(x)
	}
}

func float32ToBytes(n float32, size int) []byte {
	buf := bytes.NewBuffer([]byte{})
	binary.Write(buf, binary.LittleEndian, n)
	return buf.Bytes()[0:size]
}

func bytesToFloat32(b []byte) float32 {
	var x float32
	buf := bytes.NewBuffer(b)
	binary.Read(buf, binary.LittleEndian, &x)
	return x
}

func float64ToBytes(n float64, size int) []byte {
	buf := bytes.NewBuffer([]byte{})
	binary.Write(buf, binary.LittleEndian, n)
	return buf.Bytes()[0:size]
}

func bytesToFloat64(b []byte) float64 {
	var x float64
	buf := bytes.NewBuffer(b)
	binary.Read(buf, binary.LittleEndian, &x)
	return x
}
