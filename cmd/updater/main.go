package main

import (
	"fmt"
	"os"

	"github.com/huntresslabs/win-service-updater/updater"
)

// rc, _ := try(WYUPDATE_EXE, "/quickcheck", "/justcheck", "/noerr",
// 	fmt.Sprintf("-urlargs=%s", AUTH), fmt.Sprintf("/outputinfo=%s", CHECK_LOG))

// wyupdateArgs := fmt.Sprintf("/fromservice -logfile=\"%s\" -urlargs=%s",
// 	WYUPDATE_LOG, AUTH)

func main() {
	args := updater.ParseArgs(os.Args)
	fmt.Printf("%+v", args)
}
