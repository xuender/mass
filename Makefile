PACKAGE = github.com/xuender/mass
VERSION = $(shell git describe --tags)
BUILD_TIME = $(shell date +%F' '%T)

default: lint-fix test

tools:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/cespare/reflex@latest

lint:
	golangci-lint run --timeout 60s --max-same-issues 50 ./...

lint-fix:
	golangci-lint run --timeout 60s --max-same-issues 50 --fix ./...

test:
	go test -race -v ./... -gcflags=all=-l -cover

watch-test:
	reflex -t 50ms -s -- sh -c 'go test -race -v ./...'

clean:
	rm -rf dist

proto:
	protoc --go_out=. pb/*.proto

build:
	echo ${VERSION} > build.txt
	echo ${BUILD_TIME} >> build.txt
	CGO_ENABLED=0 go build \
  -o dist/mass main.go
