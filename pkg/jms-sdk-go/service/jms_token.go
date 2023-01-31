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

func (s *JMService) GetConnectTokenInfo(tokenId string) (resp model.ConnectToken, err error) {
	data := map[string]string{
		"id": tokenId,
	}
	_, err = s.authClient.Post(SuperConnectTokenSecretURL, data, &resp)
	return
}

func (s *JMService) RenewalToken(token string) (resp TokenRenewalResponse, err error) {
	data := map[string]string{
		"id": token,
	}
	_, err = s.authClient.Patch(SuperTokenRenewalURL, data, &resp)
	return
}

type TokenRenewalResponse struct {
	Ok  bool   `json:"ok"`
	Msg string `json:"msg"`
}
