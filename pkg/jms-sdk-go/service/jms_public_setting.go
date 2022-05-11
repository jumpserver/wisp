package service

import "github.com/jumpserver/wisp/pkg/jms-sdk-go/model"

func (s *JMService) GetPublicSetting() (result model.PublicSetting, err error) {
	_, err = s.authClient.Get(PublicSettingURL, &result)
	return
}
