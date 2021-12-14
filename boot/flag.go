package boot

import (
	"bssh/conf"
	"bssh/helper"
	"flag"
	"fmt"
	"os"
)

var (
	update       bool
	printVersion bool
)

var helpMessage = `启用方式: 
本地文件模式: 
	bssh -mode local -conf xxx.json 
远程模式: 
	bssh -mode etcd -path xxx -h xxxx:2379,xxxx:2379
参数:
`


//初始化配置
func GetFlag(cFlag *conf.CFlags) {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), helpMessage)
		flag.PrintDefaults()
		fmt.Fprintf(flag.CommandLine.Output(), `version:%s,build: %s`, conf.Version, conf.Build)
	}
	flag.StringVar(&cFlag.Mode, "mode", "local", "local or etcd ...")
	flag.StringVar(&cFlag.LocalConf.ConfFile, "conf", "", "本地配置文件路径")

	flag.StringVar(&cFlag.RemoteConf.Path, "path", "", "远程配置文件路径")
	flag.StringVar(&cFlag.RemoteConf.HostList, "h", "", "远程器列表")
	flag.StringVar(&cFlag.RemoteConf.P, "p", "", "远程认证密码")
	flag.StringVar(&cFlag.RemoteConf.U, "u", "", "远程认证用户名")
	flag.BoolVar(&cFlag.RemoteConf.Set, "set", false, "set 本地文件至远程")
	flag.BoolVar(&cFlag.RemoteConf.Get, "get", false, "get 本地文件至远程")
	flag.StringVar(&cFlag.LogFile, "log", "/home/pirate/log/bssh/access.log", "日志路径")
	flag.StringVar(&cFlag.Listen, "listen", "0.0.0.0:2222", "服务监听地址")
	flag.BoolVar(&update, "update", false, "升级程序")
	flag.BoolVar(&printVersion, "v", false, "版本号")
	flag.Parse()

	getVersionOrUpdate()

	logger := helper.GetLog(cFlag.LogFile)
	logger.Info(fmt.Sprintf("version:%s build: %s", conf.Version, conf.Build))
	logger.Flush()
}

func getVersionOrUpdate() {
	if printVersion == true {
		fmt.Println("version: " + conf.Version)
		fmt.Println("build: " + conf.Build)
		os.Exit(0)
	}
	if update == true {
		//TODO 远程更新程序功能
		os.Exit(0)
	}
}

