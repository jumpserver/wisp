package client

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/jumpserver/wisp/protobuf-go/protobuf"
)

func NewClient(opts ...Options) (pb.ServiceClient, error) {
	opt := option{
		addr: "localhost:9090",
	}
	for _, setter := range opts {
		setter(&opt)
	}
	conn, err := grpc.Dial(
		opt.addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}
	client := pb.NewServiceClient(conn)
	return client, nil
}

type option struct {
	addr string
}

type Options func(*option)

func WithAddr(addr string) Options {
	return func(o *option) {
		o.addr = addr
	}
}
