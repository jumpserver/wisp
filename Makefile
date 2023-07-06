NAME=wisp
BUILDDIR=build
VERSION ?=Unknown
BuildTime:=$(shell date -u '+%Y-%m-%d %I:%M:%S%p')
COMMIT:=$(shell git rev-parse HEAD)
GOVERSION:=$(shell go version)

LDFLAGS=-w -s

GOLDFLAGS=-X 'github.com/jumpserver/wisp/cmd.Version=$(VERSION)'
GOLDFLAGS+=-X 'github.com/jumpserver/wisp/cmd.BuildTime=$(BuildTime)'
GOLDFLAGS+=-X 'github.com/jumpserver/wisp/cmd.GitCommit=$(COMMIT)'
GOLDFLAGS+=-X 'github.com/jumpserver/wisp/cmd.GoVersion=$(GOVERSION)'
WISPBUILD=CGO_ENABLED=0 go build -trimpath -ldflags "$(GOLDFLAGS) ${LDFLAGS}"

PLATFORM_LIST = \
	darwin-amd64 \
	darwin-arm64 \
	linux-amd64 \
	linux-arm64

WINDOWS_ARCH_LIST = \
	windows-amd64

all-arch: $(PLATFORM_LIST) $(WINDOWS_ARCH_LIST)

darwin-amd64:
	GOARCH=amd64 GOOS=darwin $(WISPBUILD) -o $(BUILDDIR)/$(NAME)-$@
	mkdir -p $(BUILDDIR)/$(NAME)-$(VERSION)-$@
	cp $(BUILDDIR)/$(NAME)-$@ $(BUILDDIR)/$(NAME)-$(VERSION)-$@/$(NAME)
	cd $(BUILDDIR) && tar -czvf $(NAME)-$(VERSION)-$@.tar.gz $(NAME)-$(VERSION)-$@
	rm -rf $(BUILDDIR)/$(NAME)-$(VERSION)-$@ $(BUILDDIR)/$(NAME)-$@

darwin-arm64:
	GOARCH=arm64 GOOS=darwin $(WISPBUILD) -o $(BUILDDIR)/$(NAME)-$@
	mkdir -p $(BUILDDIR)/$(NAME)-$(VERSION)-$@
	cp $(BUILDDIR)/$(NAME)-$@ $(BUILDDIR)/$(NAME)-$(VERSION)-$@/$(NAME)
	cd $(BUILDDIR) && tar -czvf $(NAME)-$(VERSION)-$@.tar.gz $(NAME)-$(VERSION)-$@
	rm -rf $(BUILDDIR)/$(NAME)-$(VERSION)-$@ $(BUILDDIR)/$(NAME)-$@

linux-amd64:
	GOARCH=amd64 GOOS=linux $(WISPBUILD) -o $(BUILDDIR)/$(NAME)-$@
	mkdir -p $(BUILDDIR)/$(NAME)-$(VERSION)-$@
	cp $(BUILDDIR)/$(NAME)-$@ $(BUILDDIR)/$(NAME)-$(VERSION)-$@/$(NAME)
	cd $(BUILDDIR) && tar -czvf $(NAME)-$(VERSION)-$@.tar.gz $(NAME)-$(VERSION)-$@
	rm -rf $(BUILDDIR)/$(NAME)-$(VERSION)-$@ $(BUILDDIR)/$(NAME)-$@

linux-arm64:
	GOARCH=arm64 GOOS=linux $(WISPBUILD) -o $(BUILDDIR)/$(NAME)-$@
	mkdir -p $(BUILDDIR)/$(NAME)-$(VERSION)-$@
	cp $(BUILDDIR)/$(NAME)-$@ $(BUILDDIR)/$(NAME)-$(VERSION)-$@/$(NAME)
	cd $(BUILDDIR) && tar -czvf $(NAME)-$(VERSION)-$@.tar.gz $(NAME)-$(VERSION)-$@
	rm -rf $(BUILDDIR)/$(NAME)-$(VERSION)-$@ $(BUILDDIR)/$(NAME)-$@

linux-loong64:
	GOARCH=loong64 GOOS=linux $(WISPBUILD) -o $(BUILDDIR)/$(NAME)-$@
	mkdir -p $(BUILDDIR)/$(NAME)-$(VERSION)-$@
	cp $(BUILDDIR)/$(NAME)-$@ $(BUILDDIR)/$(NAME)-$(VERSION)-$@/$(NAME)
	cd $(BUILDDIR) && tar -czvf $(NAME)-$(VERSION)-$@.tar.gz $(NAME)-$(VERSION)-$@
	rm -rf $(BUILDDIR)/$(NAME)-$(VERSION)-$@ $(BUILDDIR)/$(NAME)-$@

windows-amd64:
	GOARCH=amd64 GOOS=windows $(WISPBUILD) -o $(BUILDDIR)/$(NAME)-$@.exe
	mkdir -p $(BUILDDIR)/$(NAME)-$(VERSION)-$@
	cp $(BUILDDIR)/$(NAME)-$@.exe $(BUILDDIR)/$(NAME)-$(VERSION)-$@/$(NAME).exe
	cd $(BUILDDIR) && tar -czvf $(NAME)-$(VERSION)-$@.tar.gz $(NAME)-$(VERSION)-$@
	rm -rf $(BUILDDIR)/$(NAME)-$(VERSION)-$@ $(BUILDDIR)/$(NAME)-$@.exe


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
