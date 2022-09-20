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
	status.Ok = true
	ports, err := j.apiClient.GetListenPorts()
	if err != nil {
		status.Ok = false
		status.Err = err.Error()
		return &pb.ListenPortResponse{
			Status: &status,
		}, nil
	}

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
	app, err := j.apiClient.GetApplicationByPort(req.Port)
	if err != nil {
		status.Ok = false
		status.Err = err.Error()
		return &pb.PortInfoResponse{
			Status: &status,
		}, nil
	}
	pbApp := ConvertToProtobufApplication(app)
	if app.Domain != "" {
		domain, err := j.apiClient.GetDomainGateways(app.Domain)
		if err != nil {
			status.Ok = false
			status.Err = err.Error()
			return &pb.PortInfoResponse{
				Status: &status,
			}, nil
		}
		gateways = ConvertToProtobufGateways(domain.Gateways)
	}
	status.Ok = true
	info := pb.PortInfo{
		Application: pbApp,
		Gateways:    gateways,
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
