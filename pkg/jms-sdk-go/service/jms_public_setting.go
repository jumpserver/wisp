package service

import "github/jumpserver/wisp/pkg/jms-sdk-go/model"

func (s *JMService) GetPublicSetting() (result model.PublicSetting, err error) {
	var response struct {
		Data model.PublicSetting `json:"data"`
	}
	client := s.authClient.Clone()
	_, err = client.Get(PublicSettingURL, &response)
	result = response.Data
	return
}
