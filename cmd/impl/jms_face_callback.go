package impl

import (
	"context"
	"github.com/jumpserver/wisp/pkg/jms-sdk-go/service"
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
