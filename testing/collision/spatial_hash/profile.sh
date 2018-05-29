#!/bin/bash

file="test.prof"
go build -o test.out spatial_hash.go
./test.out -cpuprofile=${file}
echo ""
echo "profiling done. result can be viewed with go tool pprof ${file}"
