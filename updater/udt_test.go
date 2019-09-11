package updater

import (
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUDT(t *testing.T) {
	tmpfile, err := ioutil.TempFile("", "example")
	if err != nil {
		log.Fatal(err)
	}
	tmpfile.Close()
	defer os.Remove(tmpfile.Name())

	orig := "../test_files/updtdetails.udt"
	udt, err := ParseUDT(orig)
	assert.Nil(t, err)

	err = WriteUDT(udt, tmpfile.Name())
	assert.Nil(t, err)

	origHash, err := Sha256Hash(orig)
	assert.Nil(t, err)

	newHash, err := Sha256Hash(tmpfile.Name())
	assert.Nil(t, err)
	assert.Equal(t, origHash, newHash)
}
