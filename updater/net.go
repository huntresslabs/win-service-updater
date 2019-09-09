package updater

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"os"
)

// DownloadFile will download a url to a local file.
func DownloadFile(url string, localpath string) error {
	// Create the file
	out, err := os.Create(localpath)
	if nil != err {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if nil != err {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("Error downloading file, HTTP status = %d", resp.StatusCode)
		return err
	}

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if nil != err {
		return err
	}

	return nil
}
