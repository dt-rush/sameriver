#!/bin/bash

outputfile="test.prof"
go build -o test.out *.go
./test.out -cpuprofile=${outputfile} "$@"
echo ""
echo "profiling done. result can be viewed with go tool pprof ${outputfile}"
