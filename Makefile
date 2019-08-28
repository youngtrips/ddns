## env
CGO_ENABLED	:= 0
GOARCH		:= amd64
GOOS		:= $(shell uname -s | tr 'A-Z' 'a-z')
GO			:= go
#VERSION 	:= $(shell git describe --tags --dirty="-dev")
VERSION		:= 0.0.1


TAG_NAME=${REPOSITORY}/${PROJECT}:${TAG}

## targets
APPS=$(shell ls cmd)

all: build

build:
	@for APP in $(APPS) ; do \
		echo building $$APP ; \
		CGO_ENABLED=$(CGO_ENABLED) GOOS=$(GOOS) GOARCH=$(GOARCH) $(GO) build -ldflags "-s -w -X gohive/version.Version=$(VERSION)" -o bin/$$APP ./cmd/$$APP; \
	done

clean:
	rm -rf bin
