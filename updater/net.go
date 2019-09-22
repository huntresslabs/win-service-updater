package updater

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/hashicorp/go-multierror"
)

// DownloadFile will download a url to a local file.
func DownloadFile(urls []string, localpath string) error {
	if len(urls) == 0 {
		err := fmt.Errorf("no download urls specified")
		return err
	}

	// Create the file
	out, err := os.Create(localpath)
	if nil != err {
		err = fmt.Errorf("failed to create output file for download; %w", err)
		return err
	}
	defer out.Close()

	// attempt each url in the slice until download succeeds
	var result error
	for _, url := range urls {
		//  GET file, if we fail try next URL, otherwise return success (nil)
		err = HTTPGetFile(url, out)
		if nil != err {
			result = multierror.Append(result, err)
			continue
		} else {
			return nil
		}
	}

	return result
}

func HTTPGetFile(URL string, file *os.File) error {
	resp, err := http.Get(URL)
	if nil != err {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("error downloading file from %s; http status = %d", URL, resp.StatusCode)
		return err
	}

	// Write the body to file
	_, err = io.Copy(file, resp.Body)
	if nil != err {
		return err
	}

	return nil

}
