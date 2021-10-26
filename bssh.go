package main

import (
	"flag"
	"fmt"
	"github.com/gliderlabs/ssh"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"
)

type ConfList struct {
	Listen   string   `json:"listen"`
	User     string   `json:"user"`
	Password string   `json:"password"`
	Command  []string `json:"command"`
}

var (
	ConfFile      string
	Path          string
	HostList      string
	Config        ConfList
	RunModel      string
	etcd          bool
	u             string
	p             string
	BsshVersion   string
	BsshBuildTime string
	PrintVersion  bool
	Wg            sync.WaitGroup
	BsshServer    *ssh.Server
	Updata        bool
)

//初始化配置
func init() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "启用方式: \n本地文件模式: bssh -file xxx.json \nETCD模式: bssh -etcd -path xxx -h xxxx:2379,xxxx:2379\nETCD路径配置管理: bssh -etcd -path xxx -file xxx.json -h xxxx:2379,xxxx:2379 -runmodel set/get")
		flag.PrintDefaults()
	}
	flag.StringVar(&ConfFile, "file", "", "本地配置文件路径")
	flag.StringVar(&RunModel, "runmodel", "run", "etcd配置模式下，可get，set path路径的内容")
	flag.StringVar(&Path, "path", "", "etct 配置文件路径")
	flag.StringVar(&HostList, "h", "", "etcd 服务器列表")
	flag.BoolVar(&etcd, "etcd", false, "从etcd远程读取配置模式")
	flag.StringVar(&u, "u", "", "etct 认证用户名")
	flag.StringVar(&p, "p", "", "etcd 认证密码")
	flag.BoolVar(&Updata, "update", false, "升级程序")
	flag.BoolVar(&PrintVersion, "v", false, "版本号")
	flag.Parse()

	if PrintVersion == true {
		fmt.Printf("version:%s \nbuild: %s", BsshVersion, BsshBuildTime)
		os.Exit(0)
	}

	if Updata == true {
		//TODO 远程更新程序功能
		os.Exit(0)
	}

	if etcd == true {
		EtcdModel()
	} else {
		LocalModel()
	}
	if Config.User == "" {
		panic("ssh 连接 User 没配置")
	}
	if Config.Password == "" {
		panic("ssh 连接 Password 没配置")
	}
}

func main() {
	for {
		ssh.Handle(func(session ssh.Session) {
			_, _, isPty := session.Pty()
			if session.User() != Config.User {
				printStdout(session, fmt.Sprintf("Not Allow User:%s requested.", session.User()), 1)
			} else if isPty == false {
				commandString := strings.Join(session.Command(), " ")
				iand := strings.Index(commandString, "&&")
				var commandList []string
				if iand != 0 {
					commandList = strings.Split(commandString, "&&")
				} else {
					commandList = append(commandList, commandString)
				}
				var _stdout []byte
				for _, commands := range commandList {
					ipipe := strings.Index(commands, "|")
					pipeCommandList := []string{}
					if ipipe != 0 {
						pipeCommandList = strings.Split(commandString, "|")
					} else {
						
					}

					command := filterCommand(commands, session)
					if len(command) == 0 {
						continue
					}
					cmd := exec.Command(strings.Trim(command[0], " "), command[1:]...)
					stdout, err := cmd.Output()
					if err != nil {
						stdout = []byte(err.Error())
					}
					_stdout = append(_stdout, stdout...)
					cmd.Wait()
				}

				io.WriteString(session, string(_stdout))
				session.Exit(0)
			} else {
				printStdout(session, "Not Allow PTY requested.\n", 1)
			}
		})
		Wg.Add(1)
		go func() {
			defer Wg.Done()
			BsshServer = &ssh.Server{
				Addr: Config.Listen,
			}
			fmt.Printf("run server ok,listen: %s \n", Config.Listen)
			BsshServer.SetOption(ssh.PasswordAuth(func(ctx ssh.Context, password string) bool {
				return Config.Password == password
			}))
			err := BsshServer.ListenAndServe()
			if err != nil && err.Error() != "ssh: Server closed" {
				fmt.Println(err)
				os.Exit(1)
			}
		}()
		Wg.Wait()
		fmt.Println("restart ok")
	}
}

func filterCommand(commands string, session ssh.Session) (command []string) {
	_commands := strings.Split(commands, " ")
	for _, vv := range _commands {
		if vv == "" {
			continue
		}
		command = append(command, strings.Trim(vv, " "))
		if len(Config.Command) != 0 {
			var allowCommand bool = false
			for _, _vv := range Config.Command {
				if _vv == command[0] {
					allowCommand = true
					break
				}
			}
			if allowCommand == false {
				printStdout(session, fmt.Sprintf("%s Command Not Allow", command[0]), 1)
			}
		}
	}
	return
}

func printStdout(session ssh.Session, message string, code int) {
	io.WriteString(session, message)
	session.Exit(code)
}
