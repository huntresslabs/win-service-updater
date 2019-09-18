// Parser for wys files
// File ID: IUSDFV2
// Compressed File ID: = { 0x50, 0x4b, 0x03, 0x04 } = { 'P', 'K', 0x03, 0x04 }
// File Extension: wys

package updater

import (
	"archive/zip"
	"encoding/binary"
	"fmt"
	"io"
	"strings"
)

const (
	DSTRING_WYS_CURRENT_LAST_VERSION      = 0x01
	DSTRING_WYS_SERVER_FILE_SITE          = 0x02
	DSTRING_WYS_MIN_CLIENT_VERSION        = 0x07
	INT_WYS_DUMMY_VAR_LEN                 = 0x0F
	DSTRING_WYS_VERSION_TO_UPDATE         = 0x0B
	DSTRING_WYS_UPDATE_FILE_SITE          = 0x03
	BYTE_WYS_RTF                          = 0x80
	DSTRING_WYS_LATEST_CHANGES            = 0x04
	LONG_WYS_UPDATE_FILE_SIZE             = 0x09
	LONG_WYS_UPDATE_FILE_ADLER32_CHECKSUM = 0x08
	BYTE_WYS_FILE_SHA1                    = 0x14
	INT_WYS_FOLDER                        = 0x0A
	DSTRING_WYS_UPDATE_ERROR_TEXT         = 0x20
	DSTRING_WYS_UPDATE_ERROR_LINK         = 0x21
	END_WYS                               = 0xFF
)

var WysTags = map[uint8]string{
	BYTE_WYS_FILE_SHA1:                    "BYTE_WYS_FILE_SHA1",
	BYTE_WYS_RTF:                          "BYTE_WYS_RTF",
	DSTRING_WYS_CURRENT_LAST_VERSION:      "DSTRING_WYS_CURRENT_LAST_VERSION",
	DSTRING_WYS_LATEST_CHANGES:            "DSTRING_WYS_LATEST_CHANGES",
	DSTRING_WYS_MIN_CLIENT_VERSION:        "DSTRING_WYS_MIN_CLIENT_VERSION",
	DSTRING_WYS_SERVER_FILE_SITE:          "DSTRING_WYS_SERVER_FILE_SITE",
	DSTRING_WYS_UPDATE_ERROR_LINK:         "DSTRING_WYS_UPDATE_ERROR_LINK",
	DSTRING_WYS_UPDATE_ERROR_TEXT:         "DSTRING_WYS_UPDATE_ERROR_TEXT",
	DSTRING_WYS_UPDATE_FILE_SITE:          "DSTRING_WYS_UPDATE_FILE_SITE",
	DSTRING_WYS_VERSION_TO_UPDATE:         "DSTRING_WYS_VERSION_TO_UPDATE",
	END_WYS:                               "END_WYS",
	INT_WYS_DUMMY_VAR_LEN:                 "INT_WYS_DUMMY_VAR_LEN",
	INT_WYS_FOLDER:                        "INT_WYS_FOLDER",
	LONG_WYS_UPDATE_FILE_ADLER32_CHECKSUM: "LONG_WYS_UPDATE_FILE_ADLER32_CHECKSUM",
	LONG_WYS_UPDATE_FILE_SIZE:             "LONG_WYS_UPDATE_FILE_SIZE",
}

type ConfigWYS struct {
	FileSha1           []byte
	RTF                []byte
	CurrentLastVersion string
	LatestChanges      string
	MinClientVersion   string
	ServerFileSite     string
	UpdateErrorLink    string
	UpdateErrorText    string
	UpdateFileSite     []string // hosts the WYU file
	VersionToUpdate    string
	DummyVarLen        int
	WysFolder          int
	UpdateFileAdler32  int64
	UpdateFileSize     int64
}

