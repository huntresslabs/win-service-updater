package updater

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"
)

const (
	INT_UDT_NUMBER_OF_REGISTRY_CHANGES       = 0x20
	INT_UDT_NUMBER_OF_FILE_INFOS             = 0x21 // (precedes file info list)
	STRING_UDT_SERVICE_TO_STOP_BEFORE_UPDATE = 0x32
	STRING_UDT_SERVICE_TO_START_AFTER_UPDATE = 0x33
	END_UDT                                  = 0xFF
)

var UdtTags = map[uint8]string{
	INT_UDT_NUMBER_OF_REGISTRY_CHANGES:       "INT_UDT_NUMBER_OF_REGISTRY_CHANGES",
	INT_UDT_NUMBER_OF_FILE_INFOS:             "INT_UDT_NUMBER_OF_FILE_INFOS", // (precedes file info list)
	STRING_UDT_SERVICE_TO_STOP_BEFORE_UPDATE: "STRING_UDT_SERVICE_TO_STOP_BEFORE_UPDATE",
	STRING_UDT_SERVICE_TO_START_AFTER_UPDATE: "STRING_UDT_SERVICE_TO_START_AFTER_UPDATE",
	END_UDT: "END_UDT",
}

type ConfigUDT struct {
	ServiceToStopBeforeUpdate []TLV
	ServiceToStartAfterUpdate []TLV
}

func ReadUdtTLV(r io.Reader) *TLV {
	var record TLV

	err := binary.Read(r, binary.BigEndian, &record.Tag)
	if err == io.EOF {
		return nil
	} else if err != nil {
		fmt.Println("\n[!] error reading TLV tag:", err.Error())
		return nil
	}

	if record.Tag == END_UDT {
		return nil
	}

	err = binary.Read(r, binary.LittleEndian, &record.Length)
	if err != nil {
		fmt.Println("\n[!] error reading TLV length:", err.Error())
		return nil
	}
	fmt.Println("[+] value length: ", record.Length)

	record.Value = make([]byte, record.Length)
	_, err = io.ReadFull(r, record.Value)
	if err != nil {
		fmt.Println("[!] error reading TLV value:", err.Error())
		return nil
	}

	fmt.Println("[+] read TLV record")
	return &record
}

func ParseUdt(path string) {
	f, err := os.Open(path)
	if nil != err {
		log.Fatal(err)
	}
	defer f.Close()

	// read HEADER
	b := make([]byte, 7)
	f.Read(b)
	fmt.Printf("[+] HEADER %s\n", b)

	var udt ConfigUDT

	for {
		tlv := ReadUdtTLV(f)
		if tlv == nil {
			break
		}

		switch tlv.Tag {
		case STRING_UDT_SERVICE_TO_STOP_BEFORE_UPDATE:
			udt.ServiceToStopBeforeUpdate = append(udt.ServiceToStopBeforeUpdate, *tlv)
		case STRING_UDT_SERVICE_TO_START_AFTER_UPDATE:
			udt.ServiceToStartAfterUpdate = append(udt.ServiceToStartAfterUpdate, *tlv)
		}
	}
}
