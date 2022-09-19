package impl

import (
	"github.com/jumpserver/wisp/pkg/jms-sdk-go/model"
	pb "github.com/jumpserver/wisp/protobuf-go/protobuf"
)

func ConvertToProtobufUser(user model.User) *pb.User {
	return &pb.User{
		Id:       user.ID,
		Name:     user.Name,
		Username: user.Username,
		Role:     user.Role,
		IsValid:  user.IsValid,
		IsActive: user.IsActive,
	}
}

func ConvertToProtobufApplication(app model.Application) *pb.Application {
	return &pb.Application{
		Id:       app.ID,
		Name:     app.Name,
		Category: app.Category,
		TypeName: app.TypeName,
		Domain:   app.Domain,
		OrgId:    app.OrgID,
		Attrs: &pb.Application_Attrs{
			Host:       app.Attrs.Host,
			Port:       int32(app.Attrs.Port),
			Database:   app.Attrs.Database,
			UseSSL:     app.Attrs.UseSSL,
			CaCert:     app.Attrs.CaCert,
			ClientCert: app.Attrs.ClientCert,
			CertKey:    app.Attrs.CertKey,
		},
	}
}

func ConvertToProtobufGateWay(gateway model.Gateway) *pb.Gateway {
	return &pb.Gateway{
		Id:         gateway.ID,
		Name:       gateway.Name,
		Ip:         gateway.IP,
		Port:       int32(gateway.Port),
		Protocol:   gateway.Protocol,
		Username:   gateway.Username,
		Password:   gateway.Password,
		PrivateKey: gateway.PrivateKey,
	}
}

func ConvertToProtobufPermission(perm model.Permission) *pb.Permission {
	return &pb.Permission{
		EnableConnect:  perm.EnableConnect(),
		EnablePaste:    perm.EnablePaste(),
		EnableCopy:     perm.EnableCopy(),
		EnableDownload: perm.EnableDownload(),
		EnableUpload:   perm.EnableUpload(),
	}
}

func ConvertToProtobufSystemUser(systemUserAuthInfo model.SystemUserAuthInfo) *pb.SystemUserAuthInfo {
	return &pb.SystemUserAuthInfo{
		Id:         systemUserAuthInfo.ID,
		Name:       systemUserAuthInfo.Name,
		Protocol:   systemUserAuthInfo.Protocol,
		Username:   systemUserAuthInfo.Username,
		Password:   systemUserAuthInfo.Password,
		PrivateKey: systemUserAuthInfo.PrivateKey,
		Token:      systemUserAuthInfo.Token,
		AdDomain:   systemUserAuthInfo.AdDomain,
		OrgName:    systemUserAuthInfo.OrgName,
		OrgId:      systemUserAuthInfo.OrgId,
	}
}

var ruleActionMap = map[model.RuleAction]pb.FilterRule_Action{
	model.ActionUnknown: pb.FilterRule_Unknown,
	model.ActionAllow:   pb.FilterRule_Allow,
	model.ActionConfirm: pb.FilterRule_Confirm,
	model.ActionDeny:    pb.FilterRule_Deny,
}

func ConvertToProtobufFilterRule(rule model.FilterRule) *pb.FilterRule {
	action, ok := ruleActionMap[rule.Action]
	if !ok {
		action = pb.FilterRule_Unknown
	}
	return &pb.FilterRule{
		Id:         rule.ID,
		Priority:   int32(rule.Priority),
		Type:       rule.Type,
		Content:    rule.Content,
		Action:     action,
		OrgId:      rule.OrgId,
		Pattern:    rule.RePattern,
		IgnoreCase: rule.IgnoreCase,
	}
}

func ConvertToProtobufExpireInfo(info model.ExpireInfo) *pb.ExpireInfo {
	return &pb.ExpireInfo{
		ExpireAt: info.ExpireAt,
	}
}

func ConvertToProtobufGateWays(gateways []model.Gateway) []*pb.Gateway {
	if len(gateways) == 0 {
		return nil
	}
	pbGateways := make([]*pb.Gateway, len(gateways))
	for i := range gateways {
		pbGateways[i] = ConvertToProtobufGateWay(gateways[i])
	}
	return pbGateways
}

func ConvertToProtobufFilterRules(rules []model.FilterRule) []*pb.FilterRule {
	if len(rules) == 0 {
		return nil
	}
	pbRules := make([]*pb.FilterRule, len(rules))
	for i := range rules {
		pbRules[i] = ConvertToProtobufFilterRule(rules[i])
	}
	return pbRules
}

func ConvertToProtobufSession(sess model.Session) *pb.Session {
	return &pb.Session{
		Id:           sess.ID,
		User:         sess.User,
		Asset:        sess.Asset,
		SystemUser:   sess.SystemUser,
		LoginFrom:    ConvertToPbLoginFrom(sess.LoginFrom),
		RemoteAddr:   sess.RemoteAddr,
		Protocol:     sess.Protocol,
		DateStart:    sess.DateStart.Unix(),
		OrgId:        sess.OrgID,
		UserId:       sess.UserID,
		AssetId:      sess.AssetID,
		SystemUserId: sess.SystemUserID,
	}
}

func ConvertToPbLoginFrom(s string) pb.Session_LoginFrom {
	return pbLoginFrom[s]
}

var pbLoginFrom = map[string]pb.Session_LoginFrom{
	model.LoginFromWT: pb.Session_WT,
	model.LoginFromST: pb.Session_ST,
	model.LoginFromRT: pb.Session_RT,
	model.LoginFromDT: pb.Session_DT,
}

func ConvertToPbTicketInfo(info *model.TicketInfo) *pb.TicketInfo {
	return &pb.TicketInfo{
		CheckReq:        ConvertToPbReqInfo(info.CheckReq),
		CancelReq:       ConvertToPbReqInfo(info.CloseReq),
		TicketDetailUrl: info.TicketDetailUrl,
		Reviewers:       info.Reviewers,
	}
}

func ConvertToPbReqInfo(reqInfo model.ReqInfo) *pb.ReqInfo {
	return &pb.ReqInfo{
		Method: reqInfo.Method,
		Url:    reqInfo.URL,
	}
}

func ConvertToPbTicketState(state *model.TicketState) *pb.TicketState {

	return &pb.TicketState{
		State:     pbTicketMap[state.State],
		Processor: state.Processor,
	}
}

var pbTicketMap = map[string]pb.TicketState_State{
	model.TicketOpen:     pb.TicketState_Open,
	model.TicketApproved: pb.TicketState_Approved,
	model.TicketRejected: pb.TicketState_Rejected,
	model.TicketClosed:   pb.TicketState_Closed,
}

func ConvertToPbSetting(setting *model.TerminalConfig) *pb.ComponentSetting {
	return &pb.ComponentSetting{MaxIdleTime: int32(setting.MaxIdleTime)}
}
