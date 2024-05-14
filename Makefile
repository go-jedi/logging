LOCAL_BIN:=$(CURDIR)/bin

install-golangci-lint:
	GOBIN=$(LOCAL_BIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

update-packages:
	go get -u ./...

lint:
	GOBIN=$(LOCAL_BIN) golangci-lint run ./... --config .golangci.yaml