package main

import (
	"os"

	"github.com/huntresslabs/win-service-updater/updater"
)

func main() {
	os.Exit(updater.Handler())
}
