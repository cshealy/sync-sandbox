#!/bin/sh

echo "--- Generating gRPC stub ---";
/usr/bin/protoc -I/protobuf -I. --go_out=plugins=grpc:. ./sync.proto;

echo "--- Generating reverse-proxy ---";
/usr/bin/protoc -I/protobuf -I. --grpc-gateway_out=logtostderr=true:. ./sync.proto;

echo "--- Complete ---"