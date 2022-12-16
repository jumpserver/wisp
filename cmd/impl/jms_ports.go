package impl

import (
	"context"

	"github.com/jumpserver/wisp/pkg/logger"
	pb "github.com/jumpserver/wisp/protobuf-go/protobuf"
)

func (j *JMServer) GetListenPorts(ctx context.Context, req *pb.Empty) (*pb.ListenPortResponse, error) {
	var (
		status pb.Status
	)
	ports, err := j.apiClient.GetListenPorts()
	if err != nil {
		status.Err = err.Error()
		return &pb.ListenPortResponse{
			Status: &status,
		}, nil
	}
	status.Ok = true
	return &pb.ListenPortResponse{
		Status: &status,
		Ports:  ports,
	}, nil
}

func (j *JMServer) GetPortInfo(ctx context.Context,
	req *pb.PortInfoRequest) (*pb.PortInfoResponse, error) {
	var (
		status   pb.Status
		gateways []*pb.Gateway
	)
	// todo: 网域网关的情况
	app, err := j.apiClient.GetAssetByPort(req.Port)
	if err != nil {
		status.Ok = false
		status.Err = err.Error()
		return &pb.PortInfoResponse{
			Status: &status,
		}, nil
	}
	pbAsset := ConvertToProtobufAsset(app)
	status.Ok = true
	info := pb.PortInfo{
		Asset:    pbAsset,
		Gateways: gateways,
	}
	return &pb.PortInfoResponse{Status: &status, Data: &info}, nil
}

func (j *JMServer) HandlePortFailure(ctx context.Context,
	req *pb.PortFailureRequest) (*pb.StatusResponse, error) {
	var (
		status pb.Status
	)
	for _, portInfo := range req.GetData() {
		logger.Errorf("Listen port %d failure, reason: %s", portInfo.Port, portInfo.Reason)
	}
	status.Ok = true
	return &pb.StatusResponse{Status: &status}, nil
}
