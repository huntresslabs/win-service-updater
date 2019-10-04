package updater

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUpdate(t *testing.T) {
	tempExtract := GetTempDir()
	defer os.RemoveAll(tempExtract)

	tempInstall := GetTempDir()
	defer os.RemoveAll(tempInstall)

	src := "../test_files/widgetX.1.0.1.wyu"
	_, files, err := Unzip(src, tempExtract)
	assert.Nil(t, err)

	udt, updateFiles, err := GetUpdateDetails(files)
	assert.Nil(t, err)

	err = InstallUpdate(udt, updateFiles, tempInstall)
	assert.Nil(t, err)
}
