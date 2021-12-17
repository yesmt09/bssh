package app

import (
	"bssh/app/mode"
	"bssh/conf"
	"bssh/helper"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/gliderlabs/ssh"
)

var (
	BSshServer *ssh.Server
	baseConf   = &conf.Base{}
	cFlag      conf.CFlags
)

// Run 运行程序
func Run(cFlags conf.CFlags) {
	cFlag = cFlags
	if len(cFlag.Listen) == 0 {
		fmt.Println("listen address must setup")
		flag.Usage()
		os.Exit(1)
	}
	//获取配置文件
	baseConf = mode.GetConfig(cFlag)
	conf.MainLogger.Info(baseConf)
	//注册handle
	ssh.Handle(serverHandle)
	BSshServer = &ssh.Server{
		Addr: cFlag.Listen,
	}
	conf.MainLogger = helper.GetLog(cFlag.LogFile)
	conf.MainLogger.Info(fmt.Sprintf("listen :%s start ok", cFlag.Listen))
	conf.MainLogger.Info(baseConf)
	conf.MainLogger.Flush()

	BSshServer.SetOption(ssh.PasswordAuth(func(ctx ssh.Context, password string) bool {
		return baseConf.Password == password
	}))

	BSshServer.SetOption(ssh.NoPty())

	err := BSshServer.ListenAndServe()
	if err != nil && err.Error() != "ssh: Server closed" {
		fmt.Println(err)
		os.Exit(1)
	}
}

// serverHandle 注册处理的handle
func serverHandle(session ssh.Session) {
	repMsg := "error"
	repCode := 1
	//初始化log
	logger := helper.GetLog(cFlag.LogFile)
	logger.AddBase("user", session.User())
	logger.AddBase("ip", session.RemoteAddr().String())

	// 非配置文件用户不允许使用
	if session.User() != baseConf.User {
		repMsg = fmt.Sprintf("Not Allow User:%s requested.", session.User())
		repCode = 1
	} else {
		commandString := session.RawCommand()
		logger.Info(commandString)
		//在用&切割
		var commandList []string
		if strings.Index(commandString, "&&") != 0 {
			commandList = strings.Split(commandString, "&&")
		} else {
			commandList = append(commandList, commandString)
		}
		var _stdout []byte
		var err error
		for _, commands := range commandList {
			//过滤命令，有违规的命令则不执行
			_, err = helper.FilterCommand(commands, baseConf)
			if err != nil {
				repMsg = err.Error()
				break
			}
		}
		if err == nil {
			cmd := exec.Command("bash","-c", commandString)
			stdout, err := cmd.Output()
			if err != nil {
				stdout = []byte(err.Error())
			}
			_stdout = append(_stdout, stdout...)
			cmd.Wait()
			repMsg = string(_stdout)
			repCode = 0
		}
	}
	logger.Info(repMsg)
	logger.Flush()
	helper.PrintStdout(session, repMsg+"\n", repCode)
}
