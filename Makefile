all: generate clean

generate: sameriver-generate
	./sameriver-generate  -outputdir=./engine

sameriver-generate:
	go build sameriver-generate.go

clean:
	rm engine/CUSTOM_* 2>/dev/null || true


