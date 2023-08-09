package impl

import (
	"fmt"
	"strconv"

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

func ConvertToProtobufProtocols(protocol []model.Protocol) []*pb.Protocol {
	if len(protocol) == 0 {
		return nil
	}
	pbProtocols := make([]*pb.Protocol, len(protocol))
	for i := range protocol {
		pbProtocols[i] = ConvertToProtobufProtocol(&protocol[i])
	}
	return pbProtocols
}

func ConvertToProtobufProtocol(protocol *model.Protocol) *pb.Protocol {
	return &pb.Protocol{
		Id:   int32(protocol.Id),
		Name: protocol.Name,
		Port: int32(protocol.Port),
	}
}

func ConvertToProtobufAsset(asset model.Asset) *pb.Asset {
	specific := asset.SpecInfo
	secretInfo := asset.SecretInfo

	protocols := asset.Protocols
	return &pb.Asset{
		Id:        asset.ID,
		Name:      asset.Name,
		Address:   asset.Address,
		OrgId:     asset.OrgID,
		Protocols: ConvertToProtobufProtocols(protocols),
		Specific: &pb.Asset_Specific{
			DbName:           specific.DBName,
			UseSsl:           specific.UseSSL,
			CaCert:           secretInfo.CaCert,
			ClientCert:       secretInfo.ClientCert,
			ClientKey:        secretInfo.ClientKey,
			AllowInvalidCert: specific.AllowInvalidCert,
			AutoFill:         specific.AutoFill,
			UsernameSelector: specific.UsernameSelector,
			PasswordSelector: specific.PasswordSelector,
			SubmitSelector:   specific.SubmitSelector,
			HttpProxy:        specific.HttpProxy,
		},
	}
}

func ConvertToProtobufAccount(account model.Account) *pb.Account {
	secretType := account.SecretType
	return &pb.Account{
		Id:       account.ID,
		Name:     account.Name,
		Username: account.Username,
		Secret:   account.Secret,
		SecretType: &pb.LabelValue{Label: secretType.Label,
			Value: secretType.Value},
	}
}

func ConvertToProtobufGateway(gateway model.Gateway) *pb.Gateway {
	account := gateway.Account
	pbGateway := &pb.Gateway{
		Id:       gateway.ID,
		Name:     gateway.Name,
		Ip:       gateway.Address,
		Port:     int32(gateway.Protocols.GetProtocolPort("ssh")),
		Protocol: "ssh",
		Username: account.Username,
	}
	if account.IsSSHKey() {
		pbGateway.PrivateKey = account.Secret
	} else {
		pbGateway.Password = account.Secret
	}

	return pbGateway
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
	model.ActionWarning: pb.CommandACL_Warning,
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
			Pattern:    group.Pattern,
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
		if gateways[i].Address == "" {
			continue
		}
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

func ConvertToPbLoginFrom(s model.LabelFiled) pb.Session_LoginFrom {
	return pbLoginFrom[s]
}

var pbLoginFrom = map[model.LabelFiled]pb.Session_LoginFrom{
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
	stateKey := string(state.State)
	return &pb.TicketState{
		State:     pbTicketMap[stateKey],
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
	return &pb.ComponentSetting{
		MaxIdleTime:    int32(setting.MaxIdleTime),
		MaxSessionTime: int32(setting.MaxSessionTime),
	}
}

func ConvertToPbPlatform(platform *model.Platform) *pb.Platform {
	return &pb.Platform{
		Id:        int32(platform.ID),
		Name:      platform.Name,
		Category:  platform.Category.Value,
		Charset:   platform.Charset.Value,
		Type:      platform.Charset.Value,
		Protocols: ConvertToPlatformProtobufProtocols(platform.Protocols),
	}
}

func ConvertToPlatformProtobufProtocols(protocols []model.PlatformProtocol) []*pb.PlatformProtocol {
	pbPlatformProtocols := make([]*pb.PlatformProtocol, 0, len(protocols))
	for i := range protocols {
		protocol := protocols[i]
		pbSetting := make(map[string]string, len(protocol.Setting))
		for k, v := range protocol.Setting {
			switch v.(type) {
			case int32, int64:
				pbSetting[k] = strconv.Itoa(int(v.(int32)))
			case string:
				pbSetting[k] = v.(string)
			case bool:
				pbSetting[k] = strconv.FormatBool(v.(bool))
			default:
				pbSetting[k] = fmt.Sprintf("%v", v)
			}
		}
		pbProtocol := &pb.PlatformProtocol{
			Id:       int32(protocol.Id),
			Name:     protocol.Name,
			Port:     int32(protocol.Port),
			Settings: pbSetting,
		}
		pbPlatformProtocols = append(pbPlatformProtocols, pbProtocol)
	}
	return pbPlatformProtocols
}
