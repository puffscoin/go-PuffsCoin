# This Makefile is meant to be used by people that do not usually work
# with Go source code. If you know what GOPATH is then you probably
# don't need to bother with make.

.PHONY: gpuffs android ios gpuffs-cross swarm evm all test clean
.PHONY: gpuffs-linux gpuffs-linux-386 gpuffs-linux-amd64 gpuffs-linux-mips64 gpuffs-linux-mips64le
.PHONY: gpuffs-linux-arm gpuffs-linux-arm-5 gpuffs-linux-arm-6 gpuffs-linux-arm-7 gpuffs-linux-arm64
.PHONY: gpuffs-darwin gpuffs-darwin-386 gpuffs-darwin-amd64
.PHONY: gpuffs-windows gpuffs-windows-386 gpuffs-windows-amd64

GOBIN = $(shell pwd)/build/bin
GO ?= latest

gpuffs:
	build/env.sh go run build/ci.go install ./cmd/gpuffs
	@echo "Done building."
	@echo "Run \"$(GOBIN)/gpuffs\" to launch gpuffs."

swarm:
	build/env.sh go run build/ci.go install ./cmd/swarm
	@echo "Done building."
	@echo "Run \"$(GOBIN)/swarm\" to launch swarm."

all:
	build/env.sh go run build/ci.go install

android:
	build/env.sh go run build/ci.go aar --local
	@echo "Done building."
	@echo "Import \"$(GOBIN)/gpuffs.aar\" to use the library."

ios:
	build/env.sh go run build/ci.go xcode --local
	@echo "Done building."
	@echo "Import \"$(GOBIN)/gpuffs.framework\" to use the library."

test: all
	build/env.sh go run build/ci.go test

lint: ## Run linters.
	build/env.sh go run build/ci.go lint

clean:
	./build/clean_go_build_cache.sh
	rm -fr build/_workspace/pkg/ $(GOBIN)/*

# The devtools target installs tools required for 'go generate'.
# You need to put $GOBIN (or $GOPATH/bin) in your PATH to use 'go generate'.

devtools:
	env GOBIN= go get -u golang.org/x/tools/cmd/stringer
	env GOBIN= go get -u github.com/kevinburke/go-bindata/go-bindata
	env GOBIN= go get -u github.com/fjl/gencodec
	env GOBIN= go get -u github.com/golang/protobuf/protoc-gen-go
	env GOBIN= go install ./cmd/abigen
	@type "npm" 2> /dev/null || echo 'Please install node.js and npm'
	@type "solc" 2> /dev/null || echo 'Please install solc'
	@type "protoc" 2> /dev/null || echo 'Please install protoc'

swarm-devtools:
	env GOBIN= go install ./cmd/swarm/mimegen

# Cross Compilation Targets (xgo)

gpuffs-cross: gpuffs-linux gpuffs-darwin geth-windows gpuffs-android gpuffs-ios
	@echo "Full cross compilation done:"
	@ls -ld $(GOBIN)/gpuffs-*

gpuffs-linux: gpuffs-linux-386 gpuffs-linux-amd64 gpuffs-linux-arm gpuffs-linux-mips64 gpuffs-linux-mips64le
	@echo "Linux cross compilation done:"
	@ls -ld $(GOBIN)/gpuffs-linux-*

gpuffs-linux-386:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/386 -v ./cmd/gpuffs
	@echo "Linux 386 cross compilation done:"
	@ls -ld $(GOBIN)/gpuffs-linux-* | grep 386

gpuffs-linux-amd64:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/amd64 -v ./cmd/gpuffs
	@echo "Linux amd64 cross compilation done:"
	@ls -ld $(GOBIN)/gpuffs-linux-* | grep amd64

gpuffs-linux-arm: gpuffs-linux-arm-5 gpuffs-linux-arm-6 gpuffs-linux-arm-7 gpuffs-linux-arm64
	@echo "Linux ARM cross compilation done:"
	@ls -ld $(GOBIN)/gpuffs-linux-* | grep arm

gpuffs-linux-arm-5:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/arm-5 -v ./cmd/gpuffs
	@echo "Linux ARMv5 cross compilation done:"
	@ls -ld $(GOBIN)/gpuffs-linux-* | grep arm-5

gpuffs-linux-arm-6:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/arm-6 -v ./cmd/gpuffs
	@echo "Linux ARMv6 cross compilation done:"
	@ls -ld $(GOBIN)/gpuffs-linux-* | grep arm-6

gpuffs-linux-arm-7:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/arm-7 -v ./cmd/gpuffs
	@echo "Linux ARMv7 cross compilation done:"
	@ls -ld $(GOBIN)/gpuffs-linux-* | grep arm-7

gpuffs-linux-arm64:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/arm64 -v ./cmd/gpuffs
	@echo "Linux ARM64 cross compilation done:"
	@ls -ld $(GOBIN)/gpuffs-linux-* | grep arm64

gpuffs-linux-mips:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/mips --ldflags '-extldflags "-static"' -v ./cmd/gpuffs
	@echo "Linux MIPS cross compilation done:"
	@ls -ld $(GOBIN)/gpuffs-linux-* | grep mips

gpuffs-linux-mipsle:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/mipsle --ldflags '-extldflags "-static"' -v ./cmd/gpuffs
	@echo "Linux MIPSle cross compilation done:"
	@ls -ld $(GOBIN)/gpuffs-linux-* | grep mipsle

gpuffs-linux-mips64:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/mips64 --ldflags '-extldflags "-static"' -v ./cmd/gpuffs
	@echo "Linux MIPS64 cross compilation done:"
	@ls -ld $(GOBIN)/gpuffs-linux-* | grep mips64

gpuffs-linux-mips64le:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/mips64le --ldflags '-extldflags "-static"' -v ./cmd/gpuffs
	@echo "Linux MIPS64le cross compilation done:"
	@ls -ld $(GOBIN)/gpuffs-linux-* | grep mips64le

gpuffs-darwin: gpuffs-darwin-386 gpuffs-darwin-amd64
	@echo "Darwin cross compilation done:"
	@ls -ld $(GOBIN)/gpuffs-darwin-*

gpuffs-darwin-386:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=darwin/386 -v ./cmd/gpuffs
	@echo "Darwin 386 cross compilation done:"
	@ls -ld $(GOBIN)/gpuffs-darwin-* | grep 386

gpuffs-darwin-amd64:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=darwin/amd64 -v ./cmd/gpuffs
	@echo "Darwin amd64 cross compilation done:"
	@ls -ld $(GOBIN)/gpuffs-darwin-* | grep amd64

gpuffs-windows: geth-windows-386 geth-windows-amd64
	@echo "Windows cross compilation done:"
	@ls -ld $(GOBIN)/gpuffs-windows-*

gpuffs-windows-386:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=windows/386 -v ./cmd/gpuffs
	@echo "Windows 386 cross compilation done:"
	@ls -ld $(GOBIN)/gpuffs-windows-* | grep 386

gpuffs-windows-amd64:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=windows/amd64 -v ./cmd/gpuffs
	@echo "Windows amd64 cross compilation done:"
	@ls -ld $(GOBIN)/gpuffs-windows-* | grep amd64
