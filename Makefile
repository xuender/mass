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
	echo ${VERSION} > cmd/version.txt
	echo ${BUILD_TIME} > cmd/build.txt
	CGO_ENABLED=0 go build \
  -o dist/mass main.go

dy01:
	go run main.go del -d dy01 "DELETE LOW_PRIORITY QUICK FROM refund_business WHERE apply_time < '2023-05-23'"
dy02:
	go run main.go del -d dy02 "DELETE LOW_PRIORITY QUICK FROM refund_business WHERE apply_time < '2023-05-23'"
dy03:
	go run main.go del -d dy03 "DELETE LOW_PRIORITY QUICK FROM refund_business WHERE apply_time < '2023-05-23'"
dy04:
	go run main.go del -d dy04 "DELETE LOW_PRIORITY QUICK FROM refund_business WHERE apply_time < '2023-05-23'"
dy05:
	go run main.go del -d dy05 "DELETE LOW_PRIORITY QUICK FROM refund_business WHERE apply_time < '2023-05-23'"
dy06:
	go run main.go del -d dy06 "DELETE LOW_PRIORITY QUICK FROM refund_business WHERE apply_time < '2023-05-23'"
dy07:
	go run main.go del -d dy07 "DELETE LOW_PRIORITY QUICK FROM refund_business WHERE apply_time < '2023-05-23'"
dy08:
	go run main.go del -d dy08 "DELETE LOW_PRIORITY QUICK FROM refund_business WHERE apply_time < '2023-05-23'"
dy09:
	go run main.go del -d dy09 "DELETE LOW_PRIORITY QUICK FROM refund_business WHERE apply_time < '2023-05-23'"
dy10:
	go run main.go del -d dy10 "DELETE LOW_PRIORITY QUICK FROM refund_business WHERE apply_time < '2023-05-23'"
dy11:
	go run main.go del -d dy11 "DELETE LOW_PRIORITY QUICK FROM refund_business WHERE apply_time < '2023-05-23'"