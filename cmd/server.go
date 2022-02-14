package cmd

import (
	"context"
	"net"

	"google.golang.org/grpc"

	"github.com/jumpserver/wisp/pkg/logger"
	pb "github.com/jumpserver/wisp/protobuf-go/protobuf"
)

func NewServer(addr string, imp pb.ServiceServer) *Server {
	grpcSrv := grpc.NewServer()
	pb.RegisterServiceServer(grpcSrv, imp)
	return &Server{
		addr:       addr,
		impService: imp,
		grpcSrv:    grpcSrv,
	}
}

type Server struct {
	addr       string
	impService pb.ServiceServer
	grpcSrv    *grpc.Server

	ln net.Listener
}

func (s *Server) Run() (err error) {
	s.ln, err = net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}
	logger.Infof("Listen on: %s", s.addr)
	err = s.grpcSrv.Serve(s.ln)
	return err
}

func (s *Server) Stop() {
	if s.grpcSrv != nil {
		s.grpcSrv.GracefulStop()
	}
	if s.ln != nil {
		_ = s.ln.Close()
	}
	logger.Info("Server stop")
}

func Run(ctx context.Context, srv *Server) {
	go func() {
		<-ctx.Done()
		srv.Stop()
	}()
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
