package updater

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

const (
	DSTRING_IUC_COMPANY_NAME           = 0x01
	DSTRING_IUC_PRODUCT_NAME           = 0x02
	DSTRING_IUC_INSTALLED_VERSION      = 0x03
	STRING_IUC_GUID                    = 0x0A
	DSTRING_IUC_SERVER_FILE_SITE       = 0x04
	DSTRING_IUC_WYUPDATE_SERVER_SITE   = 0x09
	DSTRING_IUC_HEADER_IMAGE_ALIGNMENT = 0x11
	INT_IUC_HEADER_TEXT_INDENT         = 0x12
	DSTRING_IUC_HEADER_TEXT_COLOR      = 0x13
	DSTRING_IUC_HEADER_FILENAME        = 0x14
	DSTRING_IUC_SIDE_IMAGE_FILENAME    = 0x15
	DSTRING_IUC_LANGUAGE_CULTURE       = 0x18 // e.g., en-US
	DSTRING_IUC_LANGUAGE_FILENAME      = 0x16
	BOOL_IUC_HIDE_HEADER_DIVIDER       = 0x17
	BOOL_IUC_CLOSE_WYUPDATE            = 0x19
	STRING_IUC_CUSTOM_TITLE_BAR        = 0x1A
	STRING_IUC_PUBLIC_KEY              = 0x1B
	END_IUC                            = 0xFF
)

var tags = map[uint8]string{
	BOOL_IUC_CLOSE_WYUPDATE:            "BOOL_IUC_CLOSE_WYUPDATE",
	BOOL_IUC_HIDE_HEADER_DIVIDER:       "BOOL_IUC_HIDE_HEADER_DIVIDER",
	DSTRING_IUC_COMPANY_NAME:           "DSTRING_IUC_COMPANY_NAME",
	DSTRING_IUC_HEADER_FILENAME:        "DSTRING_IUC_HEADER_FILENAME",
	DSTRING_IUC_HEADER_IMAGE_ALIGNMENT: "DSTRING_IUC_HEADER_IMAGE_ALIGNMENT",
	DSTRING_IUC_HEADER_TEXT_COLOR:      "DSTRING_IUC_HEADER_TEXT_COLOR",
	DSTRING_IUC_INSTALLED_VERSION:      "DSTRING_IUC_INSTALLED_VERSION",
	DSTRING_IUC_LANGUAGE_CULTURE:       "DSTRING_IUC_LANGUAGE_CULTURE", // e.g., en-US
	DSTRING_IUC_LANGUAGE_FILENAME:      "DSTRING_IUC_LANGUAGE_FILENAME",
	DSTRING_IUC_PRODUCT_NAME:           "DSTRING_IUC_PRODUCT_NAME",
	DSTRING_IUC_SERVER_FILE_SITE:       "DSTRING_IUC_SERVER_FILE_SITE",
	DSTRING_IUC_SIDE_IMAGE_FILENAME:    "DSTRING_IUC_SIDE_IMAGE_FILENAME",
	DSTRING_IUC_WYUPDATE_SERVER_SITE:   "DSTRING_IUC_WYUPDATE_SERVER_SITE",
	INT_IUC_HEADER_TEXT_INDENT:         "INT_IUC_HEADER_TEXT_INDENT",
	STRING_IUC_CUSTOM_TITLE_BAR:        "STRING_IUC_CUSTOM_TITLE_BAR",
	STRING_IUC_GUID:                    "STRING_IUC_GUID",
	STRING_IUC_PUBLIC_KEY:              "STRING_IUC_PUBLIC_KEY",
	END_IUC:                            "END_IUC",
}

type ConfigIUC struct {
	IucCompanyName          TLV
	IucProductName          TLV
	IucInstalledVersion     TLV
	IucGUID                 TLV
	IucServerFileSite       []TLV
	IucWyupdateServerSite   []TLV
	IucHeaderImageAlignment TLV
	IucHeaderTextIndent     TLV
	IucHeaderTextColor      TLV
	IucHeaderFilename       TLV
	IucSideImageFilename    TLV
	IucLanguageCulture      TLV
	IucLanguageFilename     TLV
	IucHideHeaderDivider    TLV
	IucCloseWyupate         TLV
	IucCustomTitleBar       TLV
	IucPublicKey            TLV
}

func ReadWycTLV(r io.Reader) *TLV {
	var record TLV

	err := binary.Read(r, binary.BigEndian, &record.Tag)
	if err == io.EOF {
		return nil
	} else if err != nil {
		// fmt.Println("\n[!] error reading TLV tag:", err.Error())
		return nil
	}

	// fmt.Printf("- %s (%x)\n", tags[record.Tag], record.Tag)

	if record.Tag == END_IUC {
		return nil
	}

	switch record.Tag {
	case DSTRING_IUC_COMPANY_NAME,
		DSTRING_IUC_PRODUCT_NAME,
		DSTRING_IUC_INSTALLED_VERSION,
		DSTRING_IUC_SERVER_FILE_SITE,
		DSTRING_IUC_WYUPDATE_SERVER_SITE,
		DSTRING_IUC_HEADER_IMAGE_ALIGNMENT,
		DSTRING_IUC_HEADER_TEXT_COLOR,
		DSTRING_IUC_HEADER_FILENAME,
		DSTRING_IUC_SIDE_IMAGE_FILENAME,
		DSTRING_IUC_LANGUAGE_CULTURE,
		DSTRING_IUC_LANGUAGE_FILENAME:
		err = binary.Read(r, binary.LittleEndian, &record.DataLength)
		if err != nil {
			// fmt.Println("\n[!] error reading TLV data length:", err.Error())
			return nil
		}
		// fmt.Println("[+] value data length: ", record.DataLength)
	default:
		// fmt.Println("[!] unknown tag", record.Tag)
	}

	err = binary.Read(r, binary.LittleEndian, &record.Length)
	if err != nil {
		// fmt.Println("\n[!] error reading TLV length:", err.Error())
		return nil
	}
	// fmt.Println("[+] value length: ", record.Length)

	record.Value = make([]byte, record.Length)
	_, err = io.ReadFull(r, record.Value)
	if err != nil {
		// fmt.Println("[!] error reading TLV value:", err.Error())
		return nil
	}

	// fmt.Println("[+] read TLV record")
	return &record
}

