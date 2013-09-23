.PHONY: test fmt

concentrate: fmt concentrate.go
	go build

test: fmt
	go test

fmt:
	go fmt
