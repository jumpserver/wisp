package cmd

import (
	"context"
	"net"
	"sync"

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
		done:       make(chan struct{}),
	}
}

type Server struct {
	addr       string
	impService pb.ServiceServer
	grpcSrv    *grpc.Server

	ln net.Listener

	once sync.Once
	done chan struct{}
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
	select {
	case <-s.done:
		logger.Info("Server already stop")
	default:
		s.once.Do(func() {
			if s.grpcSrv != nil {
				s.grpcSrv.GracefulStop()
			}
			if s.ln != nil {
				_ = s.ln.Close()
			}
			close(s.done)
		})
		logger.Info("Server stop")
	}
}

func (s *Server) Wait() {
	<-s.done
}

func RunServer(ctx context.Context, srv *Server) {
	go func() {
		<-ctx.Done()
		srv.Stop()
	}()
	if err := srv.Run(); err != nil {
		logger.Fatalf("Server err: %s", err)
	}
}
