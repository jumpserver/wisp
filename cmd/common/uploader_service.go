package common

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	modelCommon "github.com/jumpserver-dev/sdk-go/common"
	"github.com/jumpserver-dev/sdk-go/model"
	"github.com/jumpserver-dev/sdk-go/service"
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
			if p.RiskLevel >= model.WarningLevel && p.RiskLevel <= model.ReviewReject {
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
		backendName := commandBackend.TypeName()
		err := commandBackend.BulkSave(cmdList)
		if err != nil && commandBackend.TypeName() != "server" {
			logger.Warnf("Uploader service command backend %s error %s. Switch default save.",
				backendName, err)
			backendName = "server"
			err = u.apiClient.PushSessionCommand(cmdList)
		}
		if err != nil {
			logger.Errorf("Uploader service command bulk save err: %s", err)
			if maxRetry > 5 {
				cmdList = cmdList[1:]
			}
			maxRetry++
			continue
		}
		logger.Infof("Uploader service command backend %s upload %d commands",
			backendName, len(cmdList))
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
	u.recordingSessionLifecycleReplay(sid, model.ReplayUploadStart, "")
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
	replayBackendName := replayBackend.TypeName()
	if replayBackendName == "null" {
		reason := string(model.ReasonErrNullStorage)
		u.recordingSessionLifecycleReplay(sid, model.ReplayUploadFailure, reason)
		_ = os.Remove(absGzFile)
		return nil
	}

	fileInfo, err := os.Stat(absGzFile)
	if err != nil {
		logger.Errorf("Uploader service replay file %s stat failed: %s", absGzFile, err)
		return err
	}

	err = replayBackend.Upload(absGzFile, target)
	if err != nil && replayBackendName != "server" {
		u.recordingSessionLifecycleReplay(sid, model.ReplayUploadFailure, err.Error())
		logger.Errorf("Uploader service replay backend %s error %s", replayBackendName, err)
		logger.Errorf("Switch default server to upload replay %s.", absGzFile)
		replayBackendName = "server"
		u.recordingSessionLifecycleReplay(sid, model.ReplayUploadStart, "")
		err = u.apiClient.Upload(sid, absGzFile)
	}
	if err != nil {
		u.recordingSessionLifecycleReplay(sid, model.ReplayUploadFailure, err.Error())
		logger.Errorf("Uploader service replay %s uploaded error: %s", absGzFile, err)
		return err
	}
	u.recordingSessionLifecycleReplay(sid, model.ReplayUploadSuccess, "")
	logger.Infof("Uploader service replay file %s upload to %s", absGzFile, replayBackendName)
	if _, err = u.apiClient.FinishReplyWithSize(sid, fileInfo.Size()); err != nil {
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

func (u *UploaderService) UploadRemainReplays(replayDir string) {
	allRemainReplays := scanRemainReplays(u.apiClient, replayDir)
	if len(allRemainReplays) <= 0 {
		return
	}
	logger.Infof("Start to upload %d replay files 10 min later", len(allRemainReplays))
	time.Sleep(10 * time.Minute)
	logger.Debugf("Upload Remain %d replay files", len(allRemainReplays))
	for replayPath := range allRemainReplays {
		remainReplay := allRemainReplays[replayPath]
		u.recordingSessionLifecycleReplay(remainReplay.Id, model.ReplayUploadStart, "")
		if err := u.uploadRemainReplay(&remainReplay); err != nil {
			logger.Errorf("Uploader service clean remain replay %s failed: %s",
				replayPath, err)
			u.recordingSessionLifecycleReplay(remainReplay.Id, model.ReplayUploadFailure, err.Error())
			continue
		}
		u.recordingSessionLifecycleReplay(remainReplay.Id, model.ReplayUploadSuccess, "")
		logger.Infof("Uploader service upload replay %s success", replayPath)
		// 上传完成 删除原录像文件

		fileInfo, err := os.Stat(replayPath)
		if err != nil {
			logger.Errorf("Uploader service replay file %s stat failed: %s", replayPath, err)
			continue
		}
		if _, err := u.apiClient.FinishReplyWithSize(remainReplay.Id, fileInfo.Size()); err != nil {
			logger.Errorf("Uploader service notify session %s replay finished failed: %s",
				remainReplay.Id, err)
		}
		if err := os.Remove(replayPath); err != nil {
			logger.Errorf("Uploader service clean remain replay %s failed: %s",
				replayPath, err)
		}
	}
	logger.Infof("Uploader service upload remain replay files done")
}

func (u *UploaderService) recordingSessionLifecycleReplay(sid string, event model.LifecycleEvent, msgErr string) {
	logObj := model.SessionLifecycleLog{Reason: msgErr}
	if err := u.apiClient.RecordSessionLifecycleLog(sid, event, logObj); err != nil {
		logger.Errorf("Record session %s activity %s failed: %s", sid, event, err)
	}
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
