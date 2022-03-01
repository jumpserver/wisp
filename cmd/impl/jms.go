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
		logger.Errorf("Get Connect token auth failed: %s", err)
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
	logger.Debugf("Upload Replay File Session %s path %s ", req.SessionId, req.ReplayFilePath)
	status := pb.Status{Ok: true}
	if err := j.uploader.UploadReplay(req.SessionId, req.ReplayFilePath); err != nil {
		status.Ok = false
		status.Err = err.Error()
		logger.Errorf("Upload Replay File Session %s path %s: %s ", req.SessionId,
			req.ReplayFilePath, err)
	}
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

func (j *JMServer) DispatchTask(stream pb.Service_DispatchTaskServer) error {
	ctx, cancel := context.WithCancel(stream.Context())
	defer cancel()
	go j.sendStreamTask(ctx, stream)
	for {
		taskReq, err := stream.Recv()
		if err != nil {
			msg := fmt.Sprintf("Dispatch Task streaming err: %v", err)
			if err == io.EOF {
				logger.Infof(msg)
				return nil
			}
			logger.Errorf(msg)
			return err
		}
		j.handleTerminalTask(taskReq)
	}
}

func (j *JMServer) ScanRemainReplays(ctx context.Context, req *pb.RemainReplayRequest) (*pb.RemainReplayResponse, error) {
	status := pb.Status{Ok: true}
	ret := j.uploader.UploadRemainReplays(req.GetReplayDir())
	if len(ret.FailureErrs) > 0 {
		status.Ok = false
		status.Err = fmt.Sprintf("there are %d errs in FailureErrs", len(ret.FailureErrs))
	}
	return &pb.RemainReplayResponse{
		Status:       &status,
		SuccessFiles: ret.SuccessFiles,
		FailureFiles: ret.FailureFiles,
		FailureErrs:  ret.FailureErrs,
	}, nil
}

func (j *JMServer) sendStreamTask(ctx context.Context, stream pb.Service_DispatchTaskServer) {
	taskChan := j.beat.GetTerminalTaskChan()
	for {
		select {
		case <-ctx.Done():
			logger.Info("Send terminal task stop")
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
func (j *JMServer) handleTerminalTask(req *pb.FinishedTaskRequest) {
	if err := j.beat.FinishTask(req.TaskId); err != nil {
		logger.Errorf("Handle task id %s failed: %s", req.TaskId, err)
	}
}

func (j *JMServer) CreateCommandTicket(ctx context.Context, req *pb.CommandConfirmRequest) (*pb.CommandConfirmResponse, error) {
	var (
		status pb.Status
	)
	sid := req.GetSessionId()
	ruleId := req.GetRuleId()
	cmd := req.GetCmd()
	res, err := j.apiClient.SubmitCommandConfirm(sid, ruleId, cmd)
	if err != nil {
		logger.Errorf("Create command ticket err: %s", err)
		status.Err = err.Error()
		return &pb.CommandConfirmResponse{Status: &status}, nil
	}
	status.Ok = true
	return &pb.CommandConfirmResponse{Status: &status,
		Info: ConvertToPbTicketInfo(&res.TicketInfo),
	}, nil
}

func (j *JMServer) CheckTicketState(ctx context.Context, req *pb.TicketRequest) (*pb.TicketStateResponse, error) {
	var (
		status pb.Status
	)
	reqInfo := ConvertToReqInfo(req.Req)

	res, err := j.apiClient.CheckConfirmStatusByRequestInfo(reqInfo)
	if err != nil {
		logger.Errorf("Check ticket status %+v err: %s", reqInfo, err)
		status.Err = err.Error()
		return &pb.TicketStateResponse{Status: &status}, nil
	}
	status.Ok = true
	return &pb.TicketStateResponse{
		Data:   ConvertToPbTicketState(&res),
		Status: &status,
	}, nil
}

func (j *JMServer) CancelTicket(ctx context.Context, req *pb.TicketRequest) (*pb.StatusResponse, error) {
	var (
		status pb.Status
	)
	status.Ok = true
	reqInfo := ConvertToReqInfo(req.Req)
	if err := j.apiClient.CancelConfirmByRequestInfo(reqInfo); err != nil {
		logger.Errorf("Cancel ticket req %+v err: %s", reqInfo, err)
		status.Ok = false
		status.Err = err.Error()
	}
	return &pb.StatusResponse{Status: &status}, nil
}
