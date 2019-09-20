package updater

import (
	"crypto/rsa"
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
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

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

func Setup() (tmpDir string, tmpFile string) {
	tmpDir, err := ioutil.TempDir("", "prefix")
	if err != nil {
		log.Fatal(err)
	}

	tmpFile, err = ioutil.TempDir("", "prefix")
	if err != nil {
		log.Fatal(err)
	}

	return tmpDir, tmpFile
}

func GetTmpDir() (tmpDir string) {
	tmpDir, err := ioutil.TempDir("", "prefix")
	if err != nil {
		log.Fatal(err)
	}
	return tmpDir
}

func fixupTestURL(uri string, testURL string) string {
	// we need the port from the test server
	tsURI, err := url.ParseRequestURI(testURL)
	if nil != err {
		log.Fatal(err)
	}
	port := tsURI.Port()

	// add the port from the test server url to the url in the wys config
	u, err := url.ParseRequestURI(uri)
	if nil != err {
		log.Fatal(err)
	}
	u.Host = fmt.Sprintf("%s:%s", u.Host, port)
	return u.String()
}

// Test functions

func TestFunctional_CompareVersions(t *testing.T) {
	info := Info{}
	wysFile := "../test_files/widgetX.1.0.1.wys"

	argv := []string{"-urlargs=12345:67890"}
	args := ParseArgs(argv)

	wys, err := info.ParseWYS(wysFile, args)
	assert.Nil(t, err)

	rc := CompareVersions("0.1.2.3", wys.VersionToUpdate)
	assert.Equal(t, A_LESS_THAN_B, rc)
}

func TestFunctional_SameVersion(t *testing.T) {
	info := Info{}
	wycFile := "../test_files/client.1.0.1.wyc"
	wysFile := "../test_files/widgetX.1.0.1.wys"

	tmpDir, instDir := Setup()
	defer os.RemoveAll(tmpDir)
	defer os.RemoveAll(instDir)

	// wys server
	tsWYS := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		dat, err := ioutil.ReadFile(wysFile)
		assert.Nil(t, err)
		w.Write(dat)
	}))
	defer tsWYS.Close()

	argv := []string{fmt.Sprintf(`-cdata="%s"`, wycFile)}
	args := ParseArgs(argv)

	iuc, err := info.ParseWYC(wycFile)
	assert.Nil(t, err)

	uri := fixupTestURL(string(iuc.IucServerFileSite[0].Value), tsWYS.URL)

	fp := fmt.Sprintf("%s/wys", tmpDir)
	err = DownloadFile([]string{uri}, fp)
	assert.Nil(t, err)

	wys, err := info.ParseWYS(fp, args)
	assert.Nil(t, err)

	// fmt.Println("installed ", string(iuc.IucInstalledVersion.Value))
	// fmt.Println("new ", wys.VersionToUpdate)
	rc := CompareVersions(string(iuc.IucInstalledVersion.Value), wys.VersionToUpdate)
	assert.Equal(t, A_EQUAL_TO_B, rc)
}

func TestFunctional_URLArgs(t *testing.T) {
	info := Info{}
	wycFile := "../test_files/client.1.0.0.wyc"
	wysFile := "../test_files/widgetX.1.0.1.wys"
	wyuFile := "../test_files/widgetX.1.0.1.wyu"

	auth := "12345:67890"

	tmpDir, instDir := Setup()
	defer os.RemoveAll(tmpDir)
	defer os.RemoveAll(instDir)

	// test server
	tsWYS := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Contains(t, r.URL.String(), auth)
		w.WriteHeader(http.StatusOK)
		dat, err := ioutil.ReadFile(wysFile)
		assert.Nil(t, err)
		w.Write(dat)
	}))
	defer tsWYS.Close()

	tsWYU := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Contains(t, r.URL.String(), auth)
		w.WriteHeader(http.StatusOK)
		dat, err := ioutil.ReadFile(wyuFile)
		assert.Nil(t, err)
		w.Write(dat)
	}))
	defer tsWYU.Close()

	argv := []string{fmt.Sprintf("-urlargs=%s", auth)}
	args := ParseArgs(argv)

	iuc, err := info.ParseWYC(wycFile)
	assert.Nil(t, err)

	urls := GetWYSURLs(iuc, args)

	// fixup URL adding port from test server
	turi := fixupTestURL(urls[0], tsWYS.URL)

	fp := fmt.Sprintf("%s/wys", tmpDir)
	err = DownloadFile([]string{turi}, fp)
	assert.Nil(t, err)

	wys, err := info.ParseWYS(fp, args)
	assert.Nil(t, err)

	// fmt.Println("installed ", string(iuc.IucInstalledVersion.Value))
	// fmt.Println("new ", wys.VersionToUpdate)
	rc := CompareVersions(string(iuc.IucInstalledVersion.Value), wys.VersionToUpdate)
	assert.Equal(t, A_LESS_THAN_B, rc)

	turi = fixupTestURL(wys.UpdateFileSite[0], tsWYU.URL)

	// download wyu
	fp = fmt.Sprintf("%s/wyu", tmpDir)
	err = DownloadFile([]string{turi}, fp)
	assert.Nil(t, err)
}

