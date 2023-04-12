package v1

import (
	"encoding/binary"
	"io"
	"math"
	"strconv"
	"strings"
)

func ParseVolume(ivol int) float64 {
	logpoint := ivol >> (8 * 3)
	hleax := (ivol >> (8 * 2)) & 0xff // [2]
	lheax := (ivol >> 8) & 0xff       // [1]
	lleax := ivol & 0xff              // [0]

	dwEcx := logpoint*2 - 0x7f
	dwEdx := logpoint*2 - 0x86
	dwEsi := logpoint*2 - 0x8e
	dwEax := logpoint*2 - 0x96

	var tmpEax int
	if dwEcx < 0 {
		tmpEax = -dwEcx
	} else {
		tmpEax = dwEcx
	}

	var dbl_xmm6 float64
	dbl_xmm6 = math.Pow(2.0, float64(tmpEax))
	if dwEcx < 0 {
		dbl_xmm6 = 1.0 / dbl_xmm6
	}

	var dbl_xmm4 float64
	if hleax > 0x80 {
		var dwtmpeax = float64(dwEdx + 1)
		tmpdbl_xmm3 := math.Pow(2.0, dwtmpeax)
		dbl_xmm0 := math.Pow(2.0, float64(dwEdx)) * 128.0
		dbl_xmm0 += float64(hleax&0x7f) * tmpdbl_xmm3
		dbl_xmm4 = dbl_xmm0
	} else {
		var dbl_xmm0 float64
		if dwEdx >= 0 {
			dbl_xmm0 = math.Pow(2.0, float64(dwEdx)) * float64(hleax)
		} else {
			dbl_xmm0 = (1 / math.Pow(2.0, float64(dwEdx))) * float64(hleax)
		}
		dbl_xmm4 = dbl_xmm0
	}

	dbl_xmm3 := math.Pow(2.0, float64(dwEsi)) * float64(lheax)
	dbl_xmm1 := math.Pow(2.0, float64(dwEax)) * float64(lleax)
	if hleax&0x80 != 0 {
		dbl_xmm3 *= 2.0
		dbl_xmm1 *= 2.0
	}

	return dbl_xmm6 + dbl_xmm4 + dbl_xmm3 + dbl_xmm1
}

type Struct struct{}

func (s *Struct) Unpark(format string, r io.Reader) interface{} {
	switch format {
	case "h", "H":
		data := make([]byte, 2)
		r.Read(data)
		return binary.LittleEndian.Uint16(data)
	case "b", "B":
		data := make([]byte, 1)
		r.Read(data)
		return uint8(data[0])
	case "i", "I", "l", "L":
		data := make([]byte, 4)
		r.Read(data)
		return binary.LittleEndian.Uint32(data)
	case "q", "Q":
		data := make([]byte, 8)
		r.Read(data)
		return binary.LittleEndian.Uint64(data)
	default:
		if strings.Contains(format, "s") {
			n, _ := strconv.Atoi(strings.TrimRight(format, "s"))
			data := make([]byte, n)
			r.Read(data)
			return string(data)
		}
	}
	return ""
}

var structInstance *Struct

func GetStruct() *Struct {
	if structInstance == nil {
		structInstance = &Struct{}
	}
	return structInstance
}
