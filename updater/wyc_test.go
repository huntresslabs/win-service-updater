package updater

import (
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWYC(t *testing.T) {
	origFile := "../test_files/iuclient.iuc"
	wyc, err := ParseWyc(origFile)
	assert.Nil(t, err)
	assert.Equal(t, wyc.IucServerFileSite[0].Value, []byte("http://127.0.0.1/update.wys"))

	tmpfile, err := ioutil.TempFile("", "example")
	if err != nil {
		log.Fatal(err)
	}
	tmpfile.Close()
	defer os.Remove(tmpfile.Name())

	err = WriteWYC(wyc, tmpfile.Name())
	assert.Nil(t, err)

	origHash, err := Sha256Hash(origFile)
	assert.Nil(t, err)

	newHash, err := Sha256Hash(tmpfile.Name())
	assert.Nil(t, err)
	assert.Equal(t, origHash, newHash)
}
