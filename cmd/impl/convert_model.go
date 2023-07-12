package impl

import (
	"github.com/jumpserver/wisp/pkg/jms-sdk-go/common"
	"github.com/jumpserver/wisp/pkg/jms-sdk-go/model"
	pb "github.com/jumpserver/wisp/protobuf-go/protobuf"
)

var modelLoginFrom = map[pb.Session_LoginFrom]model.LabelFiled{
	pb.Session_WT: model.LoginFromWT,
	pb.Session_ST: model.LoginFromST,
	pb.Session_RT: model.LoginFromRT,
	pb.Session_DT: model.LoginFromDT,
}

func ConvertModelLoginFrom(lf pb.Session_LoginFrom) model.LabelFiled {
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
