package common

import (
	"sync"
	"time"

	"github.com/jumpserver/wisp/pkg/jms-sdk-go/model"
	"github.com/jumpserver/wisp/pkg/jms-sdk-go/service"
	"github.com/jumpserver/wisp/pkg/logger"
)

func NewBeatService(apiClient *service.JMService) *BeatService {
	return &BeatService{
		sessMap:   make(map[string]struct{}),
		apiClient: apiClient,
		taskChan:  make(chan *model.TerminalTask, 5),
	}
}

type BeatService struct {
	sessMap map[string]struct{}

	apiClient *service.JMService

	taskChan chan *model.TerminalTask

	sync.Mutex
}

func (b *BeatService) KeepHeartBeat() {
	for {
		time.Sleep(30 * time.Second)
		data := b.getSessions()
		tasks, err := b.apiClient.TerminalHeartBeat(data)
		if err != nil {
			logger.Errorf("Keep heart beat err: %s", err)
			continue
		}
		if len(tasks) != 0 {
			for i := range tasks {
				select {
				case b.taskChan <- &tasks[i]:
				default:
					logger.Infof("Discard task %v", tasks[i])
				}
			}
		}
	}
}

func (b *BeatService) getSessions() []string {
	b.Lock()
	defer b.Unlock()
	sids := make([]string, 0, len(b.sessMap))
	for sid := range b.sessMap {
		sids = append(sids, sid)
	}
	return sids
}

var empty = struct{}{}

func (b *BeatService) StoreSessionId(sid string) {
	b.Lock()
	defer b.Unlock()
	b.sessMap[sid] = empty
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
