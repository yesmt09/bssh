package main

import (
	"bssh/app/mode/local"
	"bssh/conf"
	"bssh/helper"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"os"
)

var (
	set bool
	get bool
	hostList string
	username string
	password	string
	client *clientv3.Client
	file string
	path string
	version string
	BuildDate string
)

func init()  {
	flag.BoolVar(&set, "set", false, "往etcd设置内容")
	flag.BoolVar(&get, "get", false, "获取etcd内容")
	flag.StringVar(&file, "local", "", "本地配置文件")
	flag.StringVar(&path, "path", "", "etcd的配置路径")
	flag.StringVar(&hostList, "h", "", "etcd 主机列表")
	flag.StringVar(&username, "u", "", "用户名")
	flag.StringVar(&password, "p", "", "密码")
	flag.Parse()
}

func main()  {
	client = helper.ConnectToEtcd(hostList, username, password)
	if set {
		setTo()
	} else if get {
		getForm()
	} else {
		panic("set or get ")
	}
}

func setTo()  {
	if file == "" {
		flag.Usage()
		os.Exit(1)
	}
	var localMode = local.NewConf(conf.LocalFlag{ConfFile: file})
	FileConfigContent, err := localMode.GetConfig()
	content, err := json.Marshal(FileConfigContent)
	if err != nil {
		panic(fmt.Sprintln(err))
	}
	kv := clientv3.NewKV(client)
	_, err = kv.Put(context.TODO(), path, string(content))
	if err != nil {
		panic(fmt.Sprintln(err))
	}
	fmt.Println("set ok")
}

func getForm()  {
	if path == "" {
		flag.Usage()
		os.Exit(1)
	}
	config, _ := helper.GetRemoteConfigFormEtcd(client, path)
	fmt.Println(string(config))
}