func ReadWysTLV(r io.Reader) *TLV {
	var record TLV

	err := binary.Read(r, binary.BigEndian, &record.Tag)
	if err == io.EOF {
		return nil
	} else if err != nil {
		// fmt.Println("\n[!] error reading TLV tag:", err.Error())
		return nil
	}

	if record.Tag == END_WYS {
		return nil
	}

	record.TagString = WysTags[record.Tag]

	switch record.Tag {
	case DSTRING_WYS_CURRENT_LAST_VERSION,
		DSTRING_WYS_LATEST_CHANGES,
		DSTRING_WYS_MIN_CLIENT_VERSION,
		DSTRING_WYS_SERVER_FILE_SITE,
		DSTRING_WYS_UPDATE_ERROR_LINK,
		DSTRING_WYS_UPDATE_ERROR_TEXT,
		DSTRING_WYS_UPDATE_FILE_SITE,
		DSTRING_WYS_VERSION_TO_UPDATE:
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
	// fmt.Printf("[+] %s (%x) value length: %d\n", WysTags[record.Tag], record.Tag, record.Length)

	// there is no value for the dummy var
	if record.Tag == INT_WYS_DUMMY_VAR_LEN {
		return &record
	}

	record.Value = make([]byte, record.Length)
	_, err = io.ReadFull(r, record.Value)
	if err != nil {
		// fmt.Println("[!] error reading TLV value:", err.Error())
		return nil
	}

	// fmt.Println("[+] ----- read TLV record")
	// fmt.Printf("[#] %s (len: %d) - %+v\n", WysTags[record.Tag], record.Length, record)
	return &record
}

func ParseWYS(compressedWysFile string, args Args) (wys ConfigWYS, err error) {
	zipr, err := zip.OpenReader(compressedWysFile)
	if err != nil {
		return wys, err
	}
	defer zipr.Close()

	for _, f := range zipr.File {
		// there is only one file in the archive
		// "0" is the name of the uncompressed wys file
		if f.FileHeader.Name == "0" {
			fh, err := f.Open()
			if err != nil {
				return wys, err
			}
			defer fh.Close()

			// read HEADER
			b := make([]byte, 7)
			fh.Read(b)

			for {
				tlv := ReadWysTLV(fh)
				if tlv == nil {
					break
				}

				switch tlv.Tag {
				case BYTE_WYS_FILE_SHA1:
					wys.FileSha1 = ValueToByteSlice(tlv)
				case BYTE_WYS_RTF:
					wys.RTF = ValueToByteSlice(tlv)
				case DSTRING_WYS_CURRENT_LAST_VERSION:
					wys.CurrentLastVersion = ValueToString(tlv)
				case DSTRING_WYS_LATEST_CHANGES:
					wys.LatestChanges = ValueToString(tlv)
				case DSTRING_WYS_MIN_CLIENT_VERSION:
					wys.MinClientVersion = ValueToString(tlv)
				case DSTRING_WYS_SERVER_FILE_SITE:
					wys.ServerFileSite = ValueToString(tlv)
				case DSTRING_WYS_UPDATE_ERROR_LINK:
					wys.UpdateErrorLink = ValueToString(tlv)
				case DSTRING_WYS_UPDATE_ERROR_TEXT:
					wys.UpdateErrorText = ValueToString(tlv)
				case DSTRING_WYS_UPDATE_FILE_SITE:
					site := strings.Replace(ValueToString(tlv), "%urlargs%", args.Urlargs, 1)
					wys.UpdateFileSite = append(wys.UpdateFileSite, site)
				case DSTRING_WYS_VERSION_TO_UPDATE:
					wys.VersionToUpdate = ValueToString(tlv)
				case INT_WYS_DUMMY_VAR_LEN:
					// do nothing
				case INT_WYS_FOLDER:
					wys.WysFolder = ValueToInt(tlv)
				case LONG_WYS_UPDATE_FILE_ADLER32_CHECKSUM:
					wys.UpdateFileAdler32 = ValueToLong(tlv)
				case LONG_WYS_UPDATE_FILE_SIZE:
					wys.UpdateFileSize = ValueToLong(tlv)
				default:
					err := fmt.Errorf("wys tag %x not implemented", tlv.Tag)
					return wys, err
				}
			}

			return wys, nil
		}
	}

	// wys not parsed
	err = fmt.Errorf("wys not parsed")
	return wys, err

}
