all: generate clean

generate: sameriver-generate
	./sameriver-generate  -outputdir=./engine

sameriver-generate:
	go build sameriver-generate.go

dirty:
	cp /tmp/sameriver/* engine/

clean:
	mkdir /tmp/sameriver 2>/dev/null || true
	mv engine/CUSTOM_* /tmp/sameriver/ 2>/dev/null || true
	cp engine/GENERATED_* /tmp/sameriver
	git checkout HEAD -- engine/GENERATED_*
	rm sameriver-generate


