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

type FakeUpdateInfo struct {
	ConfigWYS_FileSha1 []byte
	Err                error
}

func (fakeier FakeUpdateInfo) ParseWYC(wycFile string) (iuc ConfigIUC, err error) {
	info := Info{}

	iuc, err = info.ParseWYC(wycFile)

	return iuc, err
}

func (fakeier FakeUpdateInfo) ParseWYS(wysFile string, args Args) (wys ConfigWYS, err error) {
	info := Info{}

	wys, err = info.ParseWYS(wysFile, args)

	wys.FileSha1 = fakeier.ConfigWYS_FileSha1

	return wys, err
}

func SetupTmpLog() *os.File {
	tmpFile, err := ioutil.TempFile("", "tmpLog")
	if err != nil {
		log.Fatal(err)
	}
	return tmpFile
}

func TearDown(f string) {
	err := os.Remove(f)
	if err != nil {
		log.Fatal(err)
	}
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
	f := SetupTmpLog()
	args.OutputinfoLog = f.Name()
	defer TearDown(args.OutputinfoLog)
	defer f.Close()

	info := Info{}

	exitCode, err := UpdateHandler(info, args)
	assert.Equal(t, EXIT_NO_UPDATE, exitCode)
	assert.Nil(t, err)
	assert.True(t, fileExists(args.OutputinfoLog))
}

func TestUpdateHandler_NoSignedHash(t *testing.T) {
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

	finfo := FakeUpdateInfo{
		ConfigWYS_FileSha1: make([]byte, 0),
	}

	exitCode, err := UpdateHandler(finfo, args)
	assert.Equal(t, EXIT_ERROR, exitCode)
	assert.NotNil(t, err)
}

func TestIsUpdateAvailable_NoUpdate(t *testing.T) {
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
	f := SetupTmpLog()
	args.OutputinfoLog = f.Name()
	defer TearDown(args.OutputinfoLog)
	defer f.Close()

	info := Info{}

	exitCode, _ := IsUpdateAvailable(info, args)
	assert.Equal(t, exitCode, EXIT_NO_UPDATE)
	assert.True(t, fileExists(args.OutputinfoLog))
}

func TestIsUpdateAvailable_ErrorBadWYCFile(t *testing.T) {
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
	f := SetupTmpLog()
	args.OutputinfoLog = f.Name()
	defer TearDown(args.OutputinfoLog)
	defer f.Close()

	info := Info{}

	exitCode, _ := IsUpdateAvailable(info, args)
	assert.Equal(t, exitCode, EXIT_ERROR)
	assert.True(t, fileExists(args.OutputinfoLog))
}

func TestIsUpdateAvailable_ErrorHTTP(t *testing.T) {
	// wycFile := "../test_files/client.1.0.1.wyc"
	// wysFile := "../test_files/widgetX.1.0.1.wys"

	var args Args
	f := SetupTmpLog()
	args.OutputinfoLog = f.Name()
	defer TearDown(args.OutputinfoLog)
	defer f.Close()

	info := Info{}

	exitCode, _ := IsUpdateAvailable(info, args)
	assert.Equal(t, exitCode, EXIT_ERROR)
	assert.True(t, fileExists(args.OutputinfoLog))
}

func TestIsUpdateAvailable_ErrorBadWYSFile(t *testing.T) {
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
	f := SetupTmpLog()
	args.OutputinfoLog = f.Name()
	defer TearDown(args.OutputinfoLog)
	defer f.Close()

	info := Info{}

	exitCode, _ := IsUpdateAvailable(info, args)
	assert.Equal(t, exitCode, EXIT_ERROR)
	assert.True(t, fileExists(args.OutputinfoLog))
}
