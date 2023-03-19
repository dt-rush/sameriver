.PHONY: all test deps

all: deps test

test:
	cd v4 && go test -v -coverprofile=../coverage.txt .

deps:
	./install_deps.sh
	cd v4 && go mod tidy
