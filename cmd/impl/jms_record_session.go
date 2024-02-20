package impl

import (
	"context"

	"github.com/jumpserver/wisp/pkg/jms-sdk-go/model"
	"github.com/jumpserver/wisp/pkg/logger"
	pb "github.com/jumpserver/wisp/protobuf-go/protobuf"
)

func (j *JMServer) RecordSessionLifecycleLog(ctx context.Context, req *pb.SessionLifecycleLogRequest) (*pb.StatusResponse, error) {
	var (
		status pb.Status
	)
	sessionId := req.GetSessionId()
	event := req.GetEvent()
	reason := req.GetReason()
	user := req.GetUser()
	logObj := model.SessionLifecycleLog{Reason: reason, User: user}
	lifecycleEvent := LifecycleEventMap[event]
	logger.Infof("Request record session %s lifecyle %s : %s", sessionId, event, logObj)
	if err := j.apiClient.RecordSessionLifecycleLog(sessionId, lifecycleEvent, logObj); err != nil {
		logger.Errorf("Record session %s lifecyle failed: %s", sessionId, err)
		status.Err = err.Error()
		return &pb.StatusResponse{Status: &status}, nil
	}
	status.Ok = true
	return &pb.StatusResponse{Status: &status}, nil
}
