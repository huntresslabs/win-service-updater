package updater

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// functions to decompress .wyu file
// .wyu files contain
// - updtdetails.upt (update details)
// - base/service.exe
// - base/config.ini
// - base/uninstall.exe

func findTempDir() (tempDir string) {
	tempDir = os.Getenv("TEMP")
	if len(tempDir) == 0 {
		windir := os.Getenv("SystemRoot")
		if 0 == len(windir) {
			fmt.Println("No temp directory")
			os.Exit(1)
		}
		tempDir = fmt.Sprintf("%s\\temp", windir)
	}
	return tempDir
}

// Unzip will decompress a zip archive, moving all files and folders
// within the zip file (parameter 1) to an output directory (parameter 2).
func Unzip(src string, dest string) (root string, filenames []string, err error) {
	r, err := zip.OpenReader(src)
	if err != nil {
		err := fmt.Errorf("OpenReader() failed: %w", err)
		return "", filenames, err
	}
	defer r.Close()

	for _, f := range r.File {

		// Store filename/path for returning and using later on
		fpath := filepath.Join(dest, f.Name)

		// Check for ZipSlip. More Info: http://bit.ly/2MsjAWE
		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return "", filenames, fmt.Errorf("%s: illegal file path", fpath)
		}

		filenames = append(filenames, fpath)

		if f.FileInfo().IsDir() {
			// Make Folder
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		// Make File
		if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return "", filenames, err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			err := fmt.Errorf("OpenFile() failed: %w", err)
			return "", filenames, err
		}

		rc, err := f.Open()
		if err != nil {
			err := fmt.Errorf("Open() failed: %w", err)
			return "", filenames, err
		}

		_, err = io.Copy(outFile, rc)

		// Close the file without defer to close before next iteration of loop
		outFile.Close()
		rc.Close()

		if err != nil {
			err := fmt.Errorf("Copy() failed: %w", err)
			return "", filenames, err
		}
	}
	return "", filenames, nil
}
