import struct
import sys
import logging
import os

from wybinary import WyBinary

## 1 byte: Identifier (e.g. 0xNN)
## 4 bytes: Little-endian, 32-bit signed integer length of the data (in bytes: e.g. X bytes)
## X bytes: Data
## Header: IUSDFV2

## The parser reads an uncompressed wys file and returns a dictionary
## The write takes a dictionary

DEBUG = False

WYS_CURRENT_LAST_VERSION_DSTRING = 0x01
WYS_SERVER_FILE_SITE_DSTRING = 0x02
WYS_MIN_CLIENT_VERSION_DSTRING = 0x07
WYS_DUMMY_VAR_LEN_INT = 0x0F
WYS_VERSION_TO_UPDATE_DSTRING = 0x0B
WYS_UPDATE_FILE_SITE_DSTRING = 0x03
WYS_RTF_BYTE = 0x80
WYS_LATEST_CHANGES_DSTRING = 0x04
WYS_UPDATE_FILE_SIZE_LONG = 0x09
WYS_UPDATE_FILE_ADLER32_CHECKSUM_LONG = 0x08
WYS_FILE_SHA1_BYTE = 0x14
WYS_FOLDER_INT = 0x0A
WYS_UPDATE_ERROR_TEXT_DSTRING = 0x20
WYS_UPDATE_ERROR_LINK_DSTRING = 0x21
WYS_END = 0xFF

int_to_string_mapping = {
    WYS_CURRENT_LAST_VERSION_DSTRING:      'WYS_CURRENT_LAST_VERSION_DSTRING',
    WYS_SERVER_FILE_SITE_DSTRING:          'WYS_SERVER_FILE_SITE_DSTRING',
    WYS_UPDATE_FILE_SITE_DSTRING:          'WYS_UPDATE_FILE_SITE_DSTRING',
    WYS_FOLDER_INT:                        'WYS_FOLDER_INT',
    WYS_LATEST_CHANGES_DSTRING:            'WYS_LATEST_CHANGES_DSTRING',
    WYS_UPDATE_FILE_SIZE_LONG:             'WYS_UPDATE_FILE_SIZE_LONG',
    WYS_VERSION_TO_UPDATE_DSTRING:         'WYS_VERSION_TO_UPDATE_DSTRING',
    WYS_UPDATE_FILE_ADLER32_CHECKSUM_LONG: 'WYS_UPDATE_FILE_ADLER32_CHECKSUM_LONG',
    WYS_FILE_SHA1_BYTE:                    'WYS_FILE_SHA1_BYTE',
    WYS_UPDATE_ERROR_TEXT_DSTRING:         'WYS_UPDATE_ERROR_TEXT_DSTRING',
    WYS_UPDATE_ERROR_LINK_DSTRING:         'WYS_UPDATE_ERROR_LINK_DSTRING',
    WYS_MIN_CLIENT_VERSION_DSTRING:        'WYS_MIN_CLIENT_VERSION_DSTRING',
    WYS_DUMMY_VAR_LEN_INT:                 'WYS_DUMMY_VAR_LEN_INT'}

