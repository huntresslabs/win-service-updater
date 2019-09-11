package updater

import (
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFile_Unzip(t *testing.T) {
	out, err := ioutil.TempDir("", "prefix")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(out)

	zipFile := "../test_files/test.zip"
	_, files, err := Unzip(zipFile, out)
	assert.Nil(t, err)
	// there are 4 files in the test.zip
	assert.Equal(t, 4, len(files))
}
