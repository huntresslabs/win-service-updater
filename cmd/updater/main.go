package main

import (
	"fmt"

	"log"
	"os"

	flags "github.com/jessevdk/go-flags"
)

// Args we support
// "/quickcheck"
// "/justcheck"
// "/noerr",
// "-urlargs="
// "/outputinfo="
// "/fromservice"
// "-logfile="

// https://wyday.com/wybuild/help/wyupdate-commandline.php

type Options struct {
	// Example of verbosity with level
	Quickcheck []bool `long:"quickcheck" description:"foo"`

	// Example of verbosity with level
	Justcheck []bool `long:"justcheck" description:"foo"`

	// Example of verbosity with level
	Noerr []bool `long:"noerr" description:"foo"`

	// Example of verbosity with level
	Fromservice []bool `long:"fromservice" description:"foo"`

	// Example of map with multiple default values
	Urlflags string `long:"urlflags" description:"foo"`
}

var options Options

var parser = flags.NewParser(&options, flags.Default)

func main() {
	fmt.Println("FLAGS")
	_, err := parser.Parse()
	if err != nil {
		flagsErr, ok := err.(*flags.Error)
		if ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			log.Fatal(flagsErr)
			os.Exit(1)
		}
	}
	fmt.Printf("%+v", options)
}
