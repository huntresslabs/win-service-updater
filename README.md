# win-service-updater

Partial implementation of wyUpdate functionality written in GoLang.

Differences:
- designed to only be run from a service or command-line, there is no GUI component
- only full binary replacement (no diff)
- only supports stopping/starting services before/after update

Arguments supported:
- "/quickcheck"
- "/justcheck"
- "/noerr",
- "-urlargs=_args_"
- "/outputinfo=_out_"
- "/fromservice"
- "-logfile=_log_"

Todo:
- check it update is available
  - `rc, _ := try(WYUPDATE_EXE, "/quickcheck", "/justcheck", "/noerr", fmt.Sprintf("-urlargs=%s", AUTH), fmt.Sprintf("/outputinfo=%s", CHECK_LOG))`
  - compare current version with available version
  - return code 2 means update is avaiable
- download and extract update to temp folder
- stop services specified in udt (update details)
- apply update
- rollback updade
- start services specified in udt (update details)
- update client.wyc with new version number

