.PHONY: build

clean:
	@go clean
	@go clean -testcache

test:
	@go test ./...

build:
	@go build -o /dev/null

install: test
	@go install
