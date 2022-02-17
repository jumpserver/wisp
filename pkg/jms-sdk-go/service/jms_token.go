package service

import (
	"fmt"

	"github.com/jumpserver/wisp/pkg/jms-sdk-go/model"
)

func (s *JMService) GetTokenAsset(token string) (tokenUser model.TokenUser, err error) {
	Url := fmt.Sprintf(TokenAssetURL, token)
	_, err = s.authClient.Get(Url, &tokenUser)
	return
}

func (s *JMService) GetConnectTokenAuth(token string) (resp TokenAuthInfo, err error) {
	data := map[string]string{
		"token": token,
	}
	_, err = s.authClient.Post(TokenAuthInfoURL, data, &resp)
	return
}

type TokenAuthInfo struct {
	Id          string            `json:"id"`
	Secret      string            `json:"secret"`
	TypeName    model.ConnectType `json:"type"`
	User        model.User        `json:"user"`
	Actions     []string          `json:"actions,omitempty"`
	Application model.Application `json:"application,omitempty"`
	Asset       model.Asset      `json:"asset,omitempty"`
	ExpiredAt   int64             `json:"expired_at"`
	Gateway     model.Gateway     `json:"gateway,omitempty"`

	SystemUserAuthInfo model.SystemUserAuthInfo `json:"system_user"`
}
