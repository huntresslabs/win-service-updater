package updater

import (
	"hash/adler32"
	"io/ioutil"
)

// GetAdler32 returns the Adler32 checksum of a file
func GetAdler32(file string) (uint32, error) {
	dat, err := ioutil.ReadFile(file)
	if nil != err {
		return 0, err
	}
	return adler32.Checksum(dat), nil
}

// VerifyAdler32Checksum returns true if checksum was verified
func VerifyAdler32Checksum(expected int64, file string) bool {
	cs, err := GetAdler32(file)
	if nil != err {
		return false
	}

	if cs == uint32(expected) {
		return true
	}

	return false
}