func ParseWyc(path string) (ConfigIUC, error) {
	var config ConfigIUC

	f, err := os.Open(path)
	if nil != err {
		return config, err
	}
	defer f.Close()

	// read HEADER
	b := make([]byte, 7)
	f.Read(b)

	for {
		tlv := ReadWycTLV(f)
		if tlv == nil {
			break
		}

		switch tlv.Tag {
		case DSTRING_IUC_COMPANY_NAME:
			config.IucCompanyName = *tlv
		case DSTRING_IUC_PRODUCT_NAME:
			config.IucProductName = *tlv
		case DSTRING_IUC_INSTALLED_VERSION:
			config.IucInstalledVersion = *tlv
		case DSTRING_IUC_SERVER_FILE_SITE:
			config.IucServerFileSite = append(config.IucServerFileSite, *tlv)
		case DSTRING_IUC_WYUPDATE_SERVER_SITE:
			config.IucWyupdateServerSite = append(config.IucWyupdateServerSite, *tlv)
		case DSTRING_IUC_HEADER_IMAGE_ALIGNMENT:
			config.IucHeaderImageAlignment = *tlv
		case DSTRING_IUC_HEADER_TEXT_COLOR:
			config.IucHeaderTextColor = *tlv
		case DSTRING_IUC_HEADER_FILENAME:
			config.IucHeaderFilename = *tlv
		case DSTRING_IUC_SIDE_IMAGE_FILENAME:
			config.IucSideImageFilename = *tlv
		case DSTRING_IUC_LANGUAGE_CULTURE:
			config.IucLanguageCulture = *tlv
		case DSTRING_IUC_LANGUAGE_FILENAME:
			config.IucLanguageFilename = *tlv
		case INT_IUC_HEADER_TEXT_INDENT:
			config.IucHeaderTextIndent = *tlv
		case BOOL_IUC_HIDE_HEADER_DIVIDER:
			config.IucHideHeaderDivider = *tlv
		case BOOL_IUC_CLOSE_WYUPDATE:
			config.IucCloseWyupate = *tlv
		case STRING_IUC_CUSTOM_TITLE_BAR:
			config.IucCustomTitleBar = *tlv
		case STRING_IUC_PUBLIC_KEY:
			config.IucPublicKey = *tlv
		case STRING_IUC_GUID:
			config.IucGUID = *tlv
		default:
			err = fmt.Errorf("crap")
			return config, err
		}
	}
	return config, nil
}

func WriteWYC(config ConfigIUC, path string) error {
	f, err := os.Create(path)
	if nil != err {
		return err
	}
	defer f.Close()

	// write HEADER
	f.Write([]byte("IUCDFV2"))

	// DSTRING_IUC_SIDE_IMAGE_FILENAME:
	WriteTLV(f, config.IucSideImageFilename)

	// DSTRING_IUC_HEADER_TEXT_COLOR:
	WriteTLV(f, config.IucHeaderTextColor)

	// STRING_IUC_GUID:
	WriteTLV(f, config.IucGUID)

	// DSTRING_IUC_SERVER_FILE_SITE:
	for _, s := range config.IucServerFileSite {
		WriteTLV(f, s)
	}

	// DSTRING_IUC_COMPANY_NAME:
	WriteTLV(f, config.IucCompanyName)

	// STRING_IUC_PUBLIC_KEY:
	WriteTLV(f, config.IucPublicKey)

	// DSTRING_IUC_LANGUAGE_CULTURE:
	WriteTLV(f, config.IucLanguageCulture)

	// DSTRING_IUC_PRODUCT_NAME:
	WriteTLV(f, config.IucProductName)

	// INT_IUC_HEADER_TEXT_INDENT:
	WriteTLV(f, config.IucHeaderTextIndent)

	// DSTRING_IUC_HEADER_FILENAME:
	WriteTLV(f, config.IucHeaderFilename)

	// DSTRING_IUC_INSTALLED_VERSION:
	WriteTLV(f, config.IucInstalledVersion)

	// BOOL_IUC_HIDE_HEADER_DIVIDER:
	WriteTLV(f, config.IucHideHeaderDivider)

	// DSTRING_IUC_HEADER_IMAGE_ALIGNMENT:
	WriteTLV(f, config.IucHeaderImageAlignment)

	// DSTRING_IUC_LANGUAGE_FILENAME - NOT USED
	WriteTLV(f, config.IucLanguageFilename)

	// DSTRING_IUC_WYUPDATE_SERVER_SITE - NOT USED
	for _, s := range config.IucWyupdateServerSite {
		WriteTLV(f, s)
	}

	// STRING_IUC_CUSTOM_TITLE_BAR - NOT USED
	WriteTLV(f, config.IucCustomTitleBar)

	// BOOL_IUC_CLOSE_WYUPDATE:
	WriteTLV(f, config.IucCloseWyupate)

	err = binary.Write(f, binary.BigEndian, byte(END_IUC))
	if nil != err {
		return err
	}

	// added because the test files created with python have
	// 0x00 appended
	// for i := 0; i < 3; i++ {
	// 	err = binary.Write(f, binary.BigEndian, byte(0x00))
	// 	if nil != err {
	// 		return err
	// 	}
	// }

	return nil
}
