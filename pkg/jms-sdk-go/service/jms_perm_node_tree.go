package service

import (
	"fmt"
	"github/jumpserver/wisp/pkg/jms-sdk-go/model"
)

func (s *JMService) GetNodeTreeByUserAndNodeKey(userID, nodeKey string) (nodeTrees model.NodeTreeList, err error) {
	payload := map[string]string{}
	if nodeKey != "" {
		payload["key"] = nodeKey
	}
	apiURL := fmt.Sprintf(UserPermsNodeTreeWithAssetURL, userID)
	_, err = s.authClient.Get(apiURL, &nodeTrees, payload)
	return
}
