package updater

import (
	"archive/zip"
	"crypto/sha1"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// GetExeDir returns the directory name of the executable
func GetExeDir() string {
	exe, _ := os.Executable()
	return filepath.Dir(exe)
}

// Sha1Hash returns the SHA1 hash as []byte
func Sha1Hash(filePath string) ([]byte, error) {
	// Open the passed argument and check for any error
	file, err := os.Open(filePath)
	if err != nil {
		return []byte{}, err
	}

	// Tell the program to call the following function when the current function returns
	defer file.Close()

	// Open a new hash interface to write to
	hash := sha1.New()

	// Copy the file in the hash interface and check for any error
	if _, err := io.Copy(hash, file); err != nil {
		return []byte{}, err
	}

	hashInBytes := hash.Sum(nil)

	return hashInBytes, nil
}

func findTempDir() (tempDir string) {
	tempDir = os.Getenv("TEMP")
	if len(tempDir) == 0 {
		windir := os.Getenv("SystemRoot")
		if len(windir) > 1 {
			tempDir = filepath.Join(windir, "temp")
			return tempDir
		} else {
			tempDir, err := ioutil.TempDir("", "updater")
			if nil != err {
				log.Fatal(err)
			}
			return tempDir
		}
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
