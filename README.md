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
- "-urlargs=<args>"
- "/outputinfo=<out>"
- "/fromservice"
- "-logfile=<log>"
