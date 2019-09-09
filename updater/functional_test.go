package updater

import (
	"fmt"
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

	wys := ParseWys("../test_files/wys_uncompressed.bin", args)

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
