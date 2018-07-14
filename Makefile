.PHONY: all test clean dirty generate sameriver-generate

all: clean deps generate test clean

generate: sameriver-generate
	./sameriver-generate  -outputdir=./engine

test:
	go test -v -coverprofile=coverage.txt -race ./engine

sameriver-generate:
	go build -o sameriver-generate ./generator/main

dirty:
	cp /tmp/sameriver/* engine/

clean:
	mkdir /tmp/sameriver 2>/dev/null || true
	mv engine/CUSTOM_* /tmp/sameriver/ 2>/dev/null || true
	cp engine/GENERATED_* /tmp/sameriver
	git checkout HEAD -- engine/GENERATED_*
	rm sameriver-generate 2>/dev/null || true

deps:
	./install_deps.sh
	./install_godeps.sh 
