.PHONY: all test deps

all: deps test

test:
	go test -v -coverprofile=coverage.txt -race ./engine

deps:
	./install_deps.sh
	./install_godeps.sh 
