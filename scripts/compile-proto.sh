#!/usr/bin/env bash

cd "${0%/*}"
cd ..

protoc --go_out=. --go_opt=paths=import --go_opt=module=github.com/Matt-Kelly-/go-memory-cache \
    --go-grpc_out=. --go-grpc_opt=paths=import --go-grpc_opt=module=github.com/Matt-Kelly-/go-memory-cache \
    api/service.proto