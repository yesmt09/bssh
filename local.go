package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
)

//本地配置文件模式
func LocalModel() {
	fmt.Println("local config mode")
	if ConfFile == "" {
		flag.Usage()
		os.Exit(1)
	}
	if _, err := os.Stat(ConfFile); os.IsNotExist(err) {
		flag.Usage()
		os.Exit(1)
	}
	fConfFile, _ := os.Open(ConfFile)
	defer fConfFile.Close()
	json.NewDecoder(fConfFile).Decode(&Config)
}