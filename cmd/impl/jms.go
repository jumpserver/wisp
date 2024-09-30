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

func (j *JMServer) GetTokenAuthInfo(ctx context.Context, req *pb.TokenRequest) (*pb.TokenResponse, error) {
	var status pb.Status
	var gateways []model.Gateway
	tokenAuthInfo, err := j.apiClient.GetConnectTokenInfo(req.Token)
	if err != nil {
		status.Err = err.Error()
		logger.Errorf("Get Connect token auth failed: %s", err)
		return &pb.TokenResponse{Status: &status}, nil
	}
	if tokenAuthInfo.Gateway != nil {
		gateways = append(gateways, *tokenAuthInfo.Gateway)
	}
	setting := j.uploader.GetTerminalSetting()
	dbTokenInfo := pb.TokenAuthInfo{
		KeyId:       tokenAuthInfo.Id,
		SecreteId:   tokenAuthInfo.Value,
		Asset:       ConvertToProtobufAsset(tokenAuthInfo.Asset),
		User:        ConvertToProtobufUser(tokenAuthInfo.User),
		FilterRules: ConvertToProtobufFilterRules(tokenAuthInfo.CommandFilterACLs),
		Account:     ConvertToProtobufAccount(tokenAuthInfo.Account),
		Permission:  ConvertToProtobufPermission(tokenAuthInfo.Actions),
		ExpireInfo:  ConvertToProtobufExpireInfo(tokenAuthInfo.ExpireAt),
		Gateways:    ConvertToProtobufGateways(gateways),
		Setting:     ConvertToPbSetting(&setting),
		Platform:    ConvertToPbPlatform(&tokenAuthInfo.Platform),
	}
	status.Ok = true
	logger.Debugf("Get database auth info success by token: %s", req.Token)
	return &pb.TokenResponse{Status: &status, Data: &dbTokenInfo}, nil
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
	sessionToken := common.SessionToken{
		Session: apiResp,
		TokenId: req.Data.TokenId,
	}
	j.beat.StoreSessionId(&sessionToken)
	logger.Debugf("Creat session %s", apiSess.ID)
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
			pbTask.SessionId = task.Args
			pbTask.Id = task.ID
			pbTask.TerminatedBy = task.Kwargs.TerminatedBy
			switch task.Name {
			case model.TaskKillSession:
				pbTask.Action = pb.TaskAction_KillSession
			case model.TaskLockSession:
				pbTask.CreatedBy = task.Kwargs.CreatedByUser
				pbTask.Action = pb.TaskAction_LockSession
			case model.TaskUnlockSession:
				pbTask.Action = pb.TaskAction_UnlockSession
				pbTask.CreatedBy = task.Kwargs.CreatedByUser
			case model.TaskPermExpired:
				pbTask.Action = pb.TaskAction_TokenPermExpired
				pbTask.TokenStatus = &pb.TokenStatus{
					Code:      "",
					Detail:    "",
					IsExpired: false,
				}

			case model.TaskPermValid:
				pbTask.Action = pb.TaskAction_TokenPermValid

			default:
				logger.Errorf("Unknown task name %s", task.Name)
				continue
			}
			logger.Infof("Send terminal task %s name: %s", task.ID, task.Name)
			if err := stream.Send(&pb.TaskResponse{Task: &pbTask}); err != nil {
				logger.Errorf("Send terminal task stream err: %s", err.Error())
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
	aclId := req.GetCmdAclId()
	cmd := req.GetCmd()
	res, err := j.apiClient.SubmitCommandReview(sid, aclId, cmd)
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
		pbStatus pb.Status
	)

	userId := req.GetUserId()
	assetId := req.GetAssetId()
	username := req.GetAccountUsername()
	res, err := j.apiClient.CheckIfNeedAssetLoginConfirm(userId, assetId, username)
	if err != nil {
		logger.Errorf("Check or create asset login ticket req %+v err: %s", req, err)
		pbStatus.Ok = false
		pbStatus.Err = err.Error()
		return &pb.AssetLoginTicketResponse{
			Status: &pbStatus}, nil
	}
	pbStatus.Ok = true

	return &pb.AssetLoginTicketResponse{
		NeedConfirm: res.NeedConfirm,
		TicketId:    res.TicketId,
		TicketInfo:  ConvertToPbTicketInfo(&res.TicketInfo),
		Status:      &pbStatus}, nil
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
	setting := j.uploader.GetTerminalSetting()
	status.Ok = true
	pbSetting := pb.PublicSetting{
		XpackEnabled: data.XpackEnabled,
		ValidLicense: data.ValidLicense,
		GptBaseUrl:   setting.GPTBaseURL,
		GptApiKey:    setting.GPTApiKey,
		GptProxy:     setting.GPTProxy,
		GptModel:     setting.GPTModel,
	}
	return &pb.PublicSettingResponse{Status: &status, Data: &pbSetting}, nil
}

func (j *JMServer) CheckUserByCookies(ctx context.Context, cookiesReq *pb.CookiesRequest) (*pb.UserResponse, error) {
	var pbStatus pb.Status
	cookies := cookiesReq.GetCookies()
	cookiesMap := make(map[string]string)
	for _, cookie := range cookies {
		cookiesMap[cookie.Name] = cookie.Value
	}
	user, err := j.apiClient.CheckUserCookie(cookiesMap)
	if err != nil {
		pbStatus.Err = err.Error()
		return &pb.UserResponse{Status: &pbStatus}, nil
	}
	pbStatus.Ok = true
	pbUser := ConvertToProtobufUser(*user)
	return &pb.UserResponse{Status: &pbStatus, Data: pbUser}, nil
}
