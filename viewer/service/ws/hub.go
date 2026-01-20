package ws

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/gofiber/contrib/websocket"
)

type WebSocketHub struct {
	// Registered clients.
	clients map[string]*Client

	// Inbound messages from the clients.
	received chan map[*websocket.Conn][]byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan string

	// broadcast message to all clients
	broadcast chan []byte
	Tasks     chan TaskData
	isRuning  bool
}

type Client struct {
	id   string
	conn *websocket.Conn
	mu   sync.RWMutex
	last int64
}

func NewWebSocketHub() *WebSocketHub {
	return &WebSocketHub{
		broadcast:  make(chan []byte),
		received:   make(chan map[*websocket.Conn][]byte),
		register:   make(chan *Client),
		unregister: make(chan string),
		clients:    make(map[string]*Client),
		Tasks:      make(chan TaskData, 10),
		isRuning:   false,
	}
}

var Hub *WebSocketHub

func InitWebsocketHub() {
	Hub = NewWebSocketHub()
	go Hub.Run()
}

func (h *WebSocketHub) Run() {
	if h.isRuning {
		return
	}
	defer func() {
		h.isRuning = false
	}()
	h.isRuning = true
	for {
		select {
		case client := <-h.register:
			h.clients[client.id] = client
			log.Printf("New connection with %d clients\n", len(h.clients))
		case id := <-h.unregister:
			if client, ok := h.clients[id]; ok {
				close := WebsocketMessage{
					Action: CloseMessage,
					Data:   "Close",
				}
				client.mu.Lock()
				client.conn.Conn.WriteJSON(close)
				delete(h.clients, id)
				client.conn.Conn.Close()
				client.mu.Unlock()
				log.Printf("Client disconnected with %d clients\n", len(h.clients))
			}
		case message := <-h.broadcast:
			for k, client := range h.clients {
				if err := client.conn.WriteMessage(websocket.TextMessage, message); err != nil {
					client.conn.Conn.Close()
					delete(h.clients, k)
				}
			}
		case <-h.received:

		}
	}
}

func (h *WebSocketHub) Register(id string, conn *websocket.Conn) {
	if k, ok := h.clients[id]; ok {
		k.conn.Conn.Close()
		delete(h.clients, id)
	}
	client := &Client{
		id:   id,
		conn: conn,
		mu:   sync.RWMutex{},
		last: time.Now().Unix(),
	}
	h.register <- client
	sendLog(client.conn, "Connected successfully")
	go h.healthCheck(client)
	h.listen(id, conn)
}

func (h *WebSocketHub) Broadcast(message []byte) {
	h.broadcast <- message
}

func (h *WebSocketHub) healthCheck(client *Client) {
	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		client.mu.RLock()
		last := client.last
		now := time.Now().Unix()
		if (now - last) >= 180 {
			h.unregister <- client.id
			return
		}
		pingMessage := WebsocketMessage{
			Action: PingMessage,
			Data:   "ping",
		}
		client.conn.Conn.WriteJSON(pingMessage)
		client.mu.RUnlock()
	}
}

func (h *WebSocketHub) listen(id string, conn *websocket.Conn) {
	defer func() {
		h.unregister <- id
	}()
	for {
		t, msg, err := conn.Conn.ReadMessage()
		if err != nil || t == -1 {
			break
		}
		messageData := &WebsocketMessage{}
		err = json.Unmarshal(msg, messageData)
		if err != nil {
			break
		}
		switch messageData.Action {
		case CloseMessage:
			return
		case PingMessage:
			pongMessage := WebsocketMessage{
				Action: PongMessage,
				Data:   "pong",
			}
			conn.Conn.WriteJSON(pongMessage)
		case AddTaskMessage:
			// 无法直接断言为任务数据
			buf, _ := json.Marshal(messageData.Data)
			taskData := TaskData{}
			err := json.Unmarshal(buf, &taskData)
			if err != nil {
				log.Println("task data unmarshal err:", err)
			}
			h.Tasks <- taskData
		}
		client, ok := h.clients[id]
		if ok {
			client.mu.Lock()
			client.last = time.Now().Unix()
			client.mu.Unlock()
			h.clients[id] = client
		}
		h.received <- map[*websocket.Conn][]byte{conn: msg}
	}
}

// func (h *WebSocketHub) Handle() {
// 	for msg := range  {
// 		for client, msg := range msg {
// 			// 最强人工智能
// 			if err := client.WriteMessage(websocket.TextMessage, msg); err != nil {
// 				client.Conn.Close()
// 			}
// 		}
// 	}
// }

// func (h *WebSocketHub) broadcastLog(content string) {
// 	logItem := NewLogDataToMessage(content)
// 	j, err := json.Marshal(logItem)
// 	if err != nil {
// 		return
// 	}
// 	h.Broadcast(j)
// }

func sendLog(conn *websocket.Conn, content string) {
	logItem := NewLogDataToMessage(content)
	conn.WriteJSON(logItem)
}
