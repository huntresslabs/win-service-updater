package updater

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func GetExeDir() string {
	exe, _ := os.Executable()
	return filepath.Dir(exe)
}

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

// Unzip will decompress a zip archive, moving all compressed files/folders
// to the specified output directory.
func Unzip(srcArchive string, destDir string) (root string, filenames []string, err error) {
	r, err := zip.OpenReader(srcArchive)
	if err != nil {
		err := fmt.Errorf("OpenReader() failed: %w", err)
		return "", filenames, err
	}
	defer r.Close()

	for _, f := range r.File {

		// Store filename/path for returning and using later on
		fpath := filepath.Join(destDir, f.Name)

		// Check for ZipSlip. More Info: http://bit.ly/2MsjAWE
		if !strings.HasPrefix(fpath, filepath.Clean(destDir)+string(os.PathSeparator)) {
			return "", filenames, fmt.Errorf("%s: illegal file path", fpath)
		}

		err := writeCompressedFile(f, fpath)
		if nil != err {
			return "", filenames, err
		}
		filenames = append(filenames, fpath)
	}
	return "", filenames, nil
}

// writeCompressedFile writes compressed file to `fpath`
func writeCompressedFile(f *zip.File, fpath string) error {
	if f.FileInfo().IsDir() {
		// Make Folder
		os.MkdirAll(fpath, os.ModePerm)
		return nil
	}

	// Make File
	if err := os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
		return err
	}

	outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
	if err != nil {
		err := fmt.Errorf("OpenFile() failed: %w", err)
		return err
	}
	defer outFile.Close()

	rc, err := f.Open()
	if err != nil {
		err := fmt.Errorf("Open() failed: %w", err)
		return err
	}
	defer rc.Close()

	_, err = io.Copy(outFile, rc)
	if err != nil {
		err := fmt.Errorf("Copy() failed: %w", err)
		return err
	}

	return nil
}
