package common

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/jumpserver/wisp/pkg/jms-sdk-go/model"
	"github.com/jumpserver/wisp/pkg/jms-sdk-go/service"
	"github.com/jumpserver/wisp/pkg/logger"
)

type UploaderService struct {
	replayBackend  ReplayStorage
	commandBackend CommandStorage
	replayChan     chan string // Todo:
	commandChan    chan *model.Command
	mux            sync.Mutex

	commandCfg     atomic.Value // model.CommandConfig
	replayCfg      atomic.Value // model.ReplayConfig

	apiClient *service.JMService

	wg sync.WaitGroup
}

func (u *UploaderService) WatchConfig(ctx context.Context) {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			logger.Info("Upload service watch config done")
		case <-ticker.C:
			termCfg, err := u.apiClient.GetTerminalConfig()
			if err != nil {
				logger.Errorf("Upload service watch config err: %s", err)
				continue
			}
			u.commandCfg.Store(termCfg.CommandStorage)
			u.replayCfg.Store(termCfg.ReplayStorage)
		}
	}
}

func (u *UploaderService) uploadCommand() {

}

func (u *UploaderService) uploadReplay() {

}

func (u *UploaderService) getCommandBackend() CommandStorage {
	cfg := u.commandCfg.Load().(model.CommandConfig)
	return NewCommandBackend(u.apiClient, &cfg)
}

func (u *UploaderService) getReplayBackend() ReplayStorage {
	cfg := u.replayCfg.Load().(model.ReplayConfig)
	return NewReplayBackend(u.apiClient, &cfg)
}

func (u *UploaderService) run(ctx context.Context) {
	cmdList := make([]*model.Command, 0, 10)
	notificationList := make([]*model.Command, 0, 10)
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	maxRetry := 0
	for {
		select {
		case <-ctx.Done():
			u.wg.Wait()
			logger.Info("Upload Service run done")
			return
		case p := <-u.commandChan:
			if p.RiskLevel == model.DangerLevel {
				notificationList = append(notificationList, p)
			}
			cmdList = append(cmdList, p)
			if len(cmdList) < 5 {
				continue
			}
		case <-ticker.C:
			if len(cmdList) == 0 {
				continue
			}
		}
		if len(notificationList) > 0 {
			if err := u.apiClient.NotifyCommand(notificationList); err == nil {
				notificationList = notificationList[:0]
			} else {
				logger.Errorf("Upload Service command notify err: %s", err)
			}
		}
		commandBackend := u.getCommandBackend()
		if err := commandBackend.BulkSave(cmdList); err != nil {
			logger.Errorf("Upload Service command bulk save err: %s", err)
			maxRetry++
		}
		logger.Infof("Upload Service command backend %s upload %d commands",
			commandBackend.TypeName(), len(cmdList))
		cmdList = cmdList[:0]
		maxRetry = 0

	}
}
