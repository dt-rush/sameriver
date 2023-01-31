.PHONY: all test deps

all: deps test

test:
	go test -v -coverprofile=coverage.txt ./engine

deps:
	./install_deps.sh
	go mod tidy
