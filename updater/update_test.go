package updater

import (
	"fmt"
	"testing"
)

func TestUpdate(t *testing.T) {
	err := fmt.Errorf("not tested")
	t.Fatal(err)
	// tempExtract, err := ioutil.TempDir("", "prefix")
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// defer os.RemoveAll(tempExtract)

	// tempInstall, err := ioutil.TempDir("", "prefix")
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// defer os.RemoveAll(tempInstall)

	// src := "some.wyu"
	// _, files, err := Unzip(src, tempExtract)
	// assert.Nil(t, err)

	// udt, updateFiles, err := GetUpdateDetails(files)
	// assert.Nil(t, err)

	// err = InstallUpdate(udt, updateFiles, tempInstall)
	// assert.Nil(t, err)
}
