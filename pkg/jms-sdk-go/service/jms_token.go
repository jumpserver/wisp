package service

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/jumpserver/wisp/pkg/jms-sdk-go/model"
)

func (s *JMService) GetTokenAsset(token string) (tokenUser model.TokenUser, err error) {
	Url := fmt.Sprintf(TokenAssetURL, token)
	_, err = s.authClient.Get(Url, &tokenUser)
	return
}

func (s *JMService) GetConnectTokenAuth(token string) (resp TokenAuthInfoResponse, err error) {
	data := map[string]string{
		"token": token,
	}
	_, err = s.authClient.Post(TokenAuthInfoURL, data, &resp)
	return
}

type TokenAuthInfoResponse struct {
	Info TokenAuthInfo
	Err  []string
}

/*
	接口返回可能是一个['Token not found']
*/

func (t *TokenAuthInfoResponse) UnmarshalJSON(p []byte) error {
	if index := bytes.IndexByte(p, '['); index == 0 {
		return json.Unmarshal(p, &t.Err)
	}
	return json.Unmarshal(p, &t.Info)
}

type TokenAuthInfo struct {
	Id          string            `json:"id"`
	Secret      string            `json:"secret"`
	TypeName    model.ConnectType `json:"type"`
	User        model.User        `json:"user"`
	Actions     []string          `json:"actions,omitempty"`
	Application model.Application `json:"application,omitempty"`
	Asset       *model.Asset      `json:"asset,omitempty"`
	ExpiredAt   int64             `json:"expired_at"`
	Gateway     model.Gateway     `json:"gateway,omitempty"`
	Domain      model.Domain      `json:"domain"`

	CmdFilterRules []model.FilterRule `json:"cmd_filter_rules,omitempty"`

	SystemUserAuthInfo model.SystemUserAuthInfo `json:"system_user"`
}
