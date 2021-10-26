package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"os"
	"strings"
	"time"
)

var (
	client       *clientv3.Client
	clientConfig clientv3.Config
	etcdVersion  int64
)

func EtcdModel() {
	if Path == "" || HostList == "" {
		flag.Usage()
		os.Exit(1)
	}
	clientConfig = clientv3.Config{
		Endpoints:   strings.Split(HostList, ","),
		DialTimeout: 5 * time.Second,
	}
	if u != "" {
		clientConfig.Username = u
	}
	if p != "" {
		clientConfig.Password = p
	}
	connectEtcd()
	switch RunModel {
	case "set":
		etcdSetOp()
		os.Exit(0)
		break
	case "get":
		etcdGetOp()
		os.Exit(0)
		break
	case "run":
		etcdRun()
		etcdDataWatch()
		break
	default:
		flag.Usage()
		os.Exit(1)
	}

}

// 连接etcd服务
func connectEtcd() {
	var err error
	if client, err = clientv3.New(clientConfig); err != nil {
		panic(fmt.Sprintln(err))
	}
}

//etcd 运行模式
func etcdRun() {
	fmt.Println("etcd mode")
	kv := clientv3.NewKV(client)
	getResp, err := kv.Get(context.TODO(), Path)
	if err != nil {
		panic(fmt.Sprintln(err))
	}
	err = json.Unmarshal(getResp.Kvs[0].Value, &Config)
	if err != nil {
		panic("etcd config json decode err")
	}
	etcdVersion = getResp.Header.GetRevision() + 1
}

//监听etcd 配置文件变化
func etcdDataWatch() {
	watcher := clientv3.NewWatcher(client)
	watchRespChan := watcher.Watch(context.TODO(), Path, clientv3.WithRev(etcdVersion))
	go func() {
		for watchResp := range watchRespChan {
			for _, event := range watchResp.Events {
				switch event.Type {
				case mvccpb.PUT:
					Wg.Add(1)
					etcdRun()
					BsshServer.Shutdown(context.TODO())
					Wg.Done()
				case mvccpb.DELETE:
					panic("删除了etcd配置文件")
				}
			}
		}
	}()
}

// 往etcd set值
func etcdSetOp() {
	if ConfFile == "" {
		panic("file err")
	}
	if _, err := os.Stat(ConfFile); os.IsNotExist(err) {
		panic(fmt.Sprintf("config file not Exist: %v", ConfFile))
	}
	confFile, _ := os.Open(ConfFile)
	defer confFile.Close()
	var fileContent ConfList
	err := json.NewDecoder(confFile).Decode(&fileContent)
	if err != nil {
		panic("err")
	}
	fmt.Println(fileContent)
	value, _ := json.Marshal(fileContent)
	kv := clientv3.NewKV(client)
	_, err = kv.Put(context.TODO(), Path, string(value))
	if err != nil {
		panic(fmt.Sprintln(err))
	}
	fmt.Println("set ok")
}

//从etcd 获取数据
func etcdGetOp() {
	kv := clientv3.NewKV(client)
	getResp, err := kv.Get(context.TODO(), Path)
	if err != nil {
		panic(fmt.Sprintln(err))
	}
	fmt.Printf(string(getResp.Kvs[0].Value))
}
