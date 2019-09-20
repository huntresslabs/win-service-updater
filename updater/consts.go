package updater

const (
	CLIENT_WYC         = "client.wyc"
	IUCLIENT_IUC       = "iuclient.iuc" // inside client.wyc
	IUC_HEADER         = "IUCDFV2"
	WYS_HEADER         = "IUSDFV2"
	UPDTDETAILS_UDT    = "updtdetails.udt" // inside .wyu archive
	UPDTDETAILS_HEADER = "IUUDFV2"
)

type Infoer interface {
	ParseWYS(string, Args) (ConfigWYS, error)
	ParseWYC(string) (ConfigIUC, error)
}

type Info struct{}