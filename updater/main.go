package updater

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

const (
	EXIT_NO_UPDATE        = 0
	EXIT_ERROR            = 1
	EXIT_UPDATE_AVALIABLE = 2
)

type UpdateInfoInterface interface {
	ParseWYC(string) (ConfigIUC, error)
}

func Handler() {
	args := ParseArgs(os.Args)
	os.Exit(UpdateHandler(args))
}

func UpdateHandler(args Args) int {
	rc, err := IsUpdateAvailable(args)
	if rc == EXIT_ERROR && nil != err {
		LogErrorMsg(args, err.Error())
		LogOutputInfoMsg(args, err.Error())
	}
	return rc
}

func LogErrorMsg(args Args, msg string) {
	if len(args.Logfile) > 0 {
		dat := []byte(msg)
		ioutil.WriteFile(args.Logfile, dat, 0644)
	}
}

func LogOutputInfoMsg(args Args, msg string) {
	if args.Outputinfo {
		if len(args.OutputinfoLog) > 0 {
			dat := []byte(msg)
			ioutil.WriteFile(args.OutputinfoLog, dat, 0644)
		} else {
			fmt.Println(msg)
		}
	}
}

func IsUpdateAvailable(args Args) (int, error) {
	// read WYC
	iuc, err := ParseWYC(args.Cdata)
	if nil != err {
		return EXIT_ERROR, err
	}

	tmpDir, instDir := Setup()
	defer os.RemoveAll(tmpDir)
	defer os.RemoveAll(instDir)

	wysTmpFile := path.Join(tmpDir, "wysTemp")
	urls := GetWYSURLs(iuc, args)

	// TODO loop through URLs here or in DownloadFile()
	err = DownloadFile(urls[0], wysTmpFile)
	if nil != err {
		return EXIT_ERROR, err
	}

	wys, err := ParseWYS(wysTmpFile, args)
	if nil != err {
		return EXIT_ERROR, err
	}

	// compare versions
	rc := CompareVersions(string(iuc.IucInstalledVersion.Value), wys.VersionToUpdate)
	switch rc {
	case A_LESS_THAN_B:
		// need update
		return EXIT_UPDATE_AVALIABLE, nil
	case A_EQUAL_TO_B:
		// no update
		return EXIT_NO_UPDATE, nil
	case A_GREATER_THAN_B:
		// no update
		return EXIT_NO_UPDATE, nil
	default:
		// unknown case
		return EXIT_ERROR, fmt.Errorf("unknown case")
	}
}
