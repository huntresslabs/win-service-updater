package updater

import (
	"strings"
)

// Args we support
// "/quickcheck"
// "/justcheck"
// "/noerr",
// "-urlargs="
// "/outputinfo="
// "/fromservice"
// "-logfile="
// "-cdata="

// https://wyday.com/wybuild/help/wyupdate-commandline.php

// /outputinfo[="<filename>"]
// When used together with "/quickcheck /justcheck", "/quickcheck /noerr", or
// even "/quickcheck /justcheck /noerr" arguments wyUpdate will output the
// update information or error to the STDOUT or to a file. For example:

type Args struct {
	Quickcheck    bool
	Justcheck     bool
	Noerr         bool
	Fromservice   bool
	Urlargs       string
	Outputinfo    bool
	OutputinfoLog string
	Logfile       string
	Cdata         string
	Server        string
}

func ParseArgs(argsSlice []string) Args {
	var args Args

	// default to client.wyc
	args.Cdata = "client.wyc"

	for _, arg := range argsSlice {
		larg := strings.ToLower(arg)

		switch {
		case larg == "/quickcheck":
			args.Quickcheck = true
		case larg == "/justcheck":
			args.Justcheck = true
		case larg == "/noerr":
			args.Noerr = true
		case larg == "/fromservice":
			args.Fromservice = true
		case strings.HasPrefix(larg, "-urlargs="):
			fields := strings.Split(larg, "=")
			args.Urlargs = fields[1]
		case strings.HasPrefix(larg, "-logfile="):
			fields := strings.Split(larg, "=")
			args.Logfile = fields[1]
		case strings.HasPrefix(larg, "/outputinfo"):
			args.Outputinfo = true
			if strings.Contains(larg, "=") {
				fields := strings.Split(larg, "=")
				args.OutputinfoLog = fields[1]
			}
		case strings.HasPrefix(larg, "-cdata="):
			fields := strings.Split(larg, "=")
			args.Cdata = fields[1]
		case strings.HasPrefix(larg, "-server="):
			fields := strings.Split(larg, "=")
			args.Server = fields[1]
		}
	}
	return args
}
