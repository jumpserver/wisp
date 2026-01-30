package common

import (
	"github.com/jumpserver-dev/sdk-go/model"
	"github.com/jumpserver-dev/sdk-go/service"
	"github.com/jumpserver-dev/sdk-go/storage"
)

func NewCommandBackend(apiClient *service.JMService, cfg *model.TerminalConfig) storage.CommandStorage {
	return storage.NewCommandStorage(apiClient, cfg)
}

func NewReplayBackend(apiClient *service.JMService, cfg *model.ReplayConfig) storage.ReplayStorage {
	return storage.NewReplayStorage(apiClient, *cfg)
}
