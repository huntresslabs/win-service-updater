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

// WYS tags
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

// WYSTags is a mapping of WYS tags to strings
var WYSTags = map[uint8]string{
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

// ConfigWYS contains the server file (WYS) details
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
	WYSFolder          int
	UpdateFileAdler32  int64
	UpdateFileSize     int64
}

func ReadWYSTLV(r io.Reader) *TLV {
	var record TLV

	err := binary.Read(r, binary.BigEndian, &record.Tag)
	if err == io.EOF {
		return nil
	} else if err != nil {
		return nil
	}

	if record.Tag == END_WYS {
		return nil
	}

	record.TagString = WYSTags[record.Tag]

	// handle d. strings with the data length
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
			return nil
		}
	default:
	}

	err = binary.Read(r, binary.LittleEndian, &record.Length)
	if err != nil {
		return nil
	}

	// there is no value for the dummy var
	if record.Tag == INT_WYS_DUMMY_VAR_LEN {
		return &record
	}

	record.Value = make([]byte, record.Length)
	_, err = io.ReadFull(r, record.Value)
	if err != nil {
		return nil
	}

	return &record
}

// ParseWYS parses a compress WYS file
func (wysInfo Info) ParseWYS(compressedWYSFile string, args Args) (wys ConfigWYS, err error) {
	zipr, err := zip.OpenReader(compressedWYSFile)
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
			header := make([]byte, 7)
			fh.Read(header)

			if string(header) != WYS_HEADER {
				err = fmt.Errorf("invalid wys header")
				return wys, err
			}

			for {
				tlv := ReadWYSTLV(fh)
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
					wys.UpdateFileSite = append(wys.UpdateFileSite, ValueToString(tlv))
				case DSTRING_WYS_VERSION_TO_UPDATE:
					wys.VersionToUpdate = ValueToString(tlv)
				case INT_WYS_DUMMY_VAR_LEN:
					// do nothing
				case INT_WYS_FOLDER:
					wys.WYSFolder = ValueToInt(tlv)
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

// GetWYUURLs returns the UpdateFileSite(s) included in the WYS file.
func (wys ConfigWYS) GetWYUURLs(args Args) (urls []string) {
	// This can only be specified in tests
	if len(args.WYUTestServer) > 0 {
		urls = append(urls, args.WYUTestServer)
		return urls
	}

	// replace %urlargs% with the args specified on the command line
	for _, s := range wys.UpdateFileSite {
		u := strings.Replace(s, "%urlargs%", args.Urlargs, 1)
		urls = append(urls, u)
	}
	return urls
}
