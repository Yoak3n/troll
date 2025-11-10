package handler

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/Yoak3n/troll/viewer/service/ws"
)

var handlerState *HandlerState

type HandlerState struct {
	tasks []ws.TaskProcessData
	mu    sync.RWMutex
	wg    sync.WaitGroup
}

func InitHandlerState() {
	handlerState = &HandlerState{
		tasks: make([]ws.TaskProcessData, 0),
		mu:    sync.RWMutex{},
		wg:    sync.WaitGroup{},
	}
	go handlerState.HandleTask()
	go handlerState.SyncHandlerState()
}

func (s *HandlerState) SyncHandlerState() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for range ticker.C {
		syncHandlerStateTasks()
	}
}

func syncHandlerStateTasks() {
	handlerState.mu.RLock()
	if len(handlerState.tasks) > 0 {
		wsMsg := ws.WebsocketMessage{
			Action: ws.TasksProcess,
			Data:   handlerState.tasks,
		}
		handlerState.mu.RUnlock()
		buf, err := json.Marshal(wsMsg)
		if err != nil {
			return
		}
		ws.Hub.Broadcast(buf)
	}

}
