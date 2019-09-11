package updater

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func init() {
	TRUE_IF_GO_TEST = true
}

func Sha256Hash(filePath string) (string, error) {
	//Initialize variable returnMD5String now in case an error has to be returned
	var sum string

	//Open the passed argument and check for any error
	file, err := os.Open(filePath)
	if err != nil {
		return sum, err
	}

	//Tell the program to call the following function when the current function returns
	defer file.Close()

	//Open a new hash interface to write to
	hash := sha256.New()

	//Copy the file in the hash interface and check for any error
	if _, err := io.Copy(hash, file); err != nil {
		return sum, err
	}

	//Get the 16 bytes hash
	hashInBytes := hash.Sum(nil)[:16]

	//Convert the bytes to a string
	sum = hex.EncodeToString(hashInBytes)
	return sum, nil
}

func TestFunctional(t *testing.T) {
	// test server
	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`body`))
	}))
	defer ts.Close()

	// we need the port from the test server
	tsURI, err := url.ParseRequestURI(ts.URL)
	assert.Nil(t, err)
	port := tsURI.Port()

	argv := []string{"-urlargs=12345:67890"}
	args := ParseArgs(argv)

	wys, err := ParseWys("../test_files/compressed.wys", args)
	assert.Nil(t, err)

	// add the port from the test server url to the url in the wys config
	u, err := url.ParseRequestURI(wys.UpdateFileSite)
	assert.Nil(t, err)
	u.Host = fmt.Sprintf("%s:%s", u.Host, port)
	turi := u.String()

	dir, err := ioutil.TempDir("", "prefix")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(dir)

	fp := fmt.Sprintf("%s/testdownload", dir)
	err = DownloadFile(turi, fp)
	assert.Nil(t, err)
}
