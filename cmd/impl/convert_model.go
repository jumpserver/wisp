package impl

import (
	"github.com/jumpserver/wisp/pkg/jms-sdk-go/common"
	"github.com/jumpserver/wisp/pkg/jms-sdk-go/model"
	pb "github.com/jumpserver/wisp/protobuf-go/protobuf"
)

var modelLoginFrom = map[pb.Session_LoginFrom]string{
	pb.Session_WT: model.LoginFromWT,
	pb.Session_ST: model.LoginFromST,
	pb.Session_RT: model.LoginFromRT,
	pb.Session_DT: model.LoginFromDT,
}

func ConvertModelLoginFrom(lf pb.Session_LoginFrom) string {
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
	}
}

func ConvertToCommand(cmd *pb.CommandRequest) model.Command {
	utc := ConvertUTCTime(cmd.Timestamp)
	return model.Command{
		SessionID:   cmd.Sid,
		OrgID:       cmd.OrgId,
		Input:       cmd.Input,
		Output:      cmd.Output,
		User:        cmd.User,
		Server:      cmd.Asset,
		SystemUser:  cmd.SystemUser,
		Timestamp:   cmd.Timestamp,
		RiskLevel:   convertRiskLevel(cmd.RiskLevel),
		DateCreated: utc.UTC(),
	}
}

func ConvertUTCTime(t int64) common.UTCTime {
	return common.ParseUnixTime(t)
}

func convertRiskLevel(lvl pb.RiskLevel) int64 {
	switch lvl {
	case pb.RiskLevel_Danger:
		return model.DangerLevel
	case pb.RiskLevel_Normal:
		return model.NormalLevel
	default:
		return model.NormalLevel

	}
}

func ConvertToReqInfo(req *pb.ReqInfo) model.ReqInfo {
	return model.ReqInfo{
		Method: req.GetMethod(),
		URL:    req.GetUrl(),
	}
}
