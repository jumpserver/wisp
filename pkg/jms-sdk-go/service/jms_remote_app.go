package service

import (
	"fmt"
	"github/jumpserver/wisp/pkg/jms-sdk-go/model"
)

func (s *JMService) GetRemoteApp(remoteAppId string) (remoteApp model.RemoteAPP, err error) {
	Url := fmt.Sprintf(RemoteAPPURL, remoteAppId)
	_, err = s.authClient.Get(Url, &remoteApp)
	return
}
