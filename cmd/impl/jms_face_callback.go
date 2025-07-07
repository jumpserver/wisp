package impl

import (
	"context"

	"github.com/jumpserver-dev/sdk-go/service"
	pb "github.com/jumpserver/wisp/protobuf-go/protobuf"
)

func (j *JMServer) FaceRecognitionCallback(ctx context.Context, req *pb.FaceRecognitionCallbackRequest) (*pb.FaceRecognitionCallbackResponse, error) {
	var status pb.Status

	result := service.FaceRecognitionResult{
		Token:        req.Token,
		Success:      req.Success,
		ErrorMessage: req.ErrorMessage,
		FaceCode:     req.FaceCode,
	}
	if err := j.apiClient.SendFaceRecognitionCallback(result); err != nil {
		status.Ok = false
		status.Err = err.Error()
	} else {
		status.Ok = true
	}
	return &pb.FaceRecognitionCallbackResponse{
		Status: &status,
	}, nil
}

func (j *JMServer) FaceMonitorCallback(ctx context.Context, req *pb.FaceMonitorCallbackRequest) (*pb.FaceMonitorCallbackResponse, error) {
	var status pb.Status

	result := service.FaceMonitorResult{
		Token:        req.Token,
		Success:      req.Success,
		ErrorMessage: req.ErrorMessage,
		FaceCodes:    req.FaceCodes,
		IsFinished:   req.IsFinished,
		Action:       req.Action,
	}
	if err := j.apiClient.SendFaceMonitorCallback(result); err != nil {
		status.Ok = false
		status.Err = err.Error()
	} else {
		status.Ok = true
	}
	return &pb.FaceMonitorCallbackResponse{
		Status: &status,
	}, nil
}

func (j *JMServer) JoinFaceMonitor(ctx context.Context, req *pb.JoinFaceMonitorRequest) (*pb.JoinFaceMonitorResponse, error) {

	var status pb.Status

	joinRequest := service.JoinFaceMonitorRequest{
		FaceMonitorToken: req.FaceMonitorToken,
		SessionId:        req.SessionId,
	}
	if err := j.apiClient.JoinFaceMonitor(joinRequest); err != nil {
		status.Ok = false
		status.Err = err.Error()
	} else {
		status.Ok = true
	}
	return &pb.JoinFaceMonitorResponse{
		Status: &status,
	}, nil
}
