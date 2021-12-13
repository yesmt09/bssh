package helper

import (
	"bssh/conf"
	"context"
	"errors"
	"fmt"
	"github.com/gliderlabs/ssh"
	"gitlab.babeltime.com/packagist/blogger"
	clientv3 "go.etcd.io/etcd/client/v3"
	"io"
	"strings"
	"time"
)

func GetLog(logFile string) blogger.BLogger {
	Bfile := blogger.NewBFile(logFile, blogger.L_INFO)
	bLogger := blogger.NewBlogger(Bfile)
	blogger.GetLogid()
	return bLogger
}

func FilterCommand(commands string, Config *conf.Base) (commandList []string, err error) {
	commandList = strings.Split(strings.Trim(commands, " "), " ")
	for key, _command := range commandList {
		commandList[key] = strings.Trim(_command, " ")
	}
	err = nil
	if len(Config.Command) != 0 {
		allowCommand := false
		for _, _vv := range Config.Command {
			if _vv == commandList[0] {
				allowCommand = true
				break
			}
		}
		if allowCommand == false {
			err = errors.New(fmt.Sprintf("%s Command Not Allow", commandList[0]))
		}
	}
	return
}

func PrintStdout(session ssh.Session, message string, code int) {
	io.WriteString(session, message)
	session.Exit(code)
}

func ConnectToEtcd(hostList string, U string, P string) (client *clientv3.Client) {
	var err error
	if hostList == "" {
		panic("hostlist 不能为空")
	}
	clientConf := clientv3.Config{
		Endpoints:   strings.Split(hostList, ","),
		DialTimeout: 5 * time.Second,
	}
	if U != "" {
		clientConf.Username = U
	}
	if P != "" {
		clientConf.Password = P
	}
	if client, err = clientv3.New(clientConf); err != nil {
		panic(fmt.Sprintln(err))
	}
	return client
}

func GetRemoteConfigFormEtcd(client *clientv3.Client, path string) ([]byte, int64) {
	kv := clientv3.NewKV(client)
	getResp, err := kv.Get(context.TODO(), path)
	if err != nil {
		panic(fmt.Sprintln(err))
	}
	return getResp.Kvs[0].Value, getResp.Header.GetRevision()
}
