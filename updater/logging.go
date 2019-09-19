package updater

import (
	"fmt"
	"io/ioutil"
)

// LogErrorMsg will write to a log file if a log file was specified
func LogErrorMsg(args Args, msg string) {
	if len(args.Logfile) > 0 {
		dat := []byte(msg)
		ioutil.WriteFile(args.Logfile, dat, 0644)
	}
}

// LogOutputInfoMsg will write a msg to STDOUT or a log file if one was specified
func LogOutputInfoMsg(args Args, msg string) {
	if args.Outputinfo {
		if len(args.OutputinfoLog) > 0 {
			dat := []byte(msg)
			ioutil.WriteFile(args.OutputinfoLog, dat, 0644)
		} else {
			fmt.Println(msg)
		}
	}
}
