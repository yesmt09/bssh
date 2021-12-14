.PHONY: build clean
BIN_FILE=bssh
BIN_DIR="bin/"
ETCD_TOOL_FILE=etcdTool
VERSION=`git rev-parse --short HEAD`
DATE= `date +%FT%T%z`
OS= `uname`
build: clean
	@go build -o ${BIN_DIR}${BIN_FILE}_${OS} -ldflags "-X bssh/conf.Version=$(VERSION) -X bssh/conf.Build=$(DATE) -w -s" main.go
	@go build -o ${BIN_DIR}${ETCD_TOOL_FILE}_${OS} -ldflags "-X main.Version=$(VERSION) -X main.Build=$(DATE) -w -s" tool/etcd.go
	@echo "go build bssh ok"
clean:
	@echo "clean rm -rf " ${BIN_DIR}"*"
	@go clean
	@rm -rf ${$BIN_DIR}"*"
