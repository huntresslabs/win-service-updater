package updater

import (
	"hash/adler32"
	"io/ioutil"

	"github.com/davecgh/go-spew/spew"
)

func GetAdler32(file string) (uint32, error) {
	dat, err := ioutil.ReadFile(file)
	if nil != err {
		return 0, err
	}
	return adler32.Checksum(dat), nil
}

func VerifyAdler32Checksum(expected int64, file string) bool {
	cs, err := GetAdler32(file)
	if nil != err {
		return false
	}

	spew.Dump(expected)
	spew.Dump(cs)

	if cs == uint32(expected) {
		return true
	}

	return false
}
