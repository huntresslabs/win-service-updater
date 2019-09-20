package updater

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
)

const (
	TLV_BOOL = iota
	TLV_BYTE
	TLV_DSTRING
	TLV_INT
	TLV_LONG
	TLV_STRING
)

type TLV struct {
	Tag        uint8
	TagString  string
	Type       int
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

func ValueToBool(tlv *TLV) []byte {
	return tlv.Value
}

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

func IntValueToBytes(tlv *TLV) []byte {
	return tlv.Value
}

func WriteTLV(f *os.File, tlv TLV) (err error) {
	if tlv.Length == 0 {
		// this tag is not needed
		return nil
	}

	err = binary.Write(f, binary.BigEndian, tlv.Tag)
	if nil != err {
		return err
	}

	if tlv.DataLength > 0 {
		err = binary.Write(f, binary.LittleEndian, tlv.DataLength)
		if nil != err {
			return err
		}
	}

	err = binary.Write(f, binary.LittleEndian, tlv.Length)
	if nil != err {
		return err
	}

	err = binary.Write(f, binary.BigEndian, tlv.Value)
	if nil != err {
		return err
	}

	return nil
}

func displayTagString(tlv *TLV) {
	fmt.Println("[+] String record:", string(tlv.Value))
}

func displayTagUint16(tlv *TLV) {
	buf := bytes.NewBuffer(tlv.Value)
	var value uint16

	err := binary.Read(buf, binary.BigEndian, &value)
	if err != nil {
		fmt.Println("[!] Invalid record:", err.Error())
	} else {
		fmt.Println("[+] Uint16 record:", value)
	}
}

// DisplayTLV will print a TLV for debugging
func DisplayTLV(tlv *TLV) {
	fmt.Printf("[+] %s (%x)\n", tlv.TagString, tlv.Tag)
	switch tlv.Type {
	case TLV_BOOL:
		fmt.Println("   -", ValueToBool(tlv))
	case TLV_BYTE:
		fmt.Println("   -", ValueToByteSlice(tlv))
	case TLV_DSTRING:
		fmt.Println("   -", ValueToString(tlv))
	case TLV_INT:
		fmt.Println("   -", ValueToInt(tlv))
	case TLV_LONG:
		fmt.Println("   -", ValueToLong(tlv))
	case TLV_STRING:
		fmt.Println("   -", ValueToString(tlv))
	default:
		fmt.Println("[!] tlv type", tlv.Type)
	}
}
