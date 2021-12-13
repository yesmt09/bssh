package conf

import "gitlab.babeltime.com/packagist/blogger"

type Base struct {
	User     string   `json:"user"`
	Password string   `json:"password"`
	Command  []string `json:"command"`
}

type RemoteFlag struct {
	U        string
	P        string
	Path     string
	HostList string
	Get      bool
	Set      bool
}

type LocalFlag struct {
	ConfFile string
}

type CFlags struct {
	Mode       string
	Listen     string
	LogFile    string
	LocalConf  LocalFlag
	RemoteConf RemoteFlag
}

var (
	Version      string
	Build        string
	MainLogger = blogger.BLogger{}
)