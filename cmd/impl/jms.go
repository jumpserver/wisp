package impl

import (
	"context"
	"fmt"
	"io"
	"net"
	"strconv"

	"github.com/jumpserver/wisp/cmd/common"
	"github.com/jumpserver/wisp/pkg/forward"
	modelCommon "github.com/jumpserver/wisp/pkg/jms-sdk-go/common"
	"github.com/jumpserver/wisp/pkg/jms-sdk-go/model"
	"github.com/jumpserver/wisp/pkg/jms-sdk-go/service"
	"github.com/jumpserver/wisp/pkg/logger"
	pb "github.com/jumpserver/wisp/protobuf-go/protobuf"
)

func NewJMServer(apiClient *service.JMService, uploader *common.UploaderService,
	beat *common.BeatService) *JMServer {
	return &JMServer{
		apiClient:    apiClient,
		uploader:     uploader,
		beat:         beat,
		forwardStore: common.NewForwardCache(),
	}
}

type JMServer struct {
	pb.UnimplementedServiceServer
	apiClient *service.JMService

	uploader *common.UploaderService
	beat     *common.BeatService

	forwardStore *common.ForwardCache
}

func (j *JMServer) GetDBTokenAuthInfo(ctx context.Context, req *pb.TokenRequest) (*pb.DBTokenResponse, error) {
	var status pb.Status
	tokenResp, err := j.apiClient.GetConnectTokenAuth(req.Token)
	if err != nil {
		status.Err = err.Error()
		logger.Errorf("Get Connect token auth failed: %s", err)
		return &pb.DBTokenResponse{Status: &status}, nil
	}
	tokenAuthInfo := tokenResp.Info
	if tokenAuthInfo.TypeName != model.ConnectApplication {
		msg := fmt.Sprintf("Bad request: token %s connect type not %s", req.Token,
			model.ConnectApplication)
		status.Err = msg
		logger.Error(msg)
		return &pb.DBTokenResponse{Status: &status}, nil
	}
	setting := j.uploader.GetTerminalSetting()
	dbTokenInfo := pb.TokenAuthInfo{
		KeyId:       tokenAuthInfo.Id,
		SecreteId:   tokenAuthInfo.Secret,
		Application: ConvertToProtobufApplication(tokenAuthInfo.Application),
		User:        ConvertToProtobufUser(tokenAuthInfo.User),
		FilterRules: ConvertToProtobufFilterRules(tokenAuthInfo.CmdFilterRules),
		SystemUser:  ConvertToProtobufSystemUser(tokenAuthInfo.SystemUserAuthInfo),
		Permission:  ConvertToProtobufPermission(model.Permission{Actions: tokenAuthInfo.Actions}),
		ExpireInfo:  ConvertToProtobufExpireInfo(model.ExpireInfo{ExpireAt: tokenAuthInfo.ExpiredAt}),
		Gateways:    ConvertToProtobufGateWays(tokenAuthInfo.Domain.Gateways),
		Setting:     ConvertToPbSetting(&setting),
	}
	status.Ok = true
	logger.Debugf("Get database auth info success by token: %s", req.Token)
	return &pb.DBTokenResponse{Status: &status, Data: &dbTokenInfo}, nil
}

