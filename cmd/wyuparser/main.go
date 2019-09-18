package main

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/huntresslabs/win-service-updater/updater"
)

func main() {
	tmpDir, err := ioutil.TempDir("", "prefix")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// extract wyu to tmpDir
	_, files, err := updater.Unzip(os.Args[1], tmpDir)
	if err != nil {
		log.Fatal(err)
	}

	udt, _, err := updater.GetUpdateDetails(files)
	if err != nil {
		log.Fatal(err)
	}
	spew.Dump("%+v", udt)
}
