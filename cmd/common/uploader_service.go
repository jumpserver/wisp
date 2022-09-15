package common

import (
	"fmt"
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

	terminalCfg atomic.Value // *model.TerminalConfig

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
				logger.Errorf("Uploader service watch config err: %s", err)
				continue
			}
			u.updateBackendCfg(&termCfg)
		}
	}
}

func (u *UploaderService) updateBackendCfg(termCfg *model.TerminalConfig) {
	u.commandCfg.Store(termCfg.CommandStorage)
	u.replayCfg.Store(termCfg.ReplayStorage)
	u.terminalCfg.Store(termCfg)
}

func (u *UploaderService) getCommandBackend() CommandStorage {
	cfg := u.commandCfg.Load().(model.CommandConfig)
	return NewCommandBackend(u.apiClient, &cfg)
}

func (u *UploaderService) getReplayBackend() ReplayStorage {
	cfg := u.replayCfg.Load().(model.ReplayConfig)
	return NewReplayBackend(u.apiClient, &cfg)
}

func (u *UploaderService) GetTerminalSetting() model.TerminalConfig {
	cfg := u.terminalCfg.Load().(*model.TerminalConfig)
	return *cfg
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
			logger.Info("Uploader service command task done")
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
				logger.Errorf("Uploader service command notify err: %s", err)
			}
		}
		commandBackend := u.getCommandBackend()
		if err := commandBackend.BulkSave(cmdList); err != nil {
			logger.Warnf("Uploader service command backend %s error %s. Switch default save.",
				commandBackend.TypeName(), err)
			defaultCommandBackend := GetDefaultCommandStorage(u.apiClient)
			err := defaultCommandBackend.BulkSave(cmdList)
			if err != nil {
				logger.Errorf("Uploader service command bulk save err: %s", err)
				maxRetry++
				if maxRetry > 5 {
					cmdList = cmdList[1:]
				}
				continue
			}
		}
		logger.Infof("Uploader service command backend %s upload %d commands",
			commandBackend.TypeName(), len(cmdList))
		cmdList = cmdList[:0]
		maxRetry = 0

	}
}

const dateTimeFormat = "2006-01-02"

func (u *UploaderService) UploadReplay(sid, replayPath string) error {
	if !HaveFile(replayPath) {
		logger.Errorf("Replay file not found: %s ", replayPath)
		return fmt.Errorf("not found %s", replayPath)
	}
	sess, err := u.apiClient.GetSessionById(sid)
	if err != nil {
		logger.Errorf("Retrieve session %s detail failed:  %s", sid, err)
		return err
	}
	today := sess.DateStart.UTC().Format(dateTimeFormat)
	absGzFile := replayPath
	if !isGzipFile(absGzFile) {
		absGzFile = absGzFile + model.SuffixGz
		if err = modelCommon.CompressToGzipFile(replayPath, absGzFile); err != nil {
			logger.Errorf("Gzip Compress replay file %s failed: %s", replayPath, err)
			return err
		}
		defer os.Remove(replayPath)
		logger.Infof("Gzip Compress completed and will remove %s", replayPath)
	}
	replayBackend := u.getReplayBackend()
	gzFilename := filepath.Base(absGzFile)
	target := strings.Join([]string{today, gzFilename}, "/")
	err = replayBackend.Upload(absGzFile, target)
	if err != nil {
		logger.Errorf("Uploader service replay file %s failed: %s", absGzFile, err)
		return err
	}
	logger.Infof("Uploader service replay file %s by %s", absGzFile, replayBackend.TypeName())

	if _, err = u.apiClient.FinishReply(sid); err != nil {
		logger.Errorf("Finish %s replay api failed: %s", sid, err)
		return err
	}
	if err = os.Remove(absGzFile); err != nil {
		logger.Errorf("Remove replay file %s failed: %s", absGzFile, err)
	}
	return nil
}

func (u *UploaderService) UploadCommand(cmd *model.Command) {
	u.commandChan <- cmd
}

