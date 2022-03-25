#!/bin/bash

# Rebuild the JS Client for sum gRPC
grpc_tools_ruby_protoc \
    --ruby_out=./src/experior/grpc \
    --grpc_out=./src/experior/grpc \
    --proto_path=`pwd`/src/sum \
    `pwd`/src/sum/sum.proto

# Rebuild the gRPC Server for sum
protoc \
    --go_out=./src/sum/rpc \
    --go_opt=paths=source_relative \
    --go-grpc_out=./src/sum/rpc \
    --go-grpc_opt=paths=source_relative \
    --proto_path=`pwd`/src/sum \
    `pwd`/src/sum/sum.proto