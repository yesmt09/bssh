package main

import (
	"encoding/json"
	"fmt"
	"github.com/gliderlabs/ssh"
	"io"
	"os"
	"os/exec"
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
			cmd := exec.Command(session.RawCommand())
			stdout, err := cmd.Output()
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			cmd.Wait()
			io.WriteString(session, string(stdout))
		} else {
			io.WriteString(session, "No Allow PTY requested.\n")
			session.Exit(1)
		}
	})

	ssh.ListenAndServe(":"+config.Port, nil,
		ssh.PasswordAuth(func(ctx ssh.Context, pass string) bool {
			return pass == config.Password
		}),
	)
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
