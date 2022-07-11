ROGUE_COLLECTOR_BINARY=rogue-collector
ROGUE_COLLECTOR_SRC=github.com/eviltomorrow/rogue/app/collector

MAINVERSION=$(shell cat ../../version)

GOPATH ?= $(shell go env GOPATH)
# Ensure GOPATH is set before running build process.
ifeq "$(GOPATH)" ""
  $(error Please set the environment variable GOPATH before running `make`)
endif

PATH := ${GOPATH}/bin:$(PATH)
GCFLAGS=-gcflags "all=-trimpath=${GOPATH}"

GITTAG := $(shell git describe --tags --always)
GITSHA := $(shell git rev-parse --short HEAD)
GITBRANCH := $(shell git rev-parse --abbrev-ref HEAD)
BUILDTIME=`date +%FT%T%z`

LDFLAGS=-ldflags "-X main.MainVersion=${MAINVERSION} -X main.GitSha=${GITSHA} -X main.GitTag=${GITTAG} -X main.GitBranch=${GITBRANCH} -X main.BuildTime=${BUILDTIME} -s -w"

release_rogue_collector:
	go build -race ${LDFLAGS} ${GCFLAGS} -o ${ROGUE_COLLECTOR_BINARY} ${ROGUE_COLLECTOR_SRC}