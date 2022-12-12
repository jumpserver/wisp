package service

import (
	"strconv"

	"github.com/jumpserver/wisp/pkg/jms-sdk-go/model"
)

func (s *JMService) GetListenPorts() (ports []int32, err error) {
	_, err = s.authClient.Get(DBListenPortsURL, &ports)
	return
}

func (s *JMService) GetAssetByPort(port int32) (app model.Asset, err error) {
	data := map[string]string{
		"port": strconv.Itoa(int(port)),
	}
	_, err = s.authClient.Get(DBPortInfoURL, &app, data)
	return
}
