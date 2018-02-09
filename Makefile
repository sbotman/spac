## DEVELOPMENT COMMANDS
.PHONY: help clean full_clean check test dev build updatedeps cover

default: test

clean: ## Cleans the project
	go clean

full_clean: ## Cleans the project and all it's artifacts.
	@rm -f bin/*
	go clean

check: ## Checks codestyle and correctness
	@which gometalinter >/dev/null; if [ $$? -eq 1 ]; then \
		go get -v -u github.com/alecthomas/gometalinter; \
		gometalinter --install --update; \
	fi
	gometalinter --vendor --disable-all --enable=vet --enable=golint ./...

test: ## Tests the project
	go test -cover ./...

dev:
	@mkdir -p bin
	go build -o bin/spac

build:
	@mkdir -p bin
	GOOS=linux GOARCH=amd64 go build -o bin/spac
	GOOS=windows GOARCH=amd64 go build -o bin/spac.exe

updatedeps: ## Updates the vendored Go dependencies and global grpc/protobuf compiler
	dep ensure -update

cover: ## Shows test coverage.
	@go tool cover 2>/dev/null; if [ $$? -eq 3 ]; then \
		go get -u golang.org/x/tools/cmd/cover; \
	fi
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out
	rm coverage.out

## MAKE HELP
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

