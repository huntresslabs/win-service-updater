package updater

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
)

func GetUpdateDetails(extractedFiles []string) (udt ConfigUDT, updates []string, err error) {
	udtFound := false

	for _, f := range extractedFiles {
		if path.Base(f) == "updtdetails.udt" {
			udt, err = ParseUDT(f)
			if err != nil {
				return ConfigUDT{}, updates, err
			}
			udtFound = true
		} else {
			updates = append(updates, f)
		}
	}

	if !udtFound {
		err := fmt.Errorf("no udt file found")
		return ConfigUDT{}, updates, err
	}

	return udt, updates, nil
}

func InstallUpdate(udt ConfigUDT, srcFiles []string, installDir string) error {
	// stop services
	for _, s := range udt.ServiceToStopBeforeUpdate {
		fmt.Printf("Stopping %s\n", ValueToString(&s))
	}

	for _, f := range srcFiles {
		fmt.Printf("Moving %s\n", f)
		err := MoveFile(f, installDir)
		if err != nil {
			return err
		}
	}

	// start services
	for _, s := range udt.ServiceToStartAfterUpdate {
		fmt.Printf("Starting %s\n", ValueToString(&s))
	}

	return nil
}

func MoveFile(file string, dstDir string) error {
	dst := filepath.Join(dstDir, filepath.Base(file))
	fmt.Println(dst)
	// Rename() returns *LinkError
	return os.Rename(file, dst)
}
