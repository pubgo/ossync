Project=github.com/pubgo/ossync
GOPath=$(shell go env GOPATH)
Version=$(shell git tag --sort=committerdate | tail -n 1)
GoROOT=$(shell go env GOROOT)
BuildTime=$(shell date "+%F %T")
CommitID=$(shell git rev-parse HEAD)
LDFLAGS=-ldflags " \
-X 'github.com/pubgo/golug/version.GoROOT=${GoROOT}' \
-X 'github.com/pubgo/golug/version.BuildTime=${BuildTime}' \
-X 'github.com/pubgo/golug/version.GoPath=${GOPath}' \
-X 'github.com/pubgo/golug/version.CommitID=${CommitID}' \
-X 'github.com/pubgo/golug/version.Project=${Project}' \
-X 'github.com/pubgo/golug/version.Version=${Version}' \
"

.PHONY: build
build:
	@go build ${LDFLAGS} -mod vendor -race -v -o main .

.PHONY: install
install:
	@go install -v ${LDFLAGS} .

.PHONY: release
release:
	@go build ${LDFLAGS} -race -v -o main .

.PHONY: test
test:
	@go test -race -v ./... -cover

.PHONY: run
run:
	@go run ${LDFLAGS} -mod vendor -race -v . ossync -t
