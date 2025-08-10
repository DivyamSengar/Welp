#!/bin/bash
# A simple script to run the benchmarks
gcc -O2 -pthread -o thread_mini thread_mini.c

# Check if the first argument is "--perf"
if [ "$1" == "--perf" ]; then
    shift
    echo "Running under perf stat for enhanced metrics..."
    # Run perf with the desired events. Should I change these?
    perf stat -e cycles,cache-references,cache-misses ./thread_mini "$@"
else
    # Run the benchmark normally.
    ./thread_mini "$@"
fi
