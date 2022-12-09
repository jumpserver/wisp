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

func ConvertToProtobufAsset(asset model.Asset) *pb.Asset {

	return nil
}

func ConvertToProtobufAccount(asset model.Account) *pb.Account {

	return nil
}

func ConvertToProtobufGateway(gateway model.Gateway) *pb.Gateway {
	return &pb.Gateway{
		Id:         gateway.ID,
		Name:       gateway.Name,
		Ip:         gateway.Address,
		Port:       int32(gateway.Protocols.GetProtocolPort("ssh")),
		Protocol:   "ssh",
		Username:   gateway.Account.Username,
		Password:   gateway.Account.Secret,
		PrivateKey: gateway.Account.Secret,
	}
}

func ConvertToProtobufPermission(perm model.Actions) *pb.Permission {
	return &pb.Permission{
		EnableConnect:  perm.EnableConnect(),
		EnablePaste:    perm.EnablePaste(),
		EnableCopy:     perm.EnableCopy(),
		EnableDownload: perm.EnableDownload(),
		EnableUpload:   perm.EnableUpload(),
	}
}

var ruleActionMap = map[model.CommandAction]pb.CommandACL_Action{
	model.ActionUnknown: pb.CommandACL_Unknown,
	model.ActionAccept:  pb.CommandACL_Accept,
	model.ActionReview:  pb.CommandACL_Review,
	model.ActionReject:  pb.CommandACL_Reject,
}

func ConvertToProtobufFilterRule(rule model.CommandACL) *pb.CommandACL {
	action, ok := ruleActionMap[rule.Action]
	if !ok {
		action = pb.CommandACL_Unknown
	}
	return &pb.CommandACL{
		Id:            rule.ID,
		Name:          rule.Name,
		Priority:      int32(rule.Priority),
		Action:        action,
		IsActive:      rule.IsActive,
		CommandGroups: ConvertToProtobufCommandGroup(rule.CommandGroups),
	}
}

func ConvertToProtobufCommandGroup(groups []model.CommandGroup) []*pb.CommandGroup {
	pbRules := make([]*pb.CommandGroup, 0, len(groups))
	for i := range groups {
		group := groups[i]
		pbRules = append(pbRules, &pb.CommandGroup{
			Id:         group.ID,
			Name:       group.Name,
			Type:       group.Type,
			IgnoreCase: group.IgnoreCase,
			Pattern:    group.RePattern,
			Content:    group.Content})
	}
	return pbRules
}

func ConvertToProtobufExpireInfo(info model.ExpireInfo) *pb.ExpireInfo {
	return &pb.ExpireInfo{
		ExpireAt: int64(info),
	}
}

func ConvertToProtobufGateways(gateways []model.Gateway) []*pb.Gateway {
	if len(gateways) == 0 {
		return nil
	}
	pbGateways := make([]*pb.Gateway, len(gateways))
	for i := range gateways {
		pbGateways[i] = ConvertToProtobufGateway(gateways[i])
	}
	return pbGateways
}

func ConvertToProtobufFilterRules(rules []model.CommandACL) []*pb.CommandACL {
	if len(rules) == 0 {
		return nil
	}
	pbRules := make([]*pb.CommandACL, len(rules))
	for i := range rules {
		pbRules[i] = ConvertToProtobufFilterRule(rules[i])
	}
	return pbRules
}

func ConvertToProtobufSession(sess model.Session) *pb.Session {
	return &pb.Session{
		Id:         sess.ID,
		User:       sess.User,
		Asset:      sess.Asset,
		Account:    sess.Account,
		LoginFrom:  ConvertToPbLoginFrom(sess.LoginFrom),
		RemoteAddr: sess.RemoteAddr,
		Protocol:   sess.Protocol,
		DateStart:  sess.DateStart.Unix(),
		OrgId:      sess.OrgID,
		UserId:     sess.UserID,
		AssetId:    sess.AssetID,
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
