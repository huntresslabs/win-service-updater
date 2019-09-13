package main

import (
	"log"
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/huntresslabs/win-service-updater/updater"
)

func main() {
	args := updater.Args{}
	wys, err := updater.ParseWYS(os.Args[1], args)
	if nil != err {
		log.Fatal(err)
	}
	spew.Dump("%+v", wys)
}
