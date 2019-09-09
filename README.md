# win-service-updater

Partial implementation of wyUpdate functionality written in GoLang.

Goals:
- Compatibility with existing wyUpdate binary files
- Drop in replacement for existing commands (for service updates only)

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
- ~~command-line argument parsing~~
- ~~parse wys file~~
- check it update is available
  - `rc, _ := try(WYUPDATE_EXE, "/quickcheck", "/justcheck", "/noerr", fmt.Sprintf("-urlargs=%s", AUTH), fmt.Sprintf("/outputinfo=%s", CHECK_LOG))`
  - compare current version with available version
  - return code 2 means update is avaiable
- download and extract update to temp folder
  - `wyupdateArgs := fmt.Sprintf("/fromservice -logfile=\"%s\" -urlargs=%s", WYUPDATE_LOG, AUTH)`
- verify signature of update (functions written)
- ~~parse update details~~
- stop services specified in udt (update details) (functions written)
- apply update
- rollback updade
- start services specified in udt (update details) (functions written)
- update client.wyc with new version number

