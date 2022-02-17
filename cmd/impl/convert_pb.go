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
		Attrs: &pb.Application_Attrs{
			Host:     app.Attrs.Host,
			Port:     int32(app.Attrs.Port),
			Database: app.Attrs.Database,
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
	if ok {
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
		LoginFrom:    sess.LoginFrom,
		RemoteAddr:   sess.RemoteAddr,
		Protocol:     sess.Protocol,
		DateStart:    sess.DateStart.Unix(),
		OrgId:        sess.OrgID,
		UserId:       sess.UserID,
		AssetId:      sess.AssetID,
		SystemUserId: sess.SystemUserID,
	}
}
