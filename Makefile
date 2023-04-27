all: build

build:
	go build -o bin/ydb-rolling-restart cmd/main.go
