.PHONY: test fmt run release clean

concentrate: fmt concentrate.go
	go build

test: fmt
	go test

fmt:
	go fmt

release: concentrate-osx.zip concentrate-linux.tar.gz

concentrate-osx.zip:
	rm -f concentrate
	GOOS=darwin GOARCH=amd64 go build
	zip -q concentrate-osx.zip concentrate

concentrate-linux.tar.gz:
	rm -f concentrate
	GOOS=linux GOARCH=amd64 go build
	tar czf concentrate-linux.tar.gz concentrate

clean:
	rm -f concentrate
	rm -f concentrate-osx.zip
	rm -f concentrate-linux.tar.gz