func TestFunctional_UpdateWithRollback(t *testing.T) {
	info := Info{}
	wycFile := "../test_files/client.1.0.0.wyc"
	wysFile := "../test_files/widgetX.1.0.1.wys"
	wyuFile := "../test_files/widgetX.1.0.1.wyu"

	auth := "12345:67890"

	tmpDir, instDir := Setup()
	defer os.RemoveAll(tmpDir)
	defer os.RemoveAll(instDir)

	// test server
	tsWYS := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		dat, err := ioutil.ReadFile(wysFile)
		assert.Nil(t, err)
		w.Write(dat)
	}))
	defer tsWYS.Close()

	tsWYU := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		dat, err := ioutil.ReadFile(wyuFile)
		assert.Nil(t, err)
		w.Write(dat)
	}))
	defer tsWYU.Close()

	argv := []string{fmt.Sprintf("-urlargs=%s", auth)}
	args := ParseArgs(argv)

	iuc, err := info.ParseWYC(wycFile)
	assert.Nil(t, err)

	// fixup URL adding port from test server
	turi := fixupTestURL(string(iuc.IucServerFileSite[0].Value), tsWYS.URL)

	fp := fmt.Sprintf("%s\\wys", tmpDir)
	err = DownloadFile([]string{turi}, fp)
	assert.Nil(t, err)

	wys, err := info.ParseWYS(fp, args)
	assert.Nil(t, err)

	// fmt.Println("installed ", string(iuc.IucInstalledVersion.Value))
	// fmt.Println("new ", wys.VersionToUpdate)
	rc := CompareVersions(string(iuc.IucInstalledVersion.Value), wys.VersionToUpdate)
	assert.Equal(t, A_LESS_THAN_B, rc)

	// download wyu
	fp = fmt.Sprintf("%s\\wyu", tmpDir)
	err = DownloadFile([]string{tsWYU.URL}, fp)
	assert.Nil(t, err)

	key, err := ParsePublicKey(string(iuc.IucPublicKey.Value))
	var rsa rsa.PublicKey
	rsa.N = key.Modulus
	rsa.E = key.Exponent

	sha1hash, err := Sha1Hash(fp)
	assert.Nil(t, err)

	// validated
	err = VerifyHash(&rsa, sha1hash, wys.FileSha1)
	assert.Nil(t, err)

	// adler32
	if wys.UpdateFileAdler32 != 0 {
		v := VerifyAdler32Checksum(wys.UpdateFileAdler32, fp)
		assert.True(t, v)
	}

	// extract wyu to tmpDir
	_, files, err := Unzip(fp, tmpDir)
	assert.Nil(t, err)

	udt, updates, err := GetUpdateDetails(files)
	assert.Nil(t, err)

	// the udt should specify stopping/starting the Spooler
	assert.Equal(t, string(udt.ServiceToStopBeforeUpdate[0].Value), "Spooler")
	assert.Equal(t, string(udt.ServiceToStartAfterUpdate[0].Value), "Spooler")

	// make the file that will be replaced
	err = ioutil.WriteFile(path.Join(instDir, "WidgetX.txt"), []byte("1.0.0"), 0644)
	assert.Nil(t, err)

	backupDir, err := BackupFiles(updates, instDir)
	assert.Nil(t, err)

	udt.ServiceToStopBeforeUpdate = []TLV{}
	udt.ServiceToStartAfterUpdate = []TLV{}
	err = InstallUpdate(udt, updates, instDir)
	assert.Nil(t, err)

	// read our "update"
	dat, err := ioutil.ReadFile(path.Join(instDir, "WidgetX.txt"))
	assert.Nil(t, err)
	assert.Equal(t, "1.0.1", string(dat))

	// rollback
	err = RollbackFiles(backupDir, instDir)
	assert.Nil(t, err)

	// original file should be restored
	dat, err = ioutil.ReadFile(path.Join(instDir, "WidgetX.txt"))
	assert.Nil(t, err)
	assert.Equal(t, "1.0.0", string(dat))
}
