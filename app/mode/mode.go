package mode

import (
	"bssh/app/mode/etcd"
	"bssh/app/mode/local"
	"bssh/conf"
	"fmt"
)

// 配置的接口
type IModeConf interface {
	GetConfig() (conf.Base, error)
	Watch(*conf.Base)
}

// 获取配置
func GetConfig(cFlag conf.CFlags) (confBase *conf.Base) {
	conf.MainLogger.Info(fmt.Sprintf("%s mode", cFlag.Mode))
	var mode IModeConf
	// 判断启动配置文件方式
	switch cFlag.Mode {
	case "etcd":
		mode = etcd.NewConf(cFlag.RemoteConf)
	case "local":
		mode = local.NewConf(cFlag.LocalConf)
	default:
		panic("err mode")
	}
	_confBase, err := mode.GetConfig()

	if err != nil {
		panic("err")
	}
	confBase = &_confBase
	//运行监听远程配置变化更新程序
	mode.Watch(confBase)
	return
}
