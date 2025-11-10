package ws

import "time"

type WebsocketMessage struct {
	Action WsAction `json:"action"`
	Data   any      `json:"data"`
}
type WsAction string

const (
	ConnectedSuccess WsAction = "Connected"
	LogMessage       WsAction = "Log"
	AddTaskMessage   WsAction = "Task"
	TaskProcess      WsAction = "Process"
	TasksProcess     WsAction = "Processes"
	CloseMessage     WsAction = "Close"
	PingMessage      WsAction = "Ping"
	PongMessage      WsAction = "Pong"
)

type LogMessageData struct {
	Time    string `json:"time"`
	Content string `json:"content"`
}

func NewLogMessageData(content string) LogMessageData {
	return LogMessageData{
		Time:    time.Now().Format("2006-01-02 15:04"),
		Content: content,
	}
}

func NewLogDataToMessage(content string) *WebsocketMessage {
	return &WebsocketMessage{
		Action: LogMessage,
		Data: LogMessageData{
			Time:    time.Now().Format("2006-01-02 15:04"),
			Content: content,
		},
	}
}

type TaskProcessData struct {
	Id        string `json:"id"`
	Label     string `json:"lable"`
	Total     int    `json:"total"`
	Current   int    `json:"current"`
	Completed bool   `json:"completed"`
}

func NewTaskProcessData(id string, label string, total int, current int, completed bool) TaskProcessData {
	return TaskProcessData{
		Id:        id,
		Label:     label,
		Total:     total,
		Current:   current,
		Completed: completed,
	}
}

func NewTasksProcessesDataToMessage(data []TaskProcessData) *WebsocketMessage {
	return &WebsocketMessage{
		Action: TasksProcess,
		Data:   data,
	}
}

// TODO 怎么设计任务呢？
type TaskData struct {
	Id string
	// 是topic还是video
	Type string
	// 指定搜索关键词、视频BV
	Data []string
	// 指话题
	Topic string
	// 指定搜索翻页页数,默认为1
	Page int
}