func (j *JMServer) RenewToken(ctx context.Context, req *pb.TokenRequest) (*pb.StatusResponse, error) {
	var status pb.Status
	res, err := j.apiClient.RenewalToken(req.Token)
	if err != nil {
		status.Err = err.Error()
		if res.Msg != "" {
			status.Err = res.Msg
		}
		logger.Errorf("Renew token %s failed: %s", req.Token, err)
		return &pb.StatusResponse{Status: &status}, nil
	}
	logger.Debugf("Renew token %s: %+v", req.Token, res)
	status.Ok = res.Ok
	if !res.Ok {
		status.Err = res.Msg
		logger.Infof("Renew token %s failed: %s", req.Token, res.Msg)
	}
	return &pb.StatusResponse{Status: &status}, nil
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

func (j *JMServer) CheckOrCreateAssetLoginTicket(ctx context.Context,
	req *pb.AssetLoginTicketRequest) (*pb.AssetLoginTicketResponse, error) {
	var (
		status pb.Status
	)

	userId := req.GetUserId()
	assetId := req.GetAssetId()
	systemUserId := req.GetSystemUserId()
	sysUsername := req.GetSystemUserUsername()
	res, err := j.apiClient.CheckIfNeedAssetLoginConfirm(userId, assetId, systemUserId, sysUsername)
	if err != nil {
		logger.Errorf("Check or create asset login ticket req %+v err: %s", req, err)
		status.Ok = false
		status.Err = err.Error()
		return &pb.AssetLoginTicketResponse{
			Status: &status}, nil
	}
	status.Ok = true

	return &pb.AssetLoginTicketResponse{
		NeedConfirm: res.NeedConfirm,
		TicketId:    res.TicketId,
		TicketInfo:  ConvertToPbTicketInfo(&res.TicketInfo),
		Status:      &status}, nil
}

func (j *JMServer) CreateForward(ctx context.Context, req *pb.ForwardRequest) (*pb.ForwardResponse, error) {
	var (
		status pb.Status
	)
	host := req.GetHost()
	port := strconv.FormatInt(int64(req.GetPort()), 10)
	dstAddr := net.JoinHostPort(host, port)
	gateways := req.GetGateways()
	client, err := common.FindAvailableGateway(gateways)
	if err != nil {
		status.Err = err.Error()
		return &pb.ForwardResponse{
			Status: &status,
		}, nil
	}
	forwardProxy := forward.SSHForward{
		Client:  client,
		DstAddr: dstAddr,
	}
	if err = forwardProxy.Start(); err != nil {
		status.Err = err.Error()
		_ = client.Close()
		logger.Errorf("Start forward proxy failed: %s", err)
		return &pb.ForwardResponse{
			Status: &status,
		}, nil
	}
	id := modelCommon.UUID()
	j.forwardStore.Add(id, &forwardProxy)
	lnAddr := forwardProxy.GetTCPAddr()
	status.Ok = true
	logger.Infof("Start forward proxy: id %s on %s", id, lnAddr.String())
	ret := &pb.ForwardResponse{
		Status: &status,
		Id:     id,
		Host:   lnAddr.IP.String(),
		Port:   int32(lnAddr.Port),
	}
	return ret, nil
}

func (j *JMServer) DeleteForward(ctx context.Context, req *pb.ForwardDeleteRequest) (*pb.StatusResponse, error) {
	var (
		status pb.Status
	)
	id := req.GetId()
	if forwardProxy := j.forwardStore.Get(id); forwardProxy != nil {
		forwardProxy.Stop()
		status.Ok = true
		j.forwardStore.Remove(id)
		logger.Infof("Forward remove id %s", id)
		return &pb.StatusResponse{Status: &status}, nil
	}
	status.Err = fmt.Sprintf("not found forward %s", id)
	logger.Errorf("Forward not found id %s", id)
	return &pb.StatusResponse{Status: &status}, nil
}

func (j *JMServer) GetPublicSetting(ctx context.Context, empty *pb.Empty) (*pb.PublicSettingResponse, error) {
	var (
		status pb.Status
	)
	data, err := j.apiClient.GetPublicSetting()
	if err != nil {
		logger.Errorf("Get public setting err: %s", err)
		status.Err = err.Error()
		return &pb.PublicSettingResponse{Status: &status}, nil
	}
	status.Ok = true
	pbSetting := pb.PublicSetting{
		XpackEnabled: data.XpackEnabled,
		ValidLicense: data.ValidLicense,
	}
	return &pb.PublicSettingResponse{Status: &status, Data: &pbSetting}, nil
}
