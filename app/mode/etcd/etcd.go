package etcd

import (
	"bssh/conf"
	"bssh/helper"
	"context"
	"encoding/json"
	"errors"
	"flag"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"os"
	"strconv"
)

type etcdConf struct {
	conf          conf.RemoteFlag
	client        *clientv3.Client
	configVersion *int64
}

func NewConf(conf conf.RemoteFlag) *etcdConf {
	if conf.Path == "" || conf.HostList == "" {
		flag.Usage()
		os.Exit(1)
	}
	client := helper.ConnectToEtcd(conf.HostList, conf.U, conf.P)
	return &etcdConf{
		conf:   conf,
		client: client,
	}
}

func (e *etcdConf) GetConfig() (conf.Base, error) {
	value, version := helper.GetRemoteConfigFormEtcd(e.client, e.conf.Path)
	var baseConf conf.Base
	err := json.Unmarshal(value, &baseConf)
	if err != nil {
		conf.MainLogger.Fatal("etcd conf json decode err")
		conf.MainLogger.Flush()
		return conf.Base{}, errors.New("etcd conf error")
	}
	version = version + 1
	e.configVersion = &version
	return baseConf, nil

}

func (e *etcdConf) Watch(baseConf *conf.Base) {
	watcher := clientv3.NewWatcher(e.client)
	watchRespChan := watcher.Watch(context.TODO(), e.conf.Path, clientv3.WithRev(*e.configVersion))
	go func() {
		for watchResp := range watchRespChan {
			for _, event := range watchResp.Events {
				switch event.Type {
				case mvccpb.PUT:
					var err error
					*baseConf, err = e.GetConfig()
					conf.MainLogger.Info(baseConf)
					if err != nil {
						break
					}
					conf.MainLogger.Info("conf update,next version:" + strconv.FormatInt(*e.configVersion,10))
					conf.MainLogger.Flush()
					watchRespChan = watcher.Watch(context.TODO(), e.conf.Path, clientv3.WithRev(*e.configVersion))
				case mvccpb.DELETE:
					conf.MainLogger.Fatal("etcd配置文件被删除")
				}
			}
		}
	}()
}
