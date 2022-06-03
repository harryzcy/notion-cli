.PHONY: lint test

lint:
	golangci-lint run -v ./...

test: lint
	go test -v ./...
