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

	wys, err := ParseWYS("../test_files/compressed.wys", args)
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

func TestFunctional_CompareVersions(t *testing.T) {
	file := "../test_files/compressed.wys"
	argv := []string{"-urlargs=12345:67890"}
	args := ParseArgs(argv)

	wys, err := ParseWYS(file, args)
	assert.Nil(t, err)

	rc := CompareVersions("0.1.2.3", wys.VersionToUpdate)
	assert.Equal(t, A_LESS_THAN_B, rc)
}

func TestFunctional_Stub(t *testing.T) {
	wycFile := "../test_files2/client1.0.0.wyc"
	wysFile := "../test_files2/wyserver1.0.1.wys"
	wyuFile := "../test_files2/widget1.1.0.1.wyu"

	// wycFile := "../huntress/client.wyc"
	// wyuFile := "../huntress/huntress-amd64-(0.9.52).wyu"
	// wysFile := "../huntress/0.9.52.wys"

	tmpDir, err := ioutil.TempDir("", "prefix")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	instDir, err := ioutil.TempDir("", "prefix")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(instDir)

	argv := []string{"-urlargs=12345:67890"}
	args := ParseArgs(argv)

	iuc, err := ParseWYC(wycFile)
	assert.Nil(t, err)

	wys, err := ParseWYS(wysFile, args)
	assert.Nil(t, err)

	// fmt.Println("installed ", string(iuc.IucInstalledVersion.Value))
	// fmt.Println("new ", wys.VersionToUpdate)
	rc := CompareVersions(string(iuc.IucInstalledVersion.Value), wys.VersionToUpdate)
	assert.Equal(t, A_LESS_THAN_B, rc)

	key, err := ParsePublicKey(string(iuc.IucPublicKey.Value))
	var rsa rsa.PublicKey
	rsa.N = key.Modulus
	rsa.E = key.Exponent
	// fmt.Printf("exponent %d\n", rsa.E)

	sha1hash, err := Sha1Hash(wyuFile)
	assert.Nil(t, err)

	// spew.Dump(sha1hash)
	// spew.Dump(wys.FileSha1)

	// validated
	err = VerifyHash(&rsa, sha1hash, wys.FileSha1)
	assert.Nil(t, err)

	// adler32
	if wys.UpdateFileAdler32 != 0 {
		v := VerifyAdler32Checksum(wys.UpdateFileAdler32, wyuFile)
		assert.True(t, v)
	}

	// extract wyu to tmpDir
	_, files, err := Unzip(wyuFile, tmpDir)
	assert.Nil(t, err)

	udt, updates, err := GetUpdateDetails(files)
	assert.Nil(t, err)

	err = InstallUpdate(udt, updates, instDir)
	assert.Nil(t, err)

	// read our "update"
	dat, err := ioutil.ReadFile(path.Join(instDir, "Widget1.txt"))
	assert.Nil(t, err)
	assert.Equal(t, "1.0.1", string(dat))
}
