package updater

import (
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
		err = fmt.Errorf("http get error %s; %w", url, err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("error downloading file from %s; http status = %d", url, resp.StatusCode)
		return err
	}

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if nil != err {
		return err
	}

	return nil
}
