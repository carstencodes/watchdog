BINARY_NAME=watchdog

_LD_DETAILS_MODULE=github.com/carstencodes/watchdog/internal/lib/common
_LD_FLAG_COPYRIGHT_YEAR=-X ${_LD_DETAILS_MODULE}.yearNow=`date -u +%Y`
_LD_FLAG_VERSION=-X ${_LD_DETAILS_MODULE}.version=`git describe --abbrev=0 --tag --dirty=-local`
LINKER_FLAGS=-ldflags "${_LD_FLAG_COPYRIGHT_YEAR} ${_LD_FLAG_VERSION}"

ifeq ($(OS),Windows_NT)
    _RUN_OS := windows
    _RUN_EXT := .exe
else
    _RUN_OS := $(shell uname | tr "[:upper:]" "[:lower:]")
endif

DEFAULT: lint vet build

build:
	GOARCH=amd64 GOOS=darwin go build -o ./bin/${BINARY_NAME}-darwin ${LINKER_FLAGS} cmd/watchdog/main.go
	GOARCH=amd64 GOOS=linux go build -o ./bin/${BINARY_NAME}-linux ${LINKER_FLAGS} cmd/watchdog/main.go
	GOARCH=amd64 GOOS=windows go build -o ./bin/${BINARY_NAME}-windows.exe ${LINKER_FLAGS} cmd/watchdog/main.go

run: build
	./bin/${BINARY_NAME}-${_RUN_OS}${_RUN_EXT}

clean:
	go clean
	rm ./bin/${BINARY_NAME}-darwin
	rm ./bin/${BINARY_NAME}-linux
	rm ./bin/${BINARY_NAME}-windows.exe

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