## VARIABLES
FOLDERS=./. 

.PHONY: help

default: test

## DEVELOPMENT COMMANDS
clean: ## Cleans the project
	@go clean

full_clean: clean ## Cleans the project (including the compiled proto and migrations)
	@rm -f **/*.pb.go

check: ## Checks codestyle and correctness
	@which gometalinter >/dev/null; if [ $$? -eq 1 ]; then \
		go get -v -u github.com/alecthomas/gometalinter; \
		gometalinter --install --update; \
	fi
	gometalinter --vendor --disable-all --enable=vet --enable=golint $(FOLDERS)

dev:
	@mkdir -p bin
	@GOOS=linux GOARCH=amd64 go build -o bin/spac
	@GOOS=windows GOARCH=amd64 go build -o bin/spac.exe

test: ## Tests the project
	@go test -cover $(FOLDERS)

updatedeps: ## Updates the vendored Go dependencies and global grpc/protobuf compiler
	@go get -u github.com/golang/protobuf/{proto,protoc-gen-go}
	@govendor fetch +external +vendor

cover: ## Shows coverage
	@go tool cover 2>/dev/null; if [ $$? -eq 3 ]; then \
		go get -u golang.org/x/tools/cmd/cover; \
	fi
	go test $(TEST) -coverprofile=coverage.out
	go tool cover -html=coverage.out
	rm coverage.out

## MAKE HELP
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

