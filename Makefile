proto_path=./protos/
proto_files=./protos/*.proto


.PHONY: proto
proto: proto-go proto-py

.PHONY: proto-go
proto-go:
	rm -rf  ./protobuf-go/protobuf/*
	protoc --proto_path=${proto_path} --go_out=./protobuf-go/  \
	--go-grpc_out=./protobuf-go/ \
	--go_opt=paths=import \
	${proto_files}

.PHONY: proto-py
proto-py:
	rm -rf ./protobuf-py/protobuf/*
	python -m grpc_tools.protoc --proto_path=${proto_path} --python_out=./protobuf-py/protobuf/ \
	--grpc_python_out=./protobuf-py/protobuf/ \
	${proto_files}

NAME=wisp
VERSION ?=Unknown
BuildTime:=$(shell date -u '+%Y-%m-%d %I:%M:%S%p')
# COMMIT:=$(shell git rev-parse HEAD)
GOVERSION:=$(shell go version)
GOLDFLAGS=-X 'github.com/jumpserver/wisp/cmd.version=$(VERSION)'
GOLDFLAGS+=-X 'github.com/jumpserver/wisp/cmd.BuildTime=$(BuildTime)'
GOLDFLAGS+=-X 'github.com/jumpserver/wisp/cmd.GitCommit=$(COMMIT)'
GOLDFLAGS+=-X 'github.com/jumpserver/wisp/cmd.GoVersion=$(GOVERSION)'

