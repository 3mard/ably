
build-client:
	go build -o bin/client cmd/client/main.go

build-server:
	go build -o bin/server cmd/server/main.go

tests:
	go test -v ./...
