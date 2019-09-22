package updater

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/go-multierror"
	"github.com/stretchr/testify/assert"
)

func TestNet_DownloadFile_Success(t *testing.T) {
	wysFile := "../test_files/widgetX.1.0.1.wys"

	server2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		dat, err := ioutil.ReadFile(wysFile)
		assert.Nil(t, err)
		w.Write(dat)
	}))
	defer server2.Close()

	f := SetupTmpLog()
	err := DownloadFile([]string{server2.URL}, f.Name())
	assert.Nil(t, err)

	origHash, err := Sha256Hash(wysFile)
	assert.Nil(t, err)

	newHash, err := Sha256Hash(f.Name())
	assert.Nil(t, err)
	assert.Equal(t, origHash, newHash)
}

func TestNet_DownloadFile_AllError(t *testing.T) {
	server1 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer server1.Close()

	server2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer server2.Close()

	f := SetupTmpLog()
	err := DownloadFile([]string{server1.URL, server2.URL, "http://foo.bar"}, f.Name())
	assert.NotNil(t, err)
	_, ok := err.(*multierror.Error)
	assert.True(t, ok)
}

func TestNet_DownloadFile_NoURLs(t *testing.T) {
	f := SetupTmpLog()
	err := DownloadFile([]string{}, f.Name())
	assert.NotNil(t, err)
}

func TestNet_DownloadFile_InvalidLocalFile(t *testing.T) {
	err := DownloadFile([]string{"http://foo.bar"}, "/Users/foo")
	assert.NotNil(t, err)
}
