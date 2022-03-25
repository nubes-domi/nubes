#!/bin/bash

# Rebuild the JS Client for sum gRPC
protoc \
    --plugin=./src/experior/node_modules/.bin/protoc-gen-ts_proto \
    --ts_proto_out=./src/experior/grpc \
    --ts_proto_opt=esModuleInterop=true \
    ./src/sum/sum.proto

# Rebuild the gRPC Server for sum
protoc \
    --go_out=./src/sum/rpc \
    --go_opt=paths=source_relative \
    --go-grpc_out=./src/sum/rpc \
    --go-grpc_opt=paths=source_relative \
    ./src/sum/sum.proto