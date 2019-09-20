package updater

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

// DownloadFile will download a url to a local file.
func DownloadFile(urls []string, localpath string) error {
	// Create the file
	out, err := os.Create(localpath)
	if nil != err {
		return err
	}
	defer out.Close()

	// attempt each url in the slice until download suceeds
	for _, url := range urls {
		resp, err := http.Get(url)
		if nil != err {
			err = fmt.Errorf("http get error %s; %w", url, err)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			err = fmt.Errorf("error downloading file from %s; http status = %d; %w", url, resp.StatusCode, err)
			continue
		}

		// Write the body to file
		_, err = io.Copy(out, resp.Body)
		if nil != err {
			return err
		}

		return nil
	}

	return err
}
