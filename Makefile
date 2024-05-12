BINARY_NAME=watchdog

PKG=github.com/carstencodes/watchdog
_VERSION=$(shell git describe --abbrev=0 --tag --dirty=-local)
_CURRENT_YEAR ?= $(shell date -u +%Y)
_LD_DETAILS_MODULE=${PKG}/internal/lib/common
_LD_FLAG_VALUE_COPYRIGHT_YEAR=${_CURRENT_YEAR}
_LD_FLAG_COPYRIGHT_YEAR=-X ${_LD_DETAILS_MODULE}.yearNow=${_LD_FLAG_VALUE_COPYRIGHT_YEAR}
_LD_FLAG_VALUE_VERSION=${_VERSION}
_LD_FLAG_VERSION=-X ${_LD_DETAILS_MODULE}.version=${_LD_FLAG_VALUE_VERSION}
LINKER_FLAGS=-ldflags "${_LD_FLAG_COPYRIGHT_YEAR} ${_LD_FLAG_VERSION}"

ifeq ($(OS),Windows_NT)
    _DIST_OS := windows
    _DIST_EXT := .exe
    _DIST_ARCH := ${PROCESSOR_ARCHITECTURE}
else
    _DIST_OS := $(shell uname | tr "[:upper:]" "[:lower:]")
	_DIST_ARCH := $(shell uname -m)
	ifeq (${_DIST_ARCH},x86_64)
		_DIST_ARCH = amd64
	else ifeq (${_DIST_ARCH},i686)
		_DIST_ARCH=386
	else ifeq (${_DIST_ARCH},aarch64)
		_DIST_ARCH=arm64
	endif
endif

_DOCKER_GID ?= $(shell getent group docker | cut -d: -f3)

DEFAULT: lint vet build vuln

build: build/windows build/linux build/darwin

build/darwin: build/darwin/amd64 build/darwin/arm64

build/darwin/amd64:
	CGO_ENABLED=0 GOARCH=amd64 GOOS=darwin go build -o ./bin/${BINARY_NAME}-darwin-amd64 ${LINKER_FLAGS} cmd/watchdog/main.go

build/darwin/arm64:
	CGO_ENABLED=0 GOARCH=arm64 GOOS=darwin go build -o ./bin/${BINARY_NAME}-darwin-arm64 ${LINKER_FLAGS} cmd/watchdog/main.go

build/linux: build/linux/386 build/linux/amd64 build/linux/arm64

build/linux/386:
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -o ./bin/${BINARY_NAME}-linux-386 ${LINKER_FLAGS} cmd/watchdog/main.go

build/linux/amd64:
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -o ./bin/${BINARY_NAME}-linux-amd64 ${LINKER_FLAGS} cmd/watchdog/main.go

build/linux/arm64:
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -o ./bin/${BINARY_NAME}-linux-arm64 ${LINKER_FLAGS} cmd/watchdog/main.go

build/windows: build/windows/amd64

build/windows/amd64:
	CGO_ENABLED=0 GOARCH=amd64 GOOS=windows go build -o ./bin/${BINARY_NAME}-windows-amd64.exe ${LINKER_FLAGS} cmd/watchdog/main.go

run:
	go run -ldflags "${_LD_FLAG_COPYRIGHT_YEAR} -X ${_LD_DETAILS_MODULE}.version=${_LD_FLAG_VALUE_VERSION}-dev" ./cmd/watchdog/main.go

dist: build
	mkdir -p ./dist
	cp ./bin/${BINARY_NAME}-${_DIST_OS}-${_DIST_ARCH}${_DIST_EXT} ./dist/${BINARY_NAME}${_DIST_EXT}

docker/build:
	docker build -f Dockerfile --target watchdog --tag carstencodes/watchdog:${_VERSION} .

docker/run: docker/build
	docker run --rm --name watchdog-local -u 1000:${_DOCKER_GID} -v /var/run/docker.sock:/var/run/docker.sock carstencodes/watchdog:${_VERSION}

docker/dev: dist
	docker run --rm --name watchdog-dev -u 1000:${_DOCKER_GID} -v /var/run/docker.sock:/var/run/docker.sock $(shell docker build -q --target dev .)

docker/clean:
	docker buildx prune

clean: docker/clean
	go clean
	rm -rf ./bin/**

test:
	go test ./...

test_coverage:
	go test ./... -coverprofile=coverage.out

dep:
	go mod download

update:
	go get -u

vuln:
	govulncheck ./...

vet:
	go vet -json --all ./...

lint:
	golangci-lint run --enable-all --out-format json --sort-results | jq .Issues
