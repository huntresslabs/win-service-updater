package updater

import (
	"fmt"
	"os"
)

const (
	EXIT_NO_UPDATE        = 0
	EXIT_ERROR            = 1
	EXIT_UPDATE_AVALIABLE = 2
)

type UpdateInfoInterface interface {
	ParseWYC(string) (ConfigIUC, error)
}

func UpdateHandler() {
	args := ParseArgs(os.Args)

	uier := UpdateInfoer{}

	exitCode := IsUpdateAvailable(uier, "./test_files2/client1.0.0.wyc", args)

	os.Exit(exitCode)
	// if exitCode == EXIT_NO_UPDATE {
	// 	os.Exit(0)
	// } else if exitCode == EXIT_ERROR {
	// 	os.Exit(1)
	// } else if exitCode == EXIT_UPDATE_AVALIABLE {
	// 	os.Exit(2)
	// }

}

func IsUpdateAvailable(uii UpdateInfoInterface, wycFile string, args Args) int {
	// read WYC
	iuc, err := uii.ParseWYC(wycFile)
	if args.Noerr && nil != err {
		if len(args.Logfile) > 0 {
			// log error
		}
		return 1
	}

	// This is test code to test if I can update the URL value from an interface
	if string(iuc.IucServerFileSite[0].Value) == "TEST_URL" {
		return 3
	} else {
		return 1
	} // end test code

	tmpDir, instDir := Setup()
	defer os.RemoveAll(tmpDir)
	defer os.RemoveAll(instDir)

	wysTmpFile := fmt.Sprintf("%s\\wys", tmpDir)
	// download WYS and extract
	err = DownloadFile(string(iuc.IucServerFileSite[0].Value), wysTmpFile)

	wys, err := ParseWYS(wysTmpFile, args)
	if args.Noerr && nil != err {
		if len(args.Logfile) > 0 {
			// log error
		}
		return 1
	}
	fmt.Printf("%+v\n", wys)

	// compare versions
	rc := CompareVersions(string(iuc.IucInstalledVersion.Value), wys.VersionToUpdate)
	// need update
	if rc == A_LESS_THAN_B {
		// Log new version to args.Outputinfo
		return 2
	}
	// no update
	return 0
}

func Update() {

}
