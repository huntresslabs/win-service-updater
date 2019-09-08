package main

import (
	"fmt"
	"os"

	"github.com/huntresslabs/win-service-updater/updater"
)

func main() {
	args := updater.ParseArgs(os.Args)
	fmt.Printf("%+v", args)
}