string_to_int_mapping = {
    'WYS_CURRENT_LAST_VERSION_DSTRING':      WYS_CURRENT_LAST_VERSION_DSTRING,
    'WYS_SERVER_FILE_SITE_DSTRING':          WYS_SERVER_FILE_SITE_DSTRING,
    'WYS_UPDATE_FILE_SITE_DSTRING':          WYS_UPDATE_FILE_SITE_DSTRING,
    'WYS_FOLDER_INT':                        WYS_FOLDER_INT,
    'WYS_LATEST_CHANGES_DSTRING':            WYS_LATEST_CHANGES_DSTRING,
    'WYS_UPDATE_FILE_SIZE_LONG':             WYS_UPDATE_FILE_SIZE_LONG,
    'WYS_VERSION_TO_UPDATE_DSTRING':         WYS_VERSION_TO_UPDATE_DSTRING,
    'WYS_UPDATE_FILE_ADLER32_CHECKSUM_LONG': WYS_UPDATE_FILE_ADLER32_CHECKSUM_LONG,
    'WYS_FILE_SHA1_BYTE':                    WYS_FILE_SHA1_BYTE,
    'WYS_UPDATE_ERROR_TEXT_DSTRING':         WYS_UPDATE_ERROR_TEXT_DSTRING,
    'WYS_UPDATE_ERROR_LINK_DSTRING':         WYS_UPDATE_ERROR_LINK_DSTRING,
    'WYS_MIN_CLIENT_VERSION_DSTRING':        WYS_MIN_CLIENT_VERSION_DSTRING,
    'WYS_DUMMY_VAR_LEN_INT':                 WYS_DUMMY_VAR_LEN_INT}

FOLDER_BASE_DIR = 1
FOLDER_SYSTEM32_DIR = 2

# create logger
logger = logging.getLogger(__name__)
if DEBUG:
    logger.setLevel(logging.DEBUG)
else:
    logger.setLevel(logging.INFO)
# create file handler which logs even debug messages
#fh = logging.FileHandler('parser.log')
#fh.setLevel(logging.DEBUG)
# create console handler with a higher log level
ch = logging.StreamHandler()
ch.setLevel(logging.DEBUG)
# create formatter and add it to the handlers
#formatter = logging.Formatter('%(asctime)s - %(name)s - %(levelname)s - %(message)s')
formatter = logging.Formatter('%(asctime)s - %(levelname)s - %(lineno)s - %(message)s')
#fh.setFormatter(formatter)
ch.setFormatter(formatter)
# add the handlers to the logger
#logger.addHandler(fh)
logger.addHandler(ch)


