package updater

import (
	"encoding/binary"
)

type TLV struct {
	Tag        uint8
	TagString  string
	DataLength uint32
	Length     uint32
	Value      []byte
}

// wyUpdate types
// int – 32‐bit, little‐endian, signed integer
// long – 64‐bit, little‐endian, signed integer

// d. string
// Int that stores String Length ‘N’ + 4,
// Int that stores String Length ‘N’,
// UTF8 string N bytes long

// string
// Int that stores String Length ‘N’,
// UTF8 string N bytes long

func ValueToInt(tlv *TLV) int {
	return int(binary.LittleEndian.Uint32(tlv.Value))
}

func ValueToLong(tlv *TLV) int64 {
	return int64(binary.LittleEndian.Uint64(tlv.Value))
}

func ValueToByteSlice(tlv *TLV) []byte {
	return tlv.Value
}

func ValueToString(tlv *TLV) string {
	return string(tlv.Value)
}

// type TLVer interface {
// 	String() string
// 	Int() int
// }

// func (tlv *TLV) Int() int {
// 	return int(binary.LittleEndian.Uint32(tlv.Value))
// }

// func (tlv *TLV) String() string {
// 	return string(tlv.Value)
// }
