BINARY_NAME=watchdog

build:
	 GOARCH=amd64 GOOS=darwin go build -o ./bin/${BINARY_NAME}-darwin cmd/watchdog/main.go
	 GOARCH=amd64 GOOS=linux go build -o ./bin/${BINARY_NAME}-linux cmd/watchdog/main.go
	 GOARCH=amd64 GOOS=windows go build -o ./bin/${BINARY_NAME}-windows cmd/watchdog/main.go

run: build
	./bin/${BINARY_NAME}

clean:
	go clean
	rm ${BINARY_NAME}-darwin
	rm ${BINARY_NAME}-linux
	rm ${BINARY_NAME}-windows

test:
	go test ./...

test_coverage:
	go test ./... -coverprofile=coverage.out

dep:
	go mod download

vet:
	go vet

lint:
	golangci-lint run --enable-all