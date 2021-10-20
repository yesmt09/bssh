package main

import (
	"encoding/json"
	"fmt"
	"github.com/gliderlabs/ssh"
	"io"
	"os"
	"os/exec"
	"strings"
)

type _config struct {
	Port     string
	Password string
	Command  []string
}

var config _config

var configPath = "./production.json"

func main() {
	_initConfig()
	ssh.Handle(func(session ssh.Session) {
		_, _, isPty := session.Pty()
		if isPty == false {
			commandString := strings.Join(session.Command(), " ")
			i := strings.Index(commandString, "&&")
			commandList := []string{}
			if i != 0 {
				commandList = strings.Split(commandString, "&&")
			} else {
				commandList = append(commandList, commandString)
			}
			var _stdout []byte
			for _, v := range commandList {
				command := strings.Split(v, " ")
				fmt.Println(strings.Join(command[1:], " "))
				cmd := exec.Command(strings.Trim(command[0]," "), command[1:]...)
				fmt.Println(strings.Join(session.Command()[1:], " "))
				stdout, _ := cmd.Output()
				_stdout = append(_stdout, stdout...)
				cmd.Wait()
			}

			io.WriteString(session, string(_stdout))
		} else {
			io.WriteString(session, "No Allow PTY requested.\n")
			session.Exit(1)
		}
	})

	ssh.ListenAndServe("127.0.0.1:"+config.Port, nil)
}

//初始化配置
func _initConfig() {
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic(fmt.Sprintf("config file not Exist: %v", configPath))
	}
	confFile, _ := os.Open(configPath)
	defer confFile.Close()
	json.NewDecoder(confFile).Decode(&config)
}
