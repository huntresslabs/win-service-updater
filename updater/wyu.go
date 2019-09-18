package updater

import "strings"

// functions to decompress .wyu file
// .wyu files contain
// - updtdetails.upt (update details)
// - base/service.exe
// - base/config.ini
// - base/uninstall.exe

// DownloadWYUFile()

// ExtractWYUFile()

// GetWYSURLs returns the ServerFileSite(s) listed in the WYC file.
func GetWYUURLs(wys ConfigWYS, args Args) (urls []string) {
	// This can only be specified in tests
	if len(args.WYUTestServer) > 0 {
		urls = append(urls, args.WYUTestServer)
		return urls
	}

	for _, s := range wys.UpdateFileSite {
		u := strings.Replace(s, "%urlargs%", args.Urlargs, 1)
		urls = append(urls, u)
	}
	return urls
}
