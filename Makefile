NAME=wisp
BUILDDIR=build
VERSION ?=Unknown
BuildTime:=$(shell date -u '+%Y-%m-%d %I:%M:%S%p')
COMMIT:=$(shell git rev-parse HEAD)
GOVERSION:=$(shell go version)

GOOS:=$(shell go env GOOS)
GOARCH:=$(shell go env GOARCH)

LDFLAGS=-w -s

GOLDFLAGS=-X 'github.com/jumpserver/wisp/cmd.Version=$(VERSION)'
GOLDFLAGS+=-X 'github.com/jumpserver/wisp/cmd.BuildTime=$(BuildTime)'
GOLDFLAGS+=-X 'github.com/jumpserver/wisp/cmd.GitCommit=$(COMMIT)'
GOLDFLAGS+=-X 'github.com/jumpserver/wisp/cmd.GoVersion=$(GOVERSION)'
WISPBUILD=CGO_ENABLED=0 go build -trimpath -ldflags "$(GOLDFLAGS) ${LDFLAGS}"

define make_artifact_full
	GOOS=$(1) GOARCH=$(2) $(WISPBUILD) -o $(BUILDDIR)/$(NAME)-$(1)-$(2)
	mkdir -p $(BUILDDIR)/$(NAME)-$(VERSION)-$(1)-$(2)
	cp $(BUILDDIR)/$(NAME)-$(1)-$(2) $(BUILDDIR)/$(NAME)-$(VERSION)-$(1)-$(2)/$(NAME)
	cd $(BUILDDIR) && tar -czvf $(NAME)-$(VERSION)-$(1)-$(2).tar.gz $(NAME)-$(VERSION)-$(1)-$(2)
	rm -rf $(BUILDDIR)/$(NAME)-$(VERSION)-$(1)-$(2) $(BUILDDIR)/$(NAME)-$(1)-$(2)
endef

build:
	GOARCH=$(GOARCH) GOOS=$(GOOS) $(WISPBUILD) -o $(BUILDDIR)/$(NAME) .

all:
	$(call make_artifact_full,darwin,amd64)
	$(call make_artifact_full,darwin,arm64)
	$(call make_artifact_full,linux,amd64)
	$(call make_artifact_full,linux,arm64)
	$(call make_artifact_full,linux,mips64le)
	$(call make_artifact_full,linux,ppc64le)
	$(call make_artifact_full,linux,s390x)
	$(call make_artifact_full,linux,riscv64)

local:
	$(call make_artifact_full,$(shell go env GOOS),$(shell go env GOARCH))

darwin-amd64:
	$(call make_artifact_full,darwin,amd64)

darwin-arm64:
	$(call make_artifact_full,darwin,arm64)

linux-amd64:
	$(call make_artifact_full,linux,amd64)

linux-arm64:
	$(call make_artifact_full,linux,arm64)

linux-loong64:
	$(call make_artifact_full,linux,loong64)

linux-mips64le:
	$(call make_artifact_full,linux,mips64le)

linux-ppc64le:
	$(call make_artifact_full,linux,ppc64le)

linux-s390x:
	$(call make_artifact_full,linux,s390x)

linux-riscv64:
	$(call make_artifact_full,linux,riscv64)

clean:
	rm -rf $(BUILDDIR)

proto_path=./protos/
proto_files=./protos/*.proto

proto_go_dir=./protobuf-go
protobuf_py_dir=./protobuf-py/protobuf

.PHONY: proto
proto: proto-go

.PHONY: proto-go
proto-go:
	@mkdir -p ${proto_go_dir}/protobuf
	protoc --proto_path=${proto_path} --go_out=${proto_go_dir}  \
	--go-grpc_out=${proto_go_dir} \
	--go_opt=paths=import \
	${proto_files}



.PHONY: proto-py
proto-py:
	@mkdir -p ./protobuf-py/protobuf
	python -m grpc_tools.protoc --proto_path=${proto_path} --python_out=./protobuf-py/protobuf \
	--pyi_out=./protobuf-py/protobuf \
	--grpc_python_out=./protobuf-py/protobuf \
	${proto_files}
