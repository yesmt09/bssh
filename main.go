package main

import (
	"bssh/app"
	"bssh/boot"
	"bssh/conf"
)

var (
	cFlag conf.CFlags
)

func main() {
	boot.GetFlag(&cFlag)
	app.Run(cFlag)
}