func (u *UploaderService) UploadRemainReplays(replayDir string) (ret RemainReplayResult) {
	allRemainReplays := scanRemainReplays(u.apiClient, replayDir)
	if len(allRemainReplays) <= 0 {
		return
	}
	successFiles := make([]string, 0, 10)
	failureFiles := make([]string, 0, 10)
	failureErrs := make([]string, 0, 10)
	logger.Debugf("Upload Remain %d replay files", len(allRemainReplays))
	for replayPath := range allRemainReplays {
		remainReplay := allRemainReplays[replayPath]
		if err := u.uploadRemainReplay(&remainReplay); err != nil {
			logger.Errorf("Uploader service clean remain replay %s failed: %s",
				replayPath, err)
			failureFiles = append(failureFiles, replayPath)
			failureErrs = append(failureErrs, err.Error())
			continue
		}
		successFiles = append(successFiles, replayPath)
		// 上传完成 删除原录像文件
		if err := os.Remove(replayPath); err != nil {
			logger.Errorf("Uploader service clean remain replay %s failed: %s",
				replayPath, err)
		}
		if _, err := u.apiClient.FinishReply(remainReplay.Id); err != nil {
			logger.Errorf("Uploader service notify session %s replay finished failed: %s",
				remainReplay.Id, err)
		}
	}
	ret.FailureErrs = failureErrs
	ret.FailureFiles = failureFiles
	ret.SuccessFiles = successFiles
	return
}

func (u *UploaderService) uploadRemainReplay(replay *RemainReplay) error {
	replayAbsGzPath := replay.AbsFilePath
	if !isGzipFile(replayAbsGzPath) {
		dirPath := filepath.Dir(replay.AbsFilePath)
		replayAbsGzPath = filepath.Join(dirPath, replay.GetGzFilename())
		if err := modelCommon.CompressToGzipFile(replay.AbsFilePath, replayAbsGzPath); err != nil {
			return fmt.Errorf("uploader service compress gzip file %s: %s", replay.AbsFilePath, err)
		}
		defer os.Remove(replayAbsGzPath)
	}
	replayBackend := u.getReplayBackend()
	return replayBackend.Upload(replayAbsGzPath, replay.TargetPath())
}

func scanRemainReplays(apiClient *service.JMService, replayDir string) map[string]RemainReplay {
	allRemainReplays := make(map[string]RemainReplay)
	_ = filepath.Walk(replayDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		var (
			sid        string
			targetDate string
			version    model.ReplayVersion
			ok         bool
		)
		if sid, ok = ParseReplaySessionID(info.Name()); !ok {
			return nil
		}
		if version, ok = ParseReplayVersion(info.Name()); !ok {
			version = model.Version2
		}
		finishedTime := modelCommon.NewUTCTime(info.ModTime())
		finishedSession, err2 := apiClient.SessionFinished(sid, finishedTime)
		if err2 != nil {
			logger.Errorf("Uploader service  mark session %s finished failed: %s", sid, err2)
			return nil
		}
		targetDate = finishedSession.DateStart.UTC().Format("2006-01-02")
		allRemainReplays[path] = RemainReplay{
			Id:          sid,
			Version:     version,
			TargetDate:  targetDate,
			AbsFilePath: path,
			IsGzip:      isGzipFile(info.Name()),
		}
		return nil
	})
	return allRemainReplays
}

func (u *UploaderService) Stop() {
	select {
	case <-u.closed:
	default:
		close(u.closed)
	}
	u.wg.Wait()
	logger.Info("Uploader service stop")
}

func HaveFile(src string) bool {
	info, err := os.Stat(src)
	return err == nil && !info.IsDir()
}

func isGzipFile(src string) bool {
	return strings.HasSuffix(src, model.SuffixGz)
}

type RemainReplayResult struct {
	SuccessFiles []string
	FailureFiles []string
	FailureErrs  []string
}
