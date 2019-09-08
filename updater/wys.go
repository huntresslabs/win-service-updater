package updater

const (
	WYS_CURRENT_LAST_VERSION_DSTRING      = 0x01
	WYS_SERVER_FILE_SITE_DSTRING          = 0x02
	WYS_MIN_CLIENT_VERSION_DSTRING        = 0x07
	WYS_DUMMY_VAR_LEN_INT                 = 0x0F
	WYS_VERSION_TO_UPDATE_DSTRING         = 0x0B
	WYS_UPDATE_FILE_SITE_DSTRING          = 0x03
	WYS_RTF_BYTE                          = 0x80
	WYS_LATEST_CHANGES_DSTRING            = 0x04
	WYS_UPDATE_FILE_SIZE_LONG             = 0x09
	WYS_UPDATE_FILE_ADLER32_CHECKSUM_LONG = 0x08
	WYS_FILE_SHA1_BYTE                    = 0x14
	WYS_FOLDER_INT                        = 0x0A
	WYS_UPDATE_ERROR_TEXT_DSTRING         = 0x20
	WYS_UPDATE_ERROR_LINK_DSTRING         = 0x21
	WYS_END                               = 0xFF
)

var WysTags = map[uint8]string{
	WYS_CURRENT_LAST_VERSION_DSTRING      = 0x01
	WYS_SERVER_FILE_SITE_DSTRING          = 0x02
	WYS_MIN_CLIENT_VERSION_DSTRING        = 0x07
	WYS_DUMMY_VAR_LEN_INT                 = 0x0F
	WYS_VERSION_TO_UPDATE_DSTRING         = 0x0B
	WYS_UPDATE_FILE_SITE_DSTRING          = 0x03
	WYS_RTF_BYTE                          = 0x80
	WYS_LATEST_CHANGES_DSTRING            = 0x04
	WYS_UPDATE_FILE_SIZE_LONG             = 0x09
	WYS_UPDATE_FILE_ADLER32_CHECKSUM_LONG = 0x08
	WYS_FILE_SHA1_BYTE                    = 0x14
	WYS_FOLDER_INT                        = 0x0A
	WYS_UPDATE_ERROR_TEXT_DSTRING         = 0x20
	WYS_UPDATE_ERROR_LINK_DSTRING         = 0x21
	WYS_END                               = 0xFF