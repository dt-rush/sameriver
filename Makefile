.PHONY: all test clean dirty generate 

all: clean generate test clean

generate: sameriver-generate
	./sameriver-generate  -outputdir=./engine

test:
	go test ./engine

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

install:
	./install_deps.sh 
