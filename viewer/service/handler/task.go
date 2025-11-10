package handler

import (
	"net/http"

	"github.com/Yoak3n/troll/viewer/service/ws"
	"github.com/google/uuid"
)

func (s *HandlerState) HandleTask() {
	for {
		task := <-ws.Hub.Tasks

		s.handleSingleTask(&task)
	}
}

func (s *HandlerState) handleSingleTask(task *ws.TaskData) {
	switch task.Type {
	case "topic":
		for _, keyword := range task.Data {
			s.handleTopicTask(keyword, task.Topic)
		}
	case "video":
	}
}

const SearchUrl = "https://api.bilibili.com/x/web-interface/wbi/search/type"

func (s *HandlerState) handleTopicTask(keyword string, topic string) {
	uuid, err := uuid.NewV7()
	if err != nil {
		return
	}
	process := ws.TaskProcessData{
		Id:    uuid.String(),
		Label: keyword,
	}
	client := http.Client{}
	req, err := http.NewRequest("GET", SearchUrl, nil)
	if err != nil {
		return
	}
	client.Do(req)
	handlerState.mu.Lock()
	handlerState.tasks = append(handlerState.tasks, process)
	handlerState.mu.Unlock()

}
