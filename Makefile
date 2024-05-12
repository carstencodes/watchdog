BINARY_NAME=watchdog

PKG=github.com/carstencodes/watchdog
_LD_DETAILS_MODULE=${PKG}/internal/lib/common
_LD_FLAG_VALUE_COPYRIGHT_YEAR=$(shell date -u +%Y)
_LD_FLAG_COPYRIGHT_YEAR=-X ${_LD_DETAILS_MODULE}.yearNow=${_LD_FLAG_VALUE_COPYRIGHT_YEAR}
_LD_FLAG_VALUE_VERSION=$(shell git describe --abbrev=0 --tag --dirty=-local)
_LD_FLAG_VERSION=-X ${_LD_DETAILS_MODULE}.version=${_LD_FLAG_VALUE_VERSION}
LINKER_FLAGS=-ldflags "${_LD_FLAG_COPYRIGHT_YEAR} ${_LD_FLAG_VERSION}"

ifeq ($(OS),Windows_NT)
    _RUN_OS := windows
    _RUN_EXT := .exe
    _RUN_ARCH := ${PROCESSOR_ARCHITECTURE}
else
    _RUN_OS := $(shell uname | tr "[:upper:]" "[:lower:]")
	_RUN_ARCH := $(shell uname -m)
	ifeq (${_RUN_ARCH},x86_64)
		_RUN_ARCH = amd64
	else ifeq(${_RUN_ARCH},i686)
		_RUN_ARCH=386
	else ifeq(${_RUN_ARCH},aarch64)
		_RUN_ARCH=arm64
	endif
endif

DEFAULT: lint vet build

build: build-windows build-linux build-darwin

build-darwin: build-darwin-amd64 build-darwin-arm64

build-darwin-amd64:
	GOARCH=amd64 GOOS=darwin go build -o ./bin/${BINARY_NAME}-darwin-amd64 ${LINKER_FLAGS} cmd/watchdog/main.go

build-darwin-arm64:
	GOARCH=arm64 GOOS=darwin go build -o ./bin/${BINARY_NAME}-darwin-arm64 ${LINKER_FLAGS} cmd/watchdog/main.go

build-linux: build-linux-386 build-linux-amd64 build-linux-arm64

build-linux-386:
	GOARCH=amd64 GOOS=linux go build -o ./bin/${BINARY_NAME}-linux-386 ${LINKER_FLAGS} cmd/watchdog/main.go

build-linux-amd64:
	GOARCH=amd64 GOOS=linux go build -o ./bin/${BINARY_NAME}-linux-amd64 ${LINKER_FLAGS} cmd/watchdog/main.go

build-linux-arm64:
	GOARCH=amd64 GOOS=linux go build -o ./bin/${BINARY_NAME}-linux-arm64 ${LINKER_FLAGS} cmd/watchdog/main.go

build-windows: build-windows-amd64

build-windows-amd64:
	GOARCH=amd64 GOOS=windows go build -o ./bin/${BINARY_NAME}-windows-amd64.exe ${LINKER_FLAGS} cmd/watchdog/main.go

run: build
	./bin/${BINARY_NAME}-${_RUN_OS}-${_RUN_ARCH}${_RUN_EXT}

clean:
	go clean
	rm -rf ./bin/**

test:
	go test ./...

test_coverage:
	go test ./... -coverprofile=coverage.out

dep:
	go mod download

vet:
	go vet -json --all ./...

lint:
	golangci-lint run --enable-all --out-format json --sort-results | jq .Issues