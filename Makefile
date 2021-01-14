.PHONY: build test coverage benchmark

build:
	go build

test:
	go test ./... -v -cover

coverage:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out

benchmark:
	go test ./... -bench=.
