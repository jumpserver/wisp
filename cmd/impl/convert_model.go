package impl

import (
	"github.com/jumpserver-dev/sdk-go/common"
	"github.com/jumpserver-dev/sdk-go/model"
	pb "github.com/jumpserver/wisp/protobuf-go/protobuf"
)

var modelLoginFrom = map[pb.Session_LoginFrom]model.LabelField{
	pb.Session_WT: model.LoginFromWT,
	pb.Session_ST: model.LoginFromST,
	pb.Session_RT: model.LoginFromRT,
	pb.Session_DT: model.LoginFromDT,
}

func ConvertModelLoginFrom(lf pb.Session_LoginFrom) model.LabelField {
	return modelLoginFrom[lf]
}

func ConvertToSession(sees *pb.Session) model.Session {
	return model.Session{
		ID:         sees.Id,
		User:       sees.User,
		Asset:      sees.Asset,
		Account:    sees.Account,
		LoginFrom:  ConvertModelLoginFrom(sees.LoginFrom),
		RemoteAddr: sees.RemoteAddr,
		Protocol:   sees.Protocol,
		DateStart:  common.ParseUnixTime(sees.DateStart),
		OrgID:      sees.OrgId,
		UserID:     sees.UserId,
		AssetID:    sees.AssetId,
		AccountID:  sees.AccountId,
		Type:       model.NORMALType,
		TokenId:    sees.TokenId,
	}
}

func ConvertToCommand(cmd *pb.CommandRequest) model.Command {
	utc := ConvertUTCTime(cmd.Timestamp)
	return model.Command{
		SessionID:      cmd.Sid,
		OrgID:          cmd.OrgId,
		Input:          cmd.Input,
		Output:         cmd.Output,
		User:           cmd.User,
		Server:         cmd.Asset,
		Account:        cmd.Account,
		Timestamp:      cmd.Timestamp,
		CmdFilterAclId: cmd.CmdAclId,
		CmdGroupId:     cmd.CmdGroupId,
		RiskLevel:      convertRiskLevel(cmd.RiskLevel),
		DateCreated:    utc.UTC(),
	}
}

func ConvertUTCTime(t int64) common.UTCTime {
	return common.ParseUnixTime(t)
}

var riskLevelMap = map[pb.RiskLevel]int64{
	pb.RiskLevel_Normal:       model.NormalLevel,
	pb.RiskLevel_Warning:      model.WarningLevel,
	pb.RiskLevel_Reject:       model.RejectLevel,
	pb.RiskLevel_ReviewReject: model.ReviewReject,
	pb.RiskLevel_ReviewAccept: model.ReviewAccept,
	pb.RiskLevel_ReviewCancel: model.ReviewCancel,
}

func convertRiskLevel(lvl pb.RiskLevel) int64 {
	if v, ok := riskLevelMap[lvl]; ok {
		return v
	}
	return model.NormalLevel
}

func ConvertToReqInfo(req *pb.ReqInfo) model.ReqInfo {
	return model.ReqInfo{
		Method: req.GetMethod(),
		URL:    req.GetUrl(),
	}
}

var LifecycleEventMap = map[pb.SessionLifecycleLogRequest_EventType]model.LifecycleEvent{
	pb.SessionLifecycleLogRequest_AssetConnectSuccess:  model.AssetConnectSuccess,
	pb.SessionLifecycleLogRequest_AssetConnectFinished: model.AssetConnectFinished,
	pb.SessionLifecycleLogRequest_CreateShareLink:      model.CreateShareLink,
	pb.SessionLifecycleLogRequest_UserJoinSession:      model.UserJoinSession,
	pb.SessionLifecycleLogRequest_UserLeaveSession:     model.UserLeaveSession,
	pb.SessionLifecycleLogRequest_AdminJoinMonitor:     model.AdminJoinMonitor,
	pb.SessionLifecycleLogRequest_AdminExitMonitor:     model.AdminExitMonitor,
	pb.SessionLifecycleLogRequest_ReplayConvertStart:   model.ReplayConvertStart,
	pb.SessionLifecycleLogRequest_ReplayConvertSuccess: model.ReplayUploadStart,
	pb.SessionLifecycleLogRequest_ReplayConvertFailure: model.ReplayConvertFailure,
	pb.SessionLifecycleLogRequest_ReplayUploadStart:    model.ReplayConvertStart,
	pb.SessionLifecycleLogRequest_ReplayUploadSuccess:  model.ReplayUploadSuccess,
	pb.SessionLifecycleLogRequest_ReplayUploadFailure:  model.ReplayUploadFailure,
}
