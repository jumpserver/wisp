package common

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/gorilla/websocket"

	"github.com/jumpserver/wisp/pkg/jms-sdk-go/common"
	"github.com/jumpserver/wisp/pkg/jms-sdk-go/model"
	"github.com/jumpserver/wisp/pkg/jms-sdk-go/service"
	"github.com/jumpserver/wisp/pkg/logger"
)

func NewBeatService(apiClient *service.JMService) *BeatService {
	return &BeatService{
		sessMap:   make(map[string]*SessionToken),
		apiClient: apiClient,
		taskChan:  make(chan *model.TerminalTask, 5),
	}
}

type SessionToken struct {
	model.Session
	TokenId string
	invalid bool
}

type BeatService struct {
	sessMap map[string]*SessionToken

	apiClient *service.JMService

	taskChan chan *model.TerminalTask

	sync.Mutex
}

func (b *BeatService) KeepHeartBeat() {
	ws, err := b.apiClient.GetWsClient()
	if err != nil {
		logger.Errorf("Start ws client failed: %s", err)
		time.Sleep(10 * time.Second)
		go b.KeepHeartBeat()
		return
	}
	logger.Info("Start ws client success")
	done := make(chan struct{}, 2)
	go b.receiveWsTask(ws, done)
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()
	if err1 := ws.WriteJSON(b.GetStatusData()); err1 != nil {
		logger.Errorf("Ws client send heartbeat data failed: %s", err1)
	}
	for {
		select {
		case <-done:
			logger.Info("Ws client closed")
			time.Sleep(10 * time.Second)
			go b.KeepHeartBeat()
			return
		case <-ticker.C:
			if err1 := ws.WriteJSON(b.GetStatusData()); err1 != nil {
				logger.Errorf("Ws client write stat data failed: %s", err1)
				continue
			}
			logger.Debug("Ws client send heartbeat success")
		}
	}
}

func (b *BeatService) receiveWsTask(ws *websocket.Conn, done chan struct{}) {
	defer close(done)
	for {
		msgType, message, err2 := ws.ReadMessage()
		if err2 != nil {
			logger.Errorf("Ws client read err: %s", err2)
			return
		}
		switch msgType {
		case websocket.PingMessage,
			websocket.PongMessage:
			logger.Debug("Ws client ping/pong Message")
			continue
		case websocket.CloseMessage:
			logger.Debug("Ws client close Message")
			return
		}
		var tasks []model.TerminalTask
		if err := json.Unmarshal(message, &tasks); err != nil {
			logger.Errorf("Ws client Unmarshal failed: %s", err)
			continue
		}
		if len(tasks) != 0 {
			for i := range tasks {
				b.sendTask(&tasks[i])
			}
		}
	}
}

func (b *BeatService) GetStatusData() interface{} {
	sessions := b.getSessionIds()
	payload := model.HeartbeatData{
		SessionOnlineIds: sessions,
		CpuUsed:          common.CpuLoad1Usage(),
		MemoryUsed:       common.MemoryUsagePercent(),
		DiskUsed:         common.DiskUsagePercent(),
		SessionOnline:    len(sessions),
	}
	return map[string]interface{}{
		"type":    "status",
		"payload": payload,
	}
}

func (b *BeatService) getSessionIds() []string {
	b.Lock()
	defer b.Unlock()
	sids := make([]string, 0, len(b.sessMap))
	for sid := range b.sessMap {
		sids = append(sids, sid)
	}
	return sids
}

func (b *BeatService) StoreSessionId(sess *SessionToken) {
	b.Lock()
	defer b.Unlock()
	b.sessMap[sess.ID] = sess
}

func (b *BeatService) GetSessions() []*SessionToken {
	b.Lock()
	defer b.Unlock()
	sids := make([]*SessionToken, 0, len(b.sessMap))
	for sid := range b.sessMap {
		sids = append(sids, b.sessMap[sid])
	}
	return sids
}

func (b *BeatService) RemoveSessionId(sid string) {
	b.Lock()
	defer b.Unlock()
	delete(b.sessMap, sid)
}

func (b *BeatService) GetTerminalTaskChan() <-chan *model.TerminalTask {
	return b.taskChan
}

func (b *BeatService) FinishTask(taskId string) error {
	return b.apiClient.FinishTask(taskId)
}

func (b *BeatService) KeepCheckTokens() {
	for {
		time.Sleep(5 * time.Minute)
		sessions := b.GetSessions()
		tokens := make(map[string]model.TokenCheckStatus, len(sessions))
		for _, s := range sessions {
			ret, ok := tokens[s.TokenId]
			if ok {
				b.handleTokenCheck(s, &ret)
				continue
			}
			ret, err := b.apiClient.CheckTokenStatus(s.TokenId)
			if err != nil && ret.Code == "" {
				logger.Errorf("Check token status failed: %s", err)
				continue
			}
			tokens[s.TokenId] = ret
			b.handleTokenCheck(s, &ret)
		}
	}
}

func (b *BeatService) sendTask(task *model.TerminalTask) {
	select {
	case b.taskChan <- task:
	default:
		logger.Errorf("Discard task %v", task)
	}
}

func (b *BeatService) handleTokenCheck(session *SessionToken, tokenStatus *model.TokenCheckStatus) {
	var action string
	switch tokenStatus.Code {
	case model.CodePermOk:
		action = model.TaskPermValid
		if !session.invalid {
			return
		}
		session.invalid = false
	default:
		if session.invalid {
			return
		}
		session.invalid = true
		action = model.TaskPermExpired
	}
	task := model.TerminalTask{
		Name:        action,
		Args:        session.ID,
		TokenStatus: *tokenStatus,
	}
	b.sendTask(&task)
}
