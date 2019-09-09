import struct
import sys
import logging
import os

## 1 byte: Identifier (e.g. 0xNN)
## 4 bytes: Little-endian, 32-bit signed integer length of the data (in bytes: e.g. X bytes)
## X bytes: Data

# create logger
logger = logging.getLogger(__name__)

class WyBinary(object):

    def __init__(self):
        pass

    def log(self, type_, value):
        if type(type_) == int:
            logger.info('{0}: {1}'.format(self.int_to_string_mapping[type_], value))
        else:
            logger.info('{0}: {1}'.format(self.string_to_int_mapping[type_], value))

    def parse(self):
        raise "Not Implemented"

    def create(self):
        raise "Not Implemented"

    def read_header(self):
        raise "Not Implemented"

    def write_header(self):
        raise "Not Implemented"

    def write_end(self):
        raise "Not Implemented"

    def read_dstring(self, t):
        logger.debug('reading dstring')
        data_len = struct.unpack('<I', self.file.read(4))[0]
        logger.debug('data len: {0}'.format(data_len))
        string_len = struct.unpack('<I', self.file.read(4))[0]
        dstring = self.file.read(string_len)
        self.data[self.int_to_string_mapping[t]] = dstring
        self.log(t, dstring)

    def read_string(self, t):
        logger.debug('reading string')
        string_len = struct.unpack('<I', self.file.read(4))[0]
        logger.debug('string len: {0}'.format(string_len))
        s = self.file.read(string_len)
        self.data[self.int_to_string_mapping[t]] = s
        self.log(t, s)

    def read_dummy(self, t):
        ## 32 int
        i = struct.unpack('<I', self.file.read(4))[0]
        self.data[self.int_to_string_mapping[t]] = i
        self.log(t, i)

    def read_int(self, t):
        ## could be 32 or 64 int, use length to determine
        length = struct.unpack('<I', self.file.read(4))[0]
        i = struct.unpack('<I', self.file.read(length))[0]
        self.data[self.int_to_string_mapping[t]] = i
        self.log(t, i)

    def read_long(self, t):
        ## could be 32 or 64 int, use length to determine
        length = struct.unpack('<I', self.file.read(4))[0]
        i = struct.unpack('<Q', self.file.read(length))[0]
        self.data[self.int_to_string_mapping[t]] = i
        self.log(t, i)

    def write_dstring(self, type_, s):
        logger.debug('writing dstring')
        string_len = struct.pack('<I', len(s))
        data_len = struct.pack('<I', len(s) + 4)
        data = struct.pack('<B', self.string_to_int_mapping[type_]) + data_len + string_len + s
        self.file.write(data)
        logger.debug('{0}: {1}'.format(type_, repr(data)))

    def write_string(self, type_, s):
        logger.debug('writing string')
        string_len = struct.pack('<I', len(s))
        data = struct.pack('<B', self.string_to_int_mapping[type_]) + string_len + s
        self.file.write(data)
        logger.debug('{0}: {1}'.format(type_, repr(data)))

    def write_int(self, type_, value):
        ## could be 32 or 64 int, use length to determine
        value_bin = struct.pack('<I', value)
        value_bin_len = struct.pack('<I', len(value_bin))
        data = struct.pack('<B', self.string_to_int_mapping[type_]) + value_bin_len + value_bin
        self.file.write(data)
        logger.debug('{0}: {1}'.format(type_, repr(data)))

    def write_long(self, type_, value):
        length_bin = struct.pack('<I', 8)
        value_bin = struct.pack('<Q', value)
        data = struct.pack('<B', self.string_to_int_mapping[type_]) + length_bin + value_bin
        self.file.write(data)
        logger.debug('{0}: {1}'.format(type_, repr(data)))