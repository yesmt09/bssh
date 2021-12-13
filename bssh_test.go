package main

import (
	"bssh/app/mode"
	"bssh/conf"
	"fmt"
	"testing"
)
func TestGetConfig(t *testing.T) {
	var cFlag = conf.CFlags{
		Mode: "local",
		LocalConf: conf.LocalFlag{
			ConfFile: "conf/dev.json",
		},
	}
	fmt.Println(mode.GetConfig(cFlag))
}