package updater

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

type FakeUpdateInfoer struct {
	URL string
	Err error
}

// func (f FakeUpdateInfoer) ParseWYC(wycFile string) (iuc ConfigIUC, err error) {
// 	uier := UpdateInfoer{}

// 	iuc, err = uier.ParseWYC(wycFile)

// 	iuc.IucServerFileSite[0].Value = []byte(f.URL)

// 	return iuc, err
// }

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	if !os.IsNotExist(err) {
		//file exists
		return true
	}

	// no such file or directory
	return false
}

func SetupTmpLog() string {
	tmpFile, err := ioutil.TempDir("", "prefix")
	if err != nil {
		log.Fatal(err)
	}
	return tmpFile
}

func TearDown(f string) {
	os.Remove(f)
}

func TestUpdateHandler(t *testing.T) {
	wycFile := "../test_files/client.1.0.1.wyc"
	wysFile := "../test_files/widgetX.1.0.1.wys"
	wyuFile := "../test_files/widgetX.1.0.1.wyu"

	// wys server
	tsWYS := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		dat, err := ioutil.ReadFile(wysFile)
		assert.Nil(t, err)
		w.Write(dat)
	}))
	defer tsWYS.Close()

	// wys server
	tsWYU := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		dat, err := ioutil.ReadFile(wyuFile)
		assert.Nil(t, err)
		w.Write(dat)
	}))
	defer tsWYU.Close()

	var args Args
	args.Cdata = wycFile
	args.Server = tsWYS.URL
	args.WYUTestServer = tsWYU.URL
	args.Outputinfo = true
	args.OutputinfoLog = SetupTmpLog()
	defer TearDown((args.OutputinfoLog))

	exitCode, err := UpdateHandler(args)
	assert.Equal(t, exitCode, EXIT_NO_UPDATE)
	assert.Nil(t, err)
	assert.True(t, fileExists(args.OutputinfoLog))
}

func TestCheckUpdateHandler_NoUpdate(t *testing.T) {
	wycFile := "../test_files/client.1.0.1.wyc"
	wysFile := "../test_files/widgetX.1.0.1.wys"

	// wys server
	tsWYS := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		dat, err := ioutil.ReadFile(wysFile)
		assert.Nil(t, err)
		w.Write(dat)
	}))
	defer tsWYS.Close()

	var args Args
	args.Cdata = wycFile
	args.Server = tsWYS.URL
	args.Outputinfo = true
	args.OutputinfoLog = SetupTmpLog()
	defer TearDown((args.OutputinfoLog))

	exitCode := CheckUpdateHandler(args)
	assert.Equal(t, exitCode, EXIT_NO_UPDATE)
	assert.True(t, fileExists(args.OutputinfoLog))
}

func TestCheckUpdateHandler_ErrorBadWYCFile(t *testing.T) {
	wycFile := "../test_files/foo"
	wysFile := "../test_files/widgetX.1.0.1.wys"

	// wys server
	tsWYS := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		dat, err := ioutil.ReadFile(wysFile)
		assert.Nil(t, err)
		w.Write(dat)
	}))
	defer tsWYS.Close()

	var args Args
	args.Cdata = wycFile
	args.Server = tsWYS.URL
	args.Outputinfo = true
	args.OutputinfoLog = SetupTmpLog()
	defer TearDown((args.OutputinfoLog))

	exitCode := CheckUpdateHandler(args)
	assert.Equal(t, exitCode, EXIT_ERROR)
	assert.True(t, fileExists(args.OutputinfoLog))
}

func TestCheckUpdateHandler_ErrorHTTP(t *testing.T) {
	// wycFile := "../test_files/client.1.0.1.wyc"
	// wysFile := "../test_files/widgetX.1.0.1.wys"

	var args Args
	args.OutputinfoLog = SetupTmpLog()
	defer TearDown((args.OutputinfoLog))

	exitCode := CheckUpdateHandler(args)
	assert.Equal(t, exitCode, EXIT_ERROR)
	assert.True(t, fileExists(args.OutputinfoLog))
}

func TestCheckUpdateHandler_ErrorBadWYSFile(t *testing.T) {
	wycFile := "../test_files/client.1.0.1.wyc"

	// wys server
	tsWYS := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("not a wys file"))
	}))
	defer tsWYS.Close()

	var args Args
	args.Cdata = wycFile
	args.Server = tsWYS.URL
	args.OutputinfoLog = SetupTmpLog()
	defer TearDown((args.OutputinfoLog))

	exitCode := CheckUpdateHandler(args)
	assert.Equal(t, exitCode, EXIT_ERROR)
	assert.True(t, fileExists(args.OutputinfoLog))
}
