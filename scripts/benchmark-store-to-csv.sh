#!/usr/bin/env bash

cd "${0%/*}"
cd ..

go test -bench=. "$@" ./internal/store | awk '/^Benchmark/ {print $1","$3}'