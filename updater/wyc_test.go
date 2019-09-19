package updater

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func appendFiles(filename string, zipw *zip.Writer) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("Failed to open %s: %s", filename, err)
	}
	defer file.Close()

	wr, err := zipw.Create(path.Base(filename))
	if err != nil {
		msg := "Failed to create entry for %s in zip file: %s"
		return fmt.Errorf(msg, filename, err)
	}

	if _, err := io.Copy(wr, file); err != nil {
		return fmt.Errorf("Failed to write %s to zip: %s", filename, err)
	}

	return nil
}

func Zip(archive string, files []string) {
	flags := os.O_WRONLY | os.O_CREATE | os.O_TRUNC
	file, err := os.OpenFile(archive, flags, 0644)
	if err != nil {
		log.Fatalf("Failed to open zip for writing: %s", err)
	}
	defer file.Close()

	zipw := zip.NewWriter(file)
	defer zipw.Close()

	for _, filename := range files {
		if err := appendFiles(filename, zipw); err != nil {
			log.Fatalf("Failed to add file %s to zip: %s", filename, err)
		}
	}
}

func TestWYC(t *testing.T) {
	origFile := "../test_files/client.1.0.0.wyc"
	wyc, err := ParseWYC(origFile)
	assert.Nil(t, err)
	assert.Equal(t, wyc.IucServerFileSite[0].Value, []byte("http://127.0.0.1/updates/wyserver.wys?key=%urlargs%"))
}

func TestWYC_UpdateWYC(t *testing.T) {
	origFile := "../test_files/client.1.0.0.wyc"

	wyc, err := ParseWYC(origFile)
	assert.Nil(t, err)

	// TODO create a function to do this, we'll need to do this after an update
	wyc.IucInstalledVersion.Value = []byte("1.2.3")
	wyc.IucInstalledVersion.DataLength = uint32(len(wyc.IucInstalledVersion.Value) + 4)
	wyc.IucInstalledVersion.Length = uint32(len(wyc.IucInstalledVersion.Value))

	new, err := UpdateWYC(wyc, origFile)
	assert.Nil(t, err)
	fmt.Println(new)

	_, err = ParseWYC(new)
	assert.Nil(t, err)
}

func TestWYC_WriteIUC(t *testing.T) {
	// create a new uiclient.iuc and compare it to the one in the archive

	origClientWYC := "../test_files/client.1.0.0.wyc"

	wyc, err := ParseWYC(origClientWYC)
	assert.Nil(t, err)

	tmpIUC, err := ioutil.TempFile("", "example")
	assert.Nil(t, err)
	tmpIUC.Close()
	// fmt.Println("tmpIUC", tmpIUC.Name())
	defer os.Remove(tmpIUC.Name())

	err = WriteIUC(wyc, tmpIUC.Name())
	assert.Nil(t, err)

	tmpDir, err := ioutil.TempDir("", "prefix")
	assert.Nil(t, err)
	defer os.RemoveAll(tmpDir)

	newHash, err := Sha256Hash(tmpIUC.Name())
	assert.Nil(t, err)

	found := false
	_, files, err := Unzip(origClientWYC, tmpDir)
	for _, f := range files {
		// fmt.Println(f)
		if filepath.Base(f) == IUCLIENT_IUC {
			origHash, err := Sha256Hash(f)
			assert.Nil(t, err)
			assert.Equal(t, origHash, newHash)
			found = true
		}
	}
	assert.True(t, found)
}
