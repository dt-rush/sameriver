.PHONY: all test deps

all: deps test

test:
	cd v3 && go test -v -coverprofile=../coverage.txt .

deps:
	./install_deps.sh
	cd v3 && go mod tidy
