package handler

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/Yoak3n/troll/viewer/service/ws"
)

var handlerState *HandlerState

type HandlerState struct {
	tasks   map[string]ws.TaskProcessData
	mu      sync.RWMutex
	wg      sync.WaitGroup
	logChan chan ws.LogMessageData
}

func InitHandlerState() {
	handlerState = &HandlerState{
		tasks:   make(map[string]ws.TaskProcessData),
		mu:      sync.RWMutex{},
		wg:      sync.WaitGroup{},
		logChan: make(chan ws.LogMessageData, 100),
	}
	go handlerState.HandleTask()
	go handlerState.sendLog()
	go handlerState.SyncHandlerState()
}

func (s *HandlerState) SyncHandlerState() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for range ticker.C {
		syncHandlerStateTasks()
	}
}

func (s *HandlerState) sendLog() {
	for logMsg := range s.logChan {
		wsMsg := ws.WebsocketMessage{
			Action: ws.LogMessage,
			Data:   logMsg,
		}
		buf, err := json.Marshal(wsMsg)
		if err != nil {
			return
		}
		ws.Hub.Broadcast(buf)
	}
}

func (s *HandlerState) Log(content string) {
	s.logChan <- ws.NewLogMessageData(content)
}

func syncHandlerStateTasks() {
	handlerState.mu.RLock()
	if len(handlerState.tasks) > 0 {
		wsMsg := ws.WebsocketMessage{
			Action: ws.TasksProcess,
			Data:   handlerState.tasks,
		}

		buf, err := json.Marshal(wsMsg)
		if err != nil {
			return
		}
		ws.Hub.Broadcast(buf)
	}
	handlerState.mu.RUnlock()
}
