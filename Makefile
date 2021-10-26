.PHONY: build clean
BIN_FILE=bssh
VERSION="v0.0.1"
DATE= `date +%FT%T%z`

build: clean
	@go build -o ${BIN_FILE} -ldflags "-X main.BsshVersion=$(VERSION) -X main.BsshBuildTime=$(DATE) -w -s" bssh.go etcd.go local.go
	@echo "go build bssh ok"
clean:
	@echo "clean rm -rf bssh"
	@go clean
	@rm -rf bssh
