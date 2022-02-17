package impl

import (
	"context"
	"fmt"
	"io"

	"github.com/jumpserver/wisp/cmd/common"
	"github.com/jumpserver/wisp/pkg/jms-sdk-go/model"
	"github.com/jumpserver/wisp/pkg/jms-sdk-go/service"
	"github.com/jumpserver/wisp/pkg/logger"
	pb "github.com/jumpserver/wisp/protobuf-go/protobuf"
)

func NewJMServer(apiClient *service.JMService, uploader *common.UploaderService,
	beat *common.BeatService) *JMServer {
	return &JMServer{
		apiClient: apiClient,
		uploader:  uploader,
		beat:      beat,
	}
}

type JMServer struct {
	pb.UnimplementedServiceServer
	apiClient *service.JMService

	uploader *common.UploaderService
	beat     *common.BeatService
}

func (j *JMServer) GetDBTokenAuthInfo(ctx context.Context, req *pb.DBTokenRequest) (*pb.DBTokenResponse, error) {
	var status pb.Status
	tokenAuthInfo, err := j.apiClient.GetConnectTokenAuth(req.Token)
	if err != nil {
		status.Err = err.Error()
		logger.Error("Get Connect token auth failed: %s", err)
		return &pb.DBTokenResponse{Status: &status}, nil
	}
	if tokenAuthInfo.TypeName != model.ConnectApplication {
		msg := fmt.Sprintf("Bad request: token %s connect type not %s", req.Token,
			model.ConnectApplication)
		status.Err = msg
		logger.Error(msg)
		return &pb.DBTokenResponse{Status: &status}, nil
	}
	dbTokenInfo := pb.DBTokenAuthInfo{
		KeyId:       tokenAuthInfo.Id,
		SecreteId:   tokenAuthInfo.Secret,
		Application: ConvertToProtobufApplication(tokenAuthInfo.Application),
		User:        ConvertToProtobufUser(tokenAuthInfo.User),
		SystemUser:  ConvertToProtobufSystemUser(tokenAuthInfo.SystemUserAuthInfo),
		Permission:  ConvertToProtobufPermission(model.Permission{Actions: tokenAuthInfo.Actions}),
		ExpireInfo:  ConvertToProtobufExpireInfo(model.ExpireInfo{ExpireAt: tokenAuthInfo.ExpiredAt}),
		Gateways:    ConvertToProtobufGateWays([]model.Gateway{tokenAuthInfo.Gateway}),
	}
	status.Ok = true
	logger.Debugf("Get database auth info success by token: %s", req.Token)
	return &pb.DBTokenResponse{Status: &status, Data: &dbTokenInfo}, nil
}

func (j *JMServer) CreateSession(ctx context.Context, req *pb.SessionCreateRequest) (*pb.SessionCreateResponse, error) {
	var (
		status pb.Status
	)
	apiSess := ConvertToSession(req.Data)
	apiResp, err := j.apiClient.CreateSession(apiSess)
	if err != nil {
		logger.Errorf("Create session failed: %s", err.Error())
		status.Err = err.Error()
		return &pb.SessionCreateResponse{Status: &status}, nil
	}
	status.Ok = true
	j.beat.StoreSessionId(apiResp.ID)
	logger.Debugf("Creat session %s", apiResp.ID)
	return &pb.SessionCreateResponse{Status: &status,
		Data: ConvertToProtobufSession(apiResp)}, nil
}

func (j *JMServer) FinishSession(ctx context.Context, req *pb.SessionFinishRequest) (*pb.SessionFinishResp, error) {
	var (
		status pb.Status
	)
	_, err := j.apiClient.SessionFinished(req.Id, ConvertUTCTime(req.DateEnd))
	if err != nil {
		logger.Errorf("Finish Session failed: %s", err.Error())
		status.Err = err.Error()
		return &pb.SessionFinishResp{Status: &status}, nil
	}
	status.Ok = true
	j.beat.RemoveSessionId(req.Id)
	logger.Debugf("Finish Session %s", req.Id)
	return &pb.SessionFinishResp{Status: &status}, nil
}

func (j *JMServer) UploadReplayFile(ctx context.Context, req *pb.ReplayRequest) (*pb.ReplayResponse, error) {
	j.uploader.UploadReplay(req.SessionId, req.ReplayFilePath)
	status := pb.Status{Ok: true}
	logger.Debugf("Upload Replay File Session %s path %s ", req.SessionId, req.ReplayFilePath)
	return &pb.ReplayResponse{Status: &status}, nil
}

func (j *JMServer) UploadCommand(ctx context.Context, req *pb.CommandRequest) (*pb.CommandResponse, error) {
	var (
		status pb.Status
	)
	cmd := ConvertToCommand(req)
	j.uploader.UploadCommand(&cmd)
	status.Ok = true
	logger.Debugf("Upload command session %s %s", req.Sid, req.Asset)
	return &pb.CommandResponse{Status: &status}, nil
}

func (j *JMServer) DispatchStreamingTask(stream pb.Service_DispatchStreamingTaskServer) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go j.sendStreamTask(ctx, stream)
	for {
		taskReq, err := stream.Recv()
		if err != nil {
			logger.Errorf("Dispatch Task streaming err: %v", err)
			if err == io.EOF {
				return nil
			}
			return err
		}
		j.handleTerminalTask(taskReq)
	}
}

func (j *JMServer) sendStreamTask(ctx context.Context, stream pb.Service_DispatchStreamingTaskServer) {
	taskChan := j.beat.GetTerminalTaskChan()
	for {
		select {
		case <-ctx.Done():
			return
		case task := <-taskChan:
			var pbTask pb.TerminalTask
			switch task.Name {
			case model.TaskKillSession:
				pbTask.SessionId = task.Args
				pbTask.Id = task.ID
				pbTask.Action = pb.TaskAction_KillSession
				pbTask.TerminatedBy = task.Kwargs.TerminatedBy
				logger.Infof("Send terminal task %s", task.ID)
				if err := stream.Send(&pb.TaskResponse{Task: &pbTask}); err != nil {
					logger.Errorf("Send terminal task stream err: %s", err.Error())
				}
			}
		}
	}
}
func (j *JMServer) handleTerminalTask(req *pb.TaskRequest) {
	if err := j.beat.FinishTask(req.TaskId); err != nil {
		logger.Errorf("Handle task id %s failed: %s", req.TaskId, err)
	}
}
