mkfile_path := $(abspath $(lastword $(MAKEFILE_LIST)))
current_dir := $(notdir $(patsubst %/,%,$(dir $(mkfile_path))))
ROOT_DIR:=$(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))

#EXECUTABLES = git go find pwd
#K := $(foreach exec,$(EXECUTABLES),\
#        $(if $(shell which $(exec)),some string,$(error "No $(exec) in PATH)))

APP_NAME=tcactl
BINARY=tcactl
VERSION=0.6.0
PLATFORMS=darwin linux windows
ARCHITECTURES=386 amd64

LDFLAGS=-ldflags "-X main.Version=${VERSION} -X main.Build=${BUILD}"

default: build

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

## Build the release and development container. The development
docker-build:
	docker-compose build --no-cache $(APP_NAME)
	docker-compose run $(APP_NAME) grunt build
	docker build -t $(APP_NAME)
#
docker-run:
	docker run -i -t --rm --env-file=./config.env -p=$(PORT):$(PORT) --name="$(APP_NAME)" $(APP_NAME)
