Project=github.com/pubgo/ossync
GOPath=$(shell go env GOPATH)
Version=$(shell git tag --sort=committerdate | tail -n 1)
GoROOT=$(shell go env GOROOT)
BuildTime=$(shell date "+%F %T")
CommitID=$(shell git rev-parse HEAD)
LDFLAGS=-ldflags " \
-X '${Project}/version.GoROOT=${GoROOT}' \
-X '${Project}/version.BuildTime=${BuildTime}' \
-X '${Project}/version.GoPath=${GOPath}' \
-X '${Project}/version.CommitID=${CommitID}' \
-X '${Project}/version.Project=${Project}' \
-X '${Project}/version.Version=${Version:-v0.0.1}' \
"

.PHONY: build
build:
	@go build ${LDFLAGS} -mod vendor -race -v -o main .

.PHONY: install
install:
	@go install ${LDFLAGS} .

.PHONY: release
release:
	@go build ${LDFLAGS} -race -v -o main .

.PHONY: test
test:
	@go test -race -v ./... -cover

.PHONY: run
run:
	@go run ${LDFLAGS} -mod vendor -race -v . ossync -t
