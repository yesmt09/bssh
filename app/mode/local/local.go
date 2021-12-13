package local

import (
	"bssh/conf"
	"encoding/json"
	"flag"
	"fmt"
	"os"
)

type localConf struct {
	conf conf.LocalFlag
}

func NewConf(conf conf.LocalFlag) *localConf {
	return &localConf{
		conf: conf,
	}
}

func (l *localConf) GetConfig() (conf.Base, error) {
	if l.conf.ConfFile == "" {
		flag.Usage()
		os.Exit(1)
	}
	openConfFile, err := os.Open(l.conf.ConfFile)
	defer openConfFile.Close()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	var baseConf conf.Base
	err = json.NewDecoder(openConfFile).Decode(&baseConf)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return baseConf, nil
}

func (l *localConf) Watch(confBase *conf.Base) {
}
