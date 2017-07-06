
all: build

build:
	go build main.go bridge.go config.go utils.go server.go
