package common

import (
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	modelCommon "github.com/jumpserver/wisp/pkg/jms-sdk-go/common"
	"github.com/jumpserver/wisp/pkg/jms-sdk-go/model"
	"github.com/jumpserver/wisp/pkg/jms-sdk-go/service"
	"github.com/jumpserver/wisp/pkg/logger"
)

func NewUploader(apiClient *service.JMService,
	cfg *model.TerminalConfig) *UploaderService {
	uploader := UploaderService{
		commandChan: make(chan *model.Command, 10),
		apiClient:   apiClient,
		closed:      make(chan struct{}),
	}
	uploader.updateBackendCfg(cfg)
	return &uploader
}

type UploaderService struct {
	commandChan chan *model.Command
	closed      chan struct{}
	apiClient   *service.JMService

	commandCfg atomic.Value // model.CommandConfig
	replayCfg  atomic.Value // model.ReplayConfig

	wg sync.WaitGroup
}

func (u *UploaderService) Start() {
	u.wg.Add(2)
	go u.run()
	go u.watchConfig()
}

func (u *UploaderService) watchConfig() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()
	defer u.wg.Done()
	for {
		select {
		case <-u.closed:
			logger.Info("Uploader service watch config task done")
			return
		case <-ticker.C:
			termCfg, err := u.apiClient.GetTerminalConfig()
			if err != nil {
				logger.Errorf("Upload service watch config err: %s", err)
				continue
			}
			u.updateBackendCfg(&termCfg)
		}
	}
}

func (u *UploaderService) updateBackendCfg(termCfg *model.TerminalConfig) {
	u.commandCfg.Store(termCfg.CommandStorage)
	u.replayCfg.Store(termCfg.ReplayStorage)
}

func (u *UploaderService) getCommandBackend() CommandStorage {
	cfg := u.commandCfg.Load().(model.CommandConfig)
	return NewCommandBackend(u.apiClient, &cfg)
}

func (u *UploaderService) getReplayBackend() ReplayStorage {
	cfg := u.replayCfg.Load().(model.ReplayConfig)
	return NewReplayBackend(u.apiClient, &cfg)
}

func (u *UploaderService) run() {
	cmdList := make([]*model.Command, 0, 10)
	notificationList := make([]*model.Command, 0, 10)
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	defer u.wg.Done()
	maxRetry := 0
	for {
		select {
		case <-u.closed:
			logger.Info("Uploader Service command task done")
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
				logger.Errorf("Uploader Service command notify err: %s", err)
			}
		}
		commandBackend := u.getCommandBackend()
		if err := commandBackend.BulkSave(cmdList); err != nil {
			logger.Errorf("Uploader Service command bulk save err: %s", err)
			maxRetry++
		}
		logger.Infof("Uploader Service command backend %s upload %d commands",
			commandBackend.TypeName(), len(cmdList))
		cmdList = cmdList[:0]
		maxRetry = 0

	}
}

const dateTimeFormat = "2006-01-02"

func (u *UploaderService) UploadReplay(sid, replayPath string) {
	u.wg.Add(1)
	defer u.wg.Done()
	if !HaveFile(replayPath) {
		logger.Errorf("Replay file not found: %s ", replayPath)
		return
	}
	sess, err := u.apiClient.GetSessionById(sid)
	if err != nil {
		logger.Errorf("Retrieve session %s detail failed:  %s", sid, err)
		return
	}
	today := sess.DateStart.UTC().Format(dateTimeFormat)
	absGzFile := replayPath
	if !isGzipFile(absGzFile) {
		absGzFile = absGzFile + model.SuffixGz
		if err = modelCommon.CompressToGzipFile(replayPath, absGzFile); err != nil {
			logger.Errorf("Gzip Compress replay file %s failed: %s", replayPath, err)
			return
		}
		defer os.Remove(replayPath)
		logger.Infof("Gzip Compress completed and will remove %s", replayPath)
	}
	replayBackend := u.getReplayBackend()
	gzFilename := filepath.Base(absGzFile)
	target := strings.Join([]string{today, gzFilename}, "/")
	err = replayBackend.Upload(absGzFile, target)
	if err != nil {
		logger.Errorf("Upload Replay file %s failed: %s", absGzFile, err)
		return
	}
	logger.Infof("Upload Replay file %s by %s", absGzFile, replayBackend.TypeName())

	if _, err = u.apiClient.FinishReply(sid); err != nil {
		logger.Errorf("Finish session replay ")
	}
	if err = os.Remove(absGzFile); err != nil {
		logger.Errorf("Remove replay file %s failed: %s", absGzFile, err)
		return
	}
}

func (u *UploaderService) UploadCommand(cmd *model.Command) {
	u.commandChan <- cmd
}

func (u *UploaderService) Stop() {
	select {
	case <-u.closed:
	default:
		close(u.closed)
	}
	u.wg.Wait()
	logger.Info("Uploader Service stop")
}

func HaveFile(src string) bool {
	info, err := os.Stat(src)
	return err == nil && !info.IsDir()
}

func isGzipFile(src string) bool {
	return strings.HasSuffix(src, model.SuffixGz)
}
