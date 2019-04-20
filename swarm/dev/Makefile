.PHONY: build cluster test

default: build

build:
	go build -o bin/swarm    github.com/puffscoin/go-puffscoin/cmd/swarm
	go build -o bin/gpuffs     github.com/puffscoin/go-puffscoin/cmd/geth
	go build -o bin/bootnode github.com/puffscoin/go-puffscoin/cmd/bootnode

cluster: build
	scripts/boot-cluster.sh

test:
	go test -v github.com/puffscoin/go-puffscoin/swarm/...