class WysBinary(WyBinary):
    """Parse/create wys files"""
    def __init__(self):
        self.int_to_string_mapping = int_to_string_mapping
        self.string_to_int_mapping = string_to_int_mapping

    def __del__(self):
        self.file.close()

    def read_dummy(self, t):
        ## 32 int
        i = struct.unpack('<I', self.file.read(4))[0]
        self.data[self.int_to_string_mapping[t]] = i
        self.log(t, i)

    def read_header(self):
        self.data['HEADER'] = self.file.read(7)

    def parse(self, filename):
        self.file = open(filename, 'rb')
        self.data = {}

        self.read_header()
        while True:
            t = struct.unpack('B', self.file.read(1))[0]
            #print(repr(t))
            if t == WYS_CURRENT_LAST_VERSION_DSTRING:
                self.read_dstring(t)
            elif t == WYS_MIN_CLIENT_VERSION_DSTRING:
                self.read_dstring(t)
            elif t == WYS_SERVER_FILE_SITE_DSTRING:
                self.read_dstring(t)
            elif t == WYS_FOLDER_INT:
                self.read_int(t)
            elif t == WYS_UPDATE_FILE_SITE_DSTRING:
                self.read_dstring(t)
            elif t == WYS_LATEST_CHANGES_DSTRING:
                self.read_dstring(t)
            elif t == WYS_UPDATE_FILE_SIZE_LONG:
                self.read_long(t)
            elif t == WYS_VERSION_TO_UPDATE_DSTRING:
                self.read_dstring(t)
            elif t == WYS_UPDATE_FILE_ADLER32_CHECKSUM_LONG:
                self.read_long(t)
            elif t == WYS_FILE_SHA1_BYTE:
                self.read_string(t)
            elif t == WYS_UPDATE_ERROR_TEXT_DSTRING:
                self.read_dstring(t)
            elif t == WYS_UPDATE_ERROR_LINK_DSTRING:
                self.read_dstring(t)
            elif t == WYS_DUMMY_VAR_LEN_INT:
                self.read_dummy(t)
            elif t == WYS_END:
                return self.data
            else:
                logger.error('Unknown type: {0}'.format(t))

    def write_header(self):
        self.file.write(self.data['HEADER'])

    def write_dummy(self, type_, value):
        ## 32 int
        i = struct.pack('<I', value)
        data = struct.pack('<B', self.string_to_int_mapping[type_]) + i
        self.file.write(data)
        logger.debug('{0}: {1}'.format(type_, repr(data)))

    def write_end(self):
        self.file.write(struct.pack('<I', WYS_END))
        self.file.close()

    def create(self, filename, iuc):
        self.file = open(filename, 'wb')
        self.data = iuc

        self.write_header()
        for k, v in self.data.iteritems():
            if k == 'WYS_CURRENT_LAST_VERSION_DSTRING':
                self.write_dstring(k, v)
            elif k == 'WYS_SERVER_FILE_SITE_DSTRING':
                self.write_dstring(k, v)
            elif k == 'WYS_FOLDER_INT':
                self.write_int(k, v)
            elif k == 'WYS_UPDATE_FILE_SITE_DSTRING':
                self.write_dstring(k, v)
            elif k == 'WYS_LATEST_CHANGES_DSTRING':
                self.write_dstring(k, v)
            elif k == 'WYS_UPDATE_FILE_SIZE_LONG':
                self.write_long(k, v)
            elif k == 'WYS_VERSION_TO_UPDATE_DSTRING':
                self.write_dstring(k, v)
            elif k == 'WYS_UPDATE_FILE_ADLER32_CHECKSUM_LONG':
                self.write_long(k, v)
            elif k == 'WYS_FILE_SHA1_BYTE':
                self.write_string(k, v)
            elif k == 'WYS_UPDATE_ERROR_TEXT_DSTRING':
                self.write_dstring(k, v)
            elif k == 'WYS_UPDATE_ERROR_LINK_DSTRING':
                self.write_dstring(k, v)
            elif k == 'WYS_MIN_CLIENT_VERSION_DSTRING':
                self.write_dstring(k, v)
            elif k == 'WYS_DUMMY_VAR_LEN_INT':
                self.write_dummy(k, v)
            elif k == 'HEADER':
                ## header already written
                pass
            else:
                logger.error('Unknown type: {0}'.format(k))
        self.write_end()

if __name__ == '__main__':
    # test by parsing file, generating a new file using the parsed values,
    # parse the new file, and then compare original and new
    # TEST_WYS_FILE = "wys.file"
    # p1 = WysBinary()
    # d1 = p1.parse(TEST_WYS_FILE)

    # w = WysBinary()
    # w.create('tmp', d1)

    # p2 = WysBinary()
    # d2 = p2.parse('tmp')
    # for k in d2.keys():
    #     print(k, repr(d2[k]))
    # assert(d1 == d2)

    # create an uncompress wys file with the following config
    WYS_CONFIG = {
        'HEADER': 'IUSDFV2',
        'WYS_MIN_CLIENT_VERSION_DSTRING': '2.6.18.4',
        'WYS_FOLDER_INT': 513,
        'WYS_CURRENT_LAST_VERSION_DSTRING': '1.0.0.1',
        'WYS_LATEST_CHANGES_DSTRING': '',
        'WYS_DUMMY_VAR_LEN_INT': 145,
        'WYS_VERSION_TO_UPDATE_DSTRING': '1.0.0.2',
        'WYS_UPDATE_FILE_SITE_DSTRING': 'https://127.0.0.1/updates/update.wyu?auth=%urlargs%',
        #'WYS_UPDATE_FILE_SITE_DSTRING': 'http://192.168.1.20:8080/%s' % UPDATE_ARCHIVE,
        'WYS_UPDATE_FILE_SIZE_LONG': 6503371,
        'WYS_UPDATE_FILE_ADLER32_CHECKSUM_LONG': 3388085357,
    }

    w = WysBinary()
    w.create('wys.test', WYS_CONFIG)
