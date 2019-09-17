package updater

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type FakeUpdateInfoer struct {
	URL string
	Err error
}

func (f FakeUpdateInfoer) ParseWYC(wycFile string) (iuc ConfigIUC, err error) {
	uier := UpdateInfoer{}

	iuc, err = uier.ParseWYC(wycFile)

	iuc.IucServerFileSite[0].Value = []byte(f.URL)

	return iuc, err
}

func TestIsUpdateAvailable_ChangeURL(t *testing.T) {
	var args Args
	wycFile := "../test_files/client.1.0.1.wyc"

	f := FakeUpdateInfoer{
		URL: "TEST_URL",
		Err: nil,
	}

	exitCode := IsUpdateAvailable(f, wycFile, args)
	assert.Equal(t, exitCode, 3)
}

// func TestIsUpdateAvailable_Error(t *testing.T) {
// 	var args Args
// 	args.Noerr = true

// 	uier := UpdateInfoer{}

// 	exitCode := IsUpdateAvailable(uier, "foo", args)
// 	assert.Equal(t, EXIT_ERROR, exitCode)
// }

// func TestIsUpdateAvailable_NoUpdate(t *testing.T) {
// 	wycFile := "../test_files/client.1.0.1.wyc"
// 	wysFile := "../test_files/widgetX.1.0.1.wys"

// 	tmpDir, instDir := Setup()
// 	defer os.RemoveAll(tmpDir)
// 	defer os.RemoveAll(instDir)

// 	// wys server
// 	tsWYS := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		w.WriteHeader(http.StatusOK)
// 		dat, err := ioutil.ReadFile(wysFile)
// 		assert.Nil(t, err)
// 		w.Write(dat)
// 	}))
// 	defer tsWYS.Close()

// 	argv := []string{fmt.Sprintf(`-cdata="%s"`, wycFile)}
// 	args := ParseArgs(argv)

// 	exitCode := IsUpdateAvailable(wycFile, args)
// 	assert.Equal(t, EXIT_NO_UPDATE, exitCode)
// }

// func TestIsUpdateAvailable_UpdateAvaliable(t *testing.T) {
// 	var args Args

// 	wysFile := "../test_files2/client1.0.1.wyc"

// 	// test server
// 	tsWYS := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		w.WriteHeader(http.StatusOK)
// 		dat, err := ioutil.ReadFile(wysFile)
// 		assert.Nil(t, err)
// 		w.Write(dat)
// 	}))
// 	defer tsWYS.Close()

// 	exitCode := IsUpdateAvailable("../test_files2/client1.0.0.wyc", args)
// 	assert.Equal(t, EXIT_UPDATE_AVALIABLE, exitCode)
// }
