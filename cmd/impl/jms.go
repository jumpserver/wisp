package impl

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/jumpserver/wisp/pkg/jms-sdk-go/service"
	pb "github.com/jumpserver/wisp/protobuf-go/protobuf"
)

type JMServer struct {
	pb.UnimplementedServiceServer
	apiClient *service.JMService
}

func (j *JMServer) GetAPIClient() *service.JMService {
	return j.apiClient
}

func (j *JMServer) GetDBTokenAuthInfo(context.Context, *pb.DBTokenRequest) (*pb.DBTokenResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetDBTokenAuthInfo not implemented")
}

func (j *JMServer) CreateSession(context.Context, *pb.SessionCreateRequest) (*pb.SessionCreateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateSession not implemented")
}

func (j *JMServer) FinishSession(context.Context, *pb.SessionFinishRequest) (*pb.SessionFinishResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FinishSession not implemented")
}

func (j *JMServer) UploadActiveSessions(context.Context, *pb.ActiveSessRequest) (*pb.ActiveSessResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UploadActiveSessions not implemented")
}

func (j *JMServer) UploadReplayFile(context.Context, *pb.ReplayRequest) (*pb.ReplayResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UploadReplayFile not implemented")
}

func (j *JMServer) UploadCommand(context.Context, *pb.CommandRequest) (*pb.CommandResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UploadCommand not implemented")
}
