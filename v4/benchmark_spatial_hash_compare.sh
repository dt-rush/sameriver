#!/bin/bash

#    single better for 10
#    parallel better for 12
#    single better for 14
#    single better for 16
#    parallel better for 18
#    single better for 20
#    parallel better for 22
#    parallel better for 24
#    parallel better for 26
#    parallel better for 28
#    parallel better for 30
#    parallel better for 32
#    parallel better for 34
#    parallel better for 36
#    parallel better for 38
#    parallel better for 40
#    single better for 38
#    single better for 40
#    parallel better for 42
#    single better for 44
#    parallel better for 46
#    single better for 48
#    parallel better for 50
#    single better for 52
#    parallel better for 54
#    parallel better for 56
#    single better for 58
#    parallel better for 60
#    parallel better for 62
#    parallel better for 64

for x in {10..100..2}; do
  parallel_c_output=$(GRIDX=$x GRIDY=$x go test -v -run=BenchmarkSpatialHashUpdateParallelC -bench=BenchmarkSpatialHashUpdateParallelC . | tail -n 3 | head -n 1)
  single_output=$(GRIDX=$x GRIDY=$x go test -v -run=BenchmarkSpatialHashUpdateSingle -bench=BenchmarkSpatialHashUpdateSingle . | tail -n 3 | head -n 1)

  parallel_c_ns=$(echo $parallel_c_output | awk '{print $(NF-1)}')
  single_ns=$(echo $single_output | awk '{print $(NF-1)}')

  if [[ $parallel_c_ns -gt $single_ns ]]; then
    echo "parallel better for $x"
  else
    echo "single better for $x"
  fi

  sleep 10 # prevent CPU overheat

done
