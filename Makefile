# make sure you set export GO111MODULE="on"
#
mkfile_path := $(abspath $(lastword $(MAKEFILE_LIST)))
current_dir := $(notdir $(patsubst %/,%,$(dir $(mkfile_path))))
ROOT_DIR:=$(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))

#EXECUTABLES = git go find pwd
#K := $(foreach exec,$(EXECUTABLES),\
#        $(if $(shell which $(exec)),some string,$(error "No $(exec) in PATH)))

APP_NAME=tcactl
BINARY=tcactl
VERSION=0.7.1
PLATFORMS=darwin linux windows
ARCHITECTURES=386 amd64
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test

LDFLAGS=-ldflags "-X main.Version=${VERSION} -X main.Build=${BUILD}"

default: build

clean:
	$(GOCLEAN)

.PHONY: all
all: clean build_all install

build:
	go build ${LDFLAGS} -o tcactl app/main/main.go

build_all:
	$(foreach GOOS, $(PLATFORMS),\
	$(foreach GOARCH, $(ARCHITECTURES), $(shell export GOOS=$(GOOS); export GOARCH=$(GOARCH); go build -v -o $(BINARY)-$(GOOS)-$(GOARCH) app/main/main.go)))

run:
	go run cmd/client/main/main.go

## Stop and remove a running container
docker-run:
	docker run --rm -it -v "$(GOPATH)":/go -w /go/src/github.com/spyroot/tcactl golang:latest go build -o "$(BINARY_UNIX)" -v

## Stop and remove a running container
stop:
	docker stop $(APP_NAME); docker rm $(APP_NAME)

## Build the release and development container.
docker-build:
#	docker-compose build --no-cache $(APP_NAME)
#	docker-compose run $(APP_NAME) grunt build
	docker build -t spyroot/$(APP_NAME):$(VERSION) .

# Docker run
docker-run:
#	docker run -v `pwd`:`pwd` -w `pwd` --name="$(APP_NAME)" -i -t --rm --env-file=./config.env $(APP_NAME) bash
#	docker run -v `pwd`:`pwd` -w `pwd` --name photon_iso_builder --rm -i -t spyroot/photon_iso_builder:1.0 bash
	docker run -v `pwd`:`pwd` -w `pwd` --name="tcactl" --env-file=./config.env --rm -i -t spyroot/$(APP_NAME):$(VERSION) /bin/sh

dep:
	go mod download

update_list:
	go list -u -m all

lint:
	golangci-lint run --enable-all