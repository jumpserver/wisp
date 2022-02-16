package common

import (
	"sync"
	"time"

	"github.com/jumpserver/wisp/pkg/jms-sdk-go/model"
	"github.com/jumpserver/wisp/pkg/jms-sdk-go/service"
	"github.com/jumpserver/wisp/pkg/logger"
)

type BeatService struct {
	sessMap map[string]struct{}

	apiClient *service.JMService

	callBack func()

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
			for _, task := range tasks {
				switch task.Name {
				case model.TaskKillSession:
					// todo: grpc 双向流 发送 core 的任务
				default:

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

var empty = struct{}{}
