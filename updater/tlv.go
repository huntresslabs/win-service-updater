package updater

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type TLV struct {
	Tag        uint8
	DataLength uint32
	Length     uint32
	Value      []byte
}

func displayTagString(tlv *TLV) {
	fmt.Println("[+] String record:", string(tlv.Value))
}

// func displayTagUint16(tlv *TLV) {
// 	buf := bytes.NewBuffer(tlv.Value)
// 	var value uint16

// 	err := binary.Read(buf, binary.BigEndian, &value)
// 	if err != nil {
// 		fmt.Println("[!] Invalid record:", err.Error())
// 	} else {
// 		fmt.Println("[+] Uint16 record:", value)
// 	}
// }

func displayTagUint32(tlv *TLV) {
	buf := bytes.NewBuffer(tlv.Value)
	var value uint32

	err := binary.Read(buf, binary.LittleEndian, &value)
	if err != nil {
		fmt.Println("[!] Invalid record:", err.Error())
	} else {
		fmt.Println("[+] Uint32 record:", value)
	}
}

// func displayTLV(tlv *TLV) {
// 	fmt.Printf("[+] tag %s (%x)\n", tags[tlv.Tag], tlv.Tag)
// 	switch tlv.Tag {
// 	case DSTRING_IUC_COMPANY_NAME,
// 		DSTRING_IUC_PRODUCT_NAME,
// 		DSTRING_IUC_INSTALLED_VERSION,
// 		DSTRING_IUC_SERVER_FILE_SITE,
// 		DSTRING_IUC_WYUPDATE_SERVER_SITE,
// 		DSTRING_IUC_HEADER_IMAGE_ALIGNMENT,
// 		DSTRING_IUC_HEADER_TEXT_COLOR,
// 		DSTRING_IUC_HEADER_FILENAME,
// 		DSTRING_IUC_SIDE_IMAGE_FILENAME,
// 		DSTRING_IUC_LANGUAGE_CULTURE,
// 		DSTRING_IUC_LANGUAGE_FILENAME:
// 		displayTagString(tlv)
// 	case INT_IUC_HEADER_TEXT_INDENT,
// 		BOOL_IUC_HIDE_HEADER_DIVIDER,
// 		BOOL_IUC_CLOSE_WYUPDATE,
// 		INT_UDT_NUMBER_OF_FILE_INFOS,
// 		INT_UDT_NUMBER_OF_REGISTRY_CHANGES:
// 		displayTagUint32(tlv)
// 	case STRING_IUC_CUSTOM_TITLE_BAR,
// 		STRING_IUC_PUBLIC_KEY,
// 		STRING_IUC_GUID,
// 		STRING_UDT_SERVICE_TO_STOP_BEFORE_UPDATE,
// 		STRING_UDT_SERVICE_TO_START_AFTER_UPDATE:
// 		displayTagString(tlv)
// 	case END_IUC:
// 		return
// 	default:
// 		fmt.Println("[!] unknown tag", tlv.Tag)
// 	}
// }
