package updater

import (
	"crypto/rsa"
	"fmt"
	"os"
	"path"
)

const (
	EXIT_SUCCESS          = 0
	EXIT_NO_UPDATE        = 0
	EXIT_ERROR            = 1
	EXIT_UPDATE_AVALIABLE = 2
)

// rc, _ := try(WYUPDATE_EXE, "/quickcheck", "/justcheck", "/noerr",
// 	fmt.Sprintf("-urlargs=%s", AUTH), fmt.Sprintf("/outputinfo=%s", CHECK_LOG))

// wyupdateArgs := fmt.Sprintf("/fromservice -logfile=\"%s\" -urlargs=%s",
// 	WYUPDATE_LOG, AUTH)

func Handler() int {
	args := ParseArgs(os.Args)
	info := Info{}

	// check for updates
	if args.Quickcheck && args.Justcheck {
		rc, err := IsUpdateAvailable(info, args)
		if nil != err {
			LogErrorMsg(args, err.Error())
			LogOutputInfoMsg(args, err.Error())
		}
		return rc
	}

	// update
	if args.Fromservice {
		rc, err := UpdateHandler(info, (args))
		if nil != err {
			LogErrorMsg(args, err.Error())
			LogOutputInfoMsg(args, err.Error())
		}
		return rc
	}

	return EXIT_ERROR
}

func UpdateHandler(infoer Infoer, args Args) (int, error) {
	tmpDir := findTempDir()
	instDir := GetExeDir()
	defer os.RemoveAll(tmpDir)

	// parse the WYC file for get update site, installed version, etc.
	iuc, err := infoer.ParseWYC(args.Cdata)
	if nil != err {
		err = fmt.Errorf("error reading %s; %w", args.Cdata, err)
		return EXIT_ERROR, err
	}

	// download the wys file (contains details about the availiable update)
	fp := fmt.Sprintf("%s/wys", tmpDir)
	urls := GetWYSURLs(iuc, args)
	err = DownloadFile(urls[0], fp)
	if nil != err {
		err = fmt.Errorf("error downloading wys file; %w", err)
		return EXIT_ERROR, err
	}

	// parse the WYS file (contains the version number of the update and the link to the update)
	wys, err := infoer.ParseWYS(fp, args)
	if nil != err {
		err = fmt.Errorf("error reading wys file (%s); %w", fp, err)
		return EXIT_ERROR, err
	}

	// fmt.Println("installed ", string(iuc.IucInstalledVersion.Value))
	// fmt.Println("new ", wys.VersionToUpdate)

	// download WYU (this is the archive with the updated files)
	fp = fmt.Sprintf("%s/wyu", tmpDir)
	urls = GetWYUURLs(wys, args)
	err = DownloadFile(urls[0], fp)
	if nil != err {
		err = fmt.Errorf("error downloading update archive; %w", err)
		return EXIT_ERROR, err
	}

	if iuc.IucPublicKey.Value != nil {
		if len(wys.FileSha1) == 0 {
			err = fmt.Errorf("error validating WYU file; PublicKey was specified but there was no signed hash")
			return EXIT_ERROR, err
		}

		// convert the public key from the WYC file to an rsa.PublicKey
		key, err := ParsePublicKey(string(iuc.IucPublicKey.Value))
		var rsa rsa.PublicKey
		rsa.N = key.Modulus
		rsa.E = key.Exponent

		// hash the downloaded WYU file
		sha1hash, err := Sha1Hash(fp)
		if nil != err {
			err = fmt.Errorf("error hashing %s; %w", fp, err)
			return EXIT_ERROR, err
		}

		// verify the signature of the WYU file (the signed hash is included in the WYS file)
		err = VerifyHash(&rsa, sha1hash, wys.FileSha1)
		if nil != err {
			err = fmt.Errorf("error verifying %s; %w", fp, err)
			return EXIT_ERROR, err
		}
	}

	// adler32 checksum
	if wys.UpdateFileAdler32 != 0 {
		v := VerifyAdler32Checksum(wys.UpdateFileAdler32, fp)
		if v != true {
			return EXIT_ERROR, err
		}
	}

	// extract the WYU to tmpDir
	_, files, err := Unzip(fp, tmpDir)
	if nil != err {
		err = fmt.Errorf("error unzipping %s; %w", fp, err)
		return EXIT_ERROR, err
	}

	// get the details of the update
	// the update "config" is "updtdetails.udt"
	// the "files" are the updated files
	udt, updates, err := GetUpdateDetails(files)
	if nil != err {
		return EXIT_ERROR, err
	}

	// backup the existing files that will be overwritten by the update
	backupDir, err := BackupFiles(updates, instDir)
	if nil != err {
		return EXIT_ERROR, err
	}

	err = InstallUpdate(udt, updates, instDir)
	if nil != err {
		err = fmt.Errorf("error applying update; %w", err)
		// TODO start services after the rollback
		RollbackFiles(backupDir, instDir)
		return EXIT_ERROR, err
	} else {
		return EXIT_SUCCESS, nil
	}
}

func IsUpdateAvailable(infoer Infoer, args Args) (int, error) {
	// read WYC
	iuc, err := infoer.ParseWYC(args.Cdata)
	if nil != err {
		err = fmt.Errorf("error reading %s; %w", args.Cdata, err)
		return EXIT_ERROR, err
	}

	tmpDir := findTempDir()
	defer os.RemoveAll(tmpDir)

	wysTmpFile := path.Join(tmpDir, "wysTemp")
	urls := GetWYSURLs(iuc, args)

	// TODO loop through URLs here or in DownloadFile()
	err = DownloadFile(urls[0], wysTmpFile)
	if nil != err {
		err = fmt.Errorf("error downloading wys file; %w", err)
		return EXIT_ERROR, err
	}

	wys, err := infoer.ParseWYS(wysTmpFile, args)
	if nil != err {
		err = fmt.Errorf("error reading %s; %w", wysTmpFile, err)
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
