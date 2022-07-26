#!/bin/bash
# shellcheck disable=SC2046
# shellcheck disable=SC2006
protoc common.proto --proto_path=./ --go_out=../
# shellcheck disable=SC2046
# shellcheck disable=SC2006
goctl rpc protoc -I=. hello.proto -v --go_out=.. --go-grpc_out=..  --zrpc_out=.. --style=goZero
# shellcheck disable=SC2046
protoc --proto_path=. --descriptor_set_out=./hello.pb hello.proto
# shellcheck disable=SC2046
protoc --proto_path=. --descriptor_set_out=./common.pb common.proto