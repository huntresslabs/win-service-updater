package updater

import (
	"crypto/rsa"
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

// rc, _ := try(WYUPDATE_EXE, "/quickcheck", "/justcheck", "/noerr",
// 	fmt.Sprintf("-urlargs=%s", AUTH), fmt.Sprintf("/outputinfo=%s", CHECK_LOG))

// wyupdateArgs := fmt.Sprintf("/fromservice -logfile=\"%s\" -urlargs=%s",
// 	WYUPDATE_LOG, AUTH)

func Handler() int {
	args := ParseArgs(os.Args)
	if args.Quickcheck && args.Justcheck {
		rc, err := IsUpdateAvailable(args)
		if nil != err {
			err = fmt.Errorf("exit code=%d; %s", rc, err)
			LogErrorMsg(args, err.Error())
			LogOutputInfoMsg(args, err.Error())
		}
		return rc
	}

	if args.Fromservice {
		rc, err := UpdateHandler((args))
		if nil != err {
			LogErrorMsg(args, err.Error())
			LogOutputInfoMsg(args, err.Error())
		}
		return rc
	}

	return EXIT_ERROR
}

func UpdateHandler(args Args) (int, error) {
	tmpDir := findTempDir()
	instDir := GetExeDir()
	defer os.RemoveAll(tmpDir)

	iuc, err := ParseWYC(args.Cdata)
	if nil != err {
		err = fmt.Errorf("error reading %s; %w", args.Cdata, err)
		return EXIT_ERROR, err
	}

	fp := fmt.Sprintf("%s/wys", tmpDir)
	urls := GetWYSURLs(iuc, args)
	err = DownloadFile(urls[0], fp)
	if nil != err {
		err = fmt.Errorf("download error; %w", err)
		return EXIT_ERROR, err
	}

	wys, err := ParseWYS(fp, args)
	if nil != err {
		err = fmt.Errorf("error reading wys file (%s); %w", fp, err)
		return EXIT_ERROR, err
	}

	// fmt.Println("installed ", string(iuc.IucInstalledVersion.Value))
	// fmt.Println("new ", wys.VersionToUpdate)

	// download wyu
	fp = fmt.Sprintf("%s/wyu", tmpDir)
	urls = GetWYUURLs(wys, args)
	err = DownloadFile(urls[0], fp)
	if nil != err {
		err = fmt.Errorf("error download update archive; %w", err)
		return EXIT_ERROR, err
	}

	key, err := ParsePublicKey(string(iuc.IucPublicKey.Value))
	var rsa rsa.PublicKey
	rsa.N = key.Modulus
	rsa.E = key.Exponent

	sha1hash, err := Sha1Hash(fp)
	if nil != err {
		err = fmt.Errorf("error hashing %s; %w", fp, err)
		return EXIT_ERROR, err
	}

	// validated
	err = VerifyHash(&rsa, sha1hash, wys.FileSha1)
	if nil != err {
		err = fmt.Errorf("error verifying %s; %w", fp, err)
		return EXIT_ERROR, err
	}

	// adler32
	if wys.UpdateFileAdler32 != 0 {
		v := VerifyAdler32Checksum(wys.UpdateFileAdler32, fp)
		if v != true {
			return EXIT_ERROR, err
		}
	}

	// extract wyu to tmpDir
	_, files, err := Unzip(fp, tmpDir)
	if nil != err {
		err = fmt.Errorf("error unzipping %s; %w", fp, err)
		return EXIT_ERROR, err
	}

	udt, updates, err := GetUpdateDetails(files)
	if nil != err {
		return EXIT_ERROR, err
	}

	backupDir, err := BackupFiles(updates, instDir)
	if nil != err {
		return EXIT_ERROR, err
	}

	err = InstallUpdate(udt, updates, instDir)
	if nil != err {
		// TODO start services
		RollbackFiles(backupDir, instDir)
		return EXIT_ERROR, err
	} else {
		return 0, nil
	}
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

	tmpDir := findTempDir()
	defer os.RemoveAll(tmpDir)

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
		err = fmt.Errorf(wys.VersionToUpdate)
		return EXIT_UPDATE_AVALIABLE, err
	case A_EQUAL_TO_B:
		// no update
		err = fmt.Errorf(wys.VersionToUpdate)
		return EXIT_NO_UPDATE, err
	case A_GREATER_THAN_B:
		// no update
		err = fmt.Errorf(string(iuc.IucInstalledVersion.Value))
		return EXIT_NO_UPDATE, err
	default:
		// unknown case
		return EXIT_ERROR, fmt.Errorf("unknown case")
	}
}
