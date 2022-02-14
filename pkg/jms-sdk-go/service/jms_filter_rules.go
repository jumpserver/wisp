package service

import (
	"fmt"

	"github.com/jumpserver/wisp/pkg/jms-sdk-go/model"
)

func (s *JMService) GetSystemUserFilterRules(systemUserID string) (rules []model.FilterRule, err error) {
	Url := fmt.Sprintf(SystemUserCmdFilterRulesListURL, systemUserID)
	_, err = s.authClient.Get(Url, &rules)
	return
}

func (s *JMService) GetCommandFilterRules(userId, sysId, assetId, appId string) (rules []model.FilterRule, err error) {
	param := make(map[string]string)
	if userId != "" {
		param["user_id"] = userId
	}
	if sysId != "" {
		param["system_user_id"] = sysId
	}
	if assetId != "" {
		param["asset_id"] = assetId
	}
	if appId != "" {
		param["application_id"] = appId
	}
	_, err = s.authClient.Get(CommandFilterRulesListURL, &rules, param)
	return
}
