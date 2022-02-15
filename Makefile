proto_path=./protos/
proto_files=./protos/*.proto

protobuf_go_dir=./protobuf-go/protobuf
protobuf_py_dir=./protobuf-py/protobuf

.PHONY: proto
proto: proto-go proto-py

.PHONY: proto-go
proto-go:
	@mkdir -p ${protobuf_go_dir}
	rm -rf  ${protobuf_go_dir}/*
	protoc --proto_path=${proto_path} --go_out=./protobuf-go/  \
	--go-grpc_out=./protobuf-go/ \
	--go_opt=paths=import \
	${proto_files}

.PHONY: proto-py
proto-py:
	@mkdir -p ${protobuf_py_dir}
	rm -rf ${protobuf_py_dir}/*
	python -m grpc_tools.protoc --proto_path=${proto_path} --python_out=${protobuf_py_dir} \
	--grpc_python_out=${protobuf_py_dir} \
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

