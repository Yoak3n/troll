package handler

import (
	"fmt"

	"github.com/Yoak3n/troll/scanner/model"
	"github.com/Yoak3n/troll/scanner/package/handler"
	"github.com/Yoak3n/troll/viewer/service/ws"
)

func (s *HandlerState) HandleTask() {
	for {
		select {
		case task := <-ws.Hub.Tasks:
			s.handleSingleTask(&task)
		}
	}
}

func (s *HandlerState) handleSingleTask(task *ws.TaskData) {
	s.Log(fmt.Sprintf("handleSingleTask %v", task))
	switch task.Type {
	case "topic":
		for _, keyword := range task.Data {
			s.handleTopicTask(keyword, task.Topic, task.Page)
		}
	case "video":
		for _, d := range task.Data {
			s.handleVideoTask(d, task.Topic)
		}
	}
}

func (s *HandlerState) handleTopicTask(keyword string, topic string, page int) {
	s.Log(fmt.Sprintf("handleTopicTask %v %v %v", keyword, topic, page))
	videos := handler.SearchVideoOfTopic(keyword, 1)
	if page > 1 {
		for i := 2; i <= page; i++ {
			videos = append(videos, handler.SearchVideoOfTopic(keyword, i)...)
		}
	}
	if videos == nil {
		return
	}
	for _, video := range videos {
		process := ws.TaskProcessData{
			Id:        video.Bvid,
			Label:     video.Title,
			Total:     video.Review,
			Current:   0,
			Completed: false,
		}
		handlerState.mu.Lock()
		handlerState.tasks[process.Id] = process
		handlerState.mu.Unlock()
		go s.taskProcess(video.Bvid, topic)
	}

}

func (s *HandlerState) taskProcess(id string, topic string) {
	handlerState.mu.RLock()
	process, ok := handlerState.tasks[id]
	handlerState.mu.RUnlock()
	if !ok {
		s.setTaskCompleted(id)
		return
	}
	videoInfo := handler.FetchVideoInfo(process.Id, topic)
	if videoInfo.Review == 0 {
		s.setTaskCompleted(id)
		return
	}
	comments := s.fetchVideoComments(id, uint(videoInfo.Avid))
	if comments == nil {
		s.setTaskCompleted(id)
		return
	}
}

func (s *HandlerState) fetchVideoComments(bvid string, avid uint) []model.CommentData {
	allComments := make([]model.CommentData, 0)
	offset := ""
	count := 0
	tempCount := 0
	cb1 := func(count int) {
		tempCount += count
		handlerState.mu.Lock()
		process, ok := handlerState.tasks[bvid]
		if ok {
			process.Current = tempCount
			handlerState.tasks[bvid] = process
		}
		handlerState.mu.Unlock()
	}
	cb2 := func(count int) {
		s.Log(fmt.Sprintf("fetchVideoComments %v %v", bvid, tempCount))
	}
	for {
		temp, newCount, returnedOffset := handler.FetchVideoComments(avid, offset, cb1, cb2)
		if newCount == 0 {
			break
		}
		count += newCount
		handlerState.mu.Lock()
		process, ok := handlerState.tasks[bvid]
		if ok {
			process.Current = count
			tempCount = count
			handlerState.tasks[bvid] = process
		}
		handlerState.mu.Unlock()
		allComments = append(allComments, temp...)
		offset = returnedOffset
	}
	s.setTaskCompleted(bvid)
	return allComments
}

func (s *HandlerState) setTaskCompleted(id string) {
	handlerState.mu.Lock()
	process, ok := handlerState.tasks[id]
	if ok {
		process.Completed = true
		process.Current = process.Total
		handlerState.tasks[id] = process
	}
	handlerState.mu.Unlock()
}

func (s *HandlerState) handleVideoTask(bvid string, topic string) {
	s.Log(fmt.Sprintf("handleVideoTask %v %v", bvid, topic))
	videoInfo := handler.FetchVideoInfo(bvid, topic)
	if videoInfo.Review == 0 {
		s.Log(fmt.Sprintf("handleVideoTask %v %v review is 0", bvid, topic))
		s.setTaskCompleted(bvid)
		return
	}
	process := ws.TaskProcessData{
		Id:        bvid,
		Label:     videoInfo.Title,
		Total:     videoInfo.Review,
		Current:   0,
		Completed: false,
	}
	handlerState.mu.Lock()
	handlerState.tasks[process.Id] = process
	handlerState.mu.Unlock()
	go s.taskProcess(bvid, topic)
}
