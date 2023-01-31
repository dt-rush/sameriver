.PHONY: all test deps

all: deps test

test:
	cd v2 && go test -v -coverprofile=../coverage.txt .

deps:
	./install_deps.sh
	cd v2 && go mod tidy
