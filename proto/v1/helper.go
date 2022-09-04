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

	return dblXmm6 + dblXmm4 + dblXmm3 + dblXmm1
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
