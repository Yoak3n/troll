package ws

import (
	"encoding/json"
	"log"

	"github.com/gofiber/contrib/websocket"
)

type WebSocketHub struct {
	// Registered clients.
	clients map[string]*websocket.Conn

	// Inbound messages from the clients.
	received chan map[*websocket.Conn][]byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *websocket.Conn

	// broadcast message to all clients
	broadcast chan []byte
	isRuning  bool
}

type Client struct {
	id   string
	conn *websocket.Conn
}

func NewWebSocketHub() *WebSocketHub {
	return &WebSocketHub{
		broadcast:  make(chan []byte),
		received:   make(chan map[*websocket.Conn][]byte),
		register:   make(chan *Client),
		unregister: make(chan *websocket.Conn),
		clients:    make(map[string]*websocket.Conn),
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
	h.isRuning = true
	for {
		select {
		case client := <-h.register:
			h.clients[client.id] = client.conn
		case client := <-h.unregister:
			for id, v := range h.clients {
				if client == v {
					delete(h.clients, id)
					client.Conn.Close()
					log.Printf("Client disconnected with %d clients\n", len(h.clients))
				}
			}
		case message := <-h.broadcast:
			for k, client := range h.clients {
				if err := client.WriteMessage(websocket.TextMessage, message); err != nil {
					client.Conn.Close()
					delete(h.clients, k)
				}
			}
		}
	}
}

func (h *WebSocketHub) Register(id string, conn *websocket.Conn) {
	if k, ok := h.clients[id]; ok {
		k.Conn.Close()
		delete(h.clients, id)
	}
	client := &Client{
		id:   id,
		conn: conn,
	}
	h.register <- client
	log.Printf("New connection with %d clients\n", len(h.clients))
	sendLog(client.conn, "Connected successfully")
	h.listen(conn)
}

func (h *WebSocketHub) Broadcast(message []byte) {
	h.broadcast <- message
}

func (h *WebSocketHub) listen(client *websocket.Conn) {
	defer func() {
		h.unregister <- client
	}()
	for {
		t, msg, err := client.Conn.ReadMessage()
		log.Println("T", t, string(msg))
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
			client.Conn.WriteJSON(pongMessage)
		}

		h.received <- map[*websocket.Conn][]byte{client: msg}
	}
}

func (h *WebSocketHub) Handle() {
	for msg := range h.received {
		for client, msg := range msg {
			// 最强人工智能
			if err := client.WriteMessage(websocket.TextMessage, msg); err != nil {
				client.Conn.Close()
			}
		}
	}
}

func (h *WebSocketHub) broadcastLog(content string) {
	logItem := NewLogDataToMessage(content)
	j, err := json.Marshal(logItem)
	if err != nil {
		return
	}
	h.Broadcast(j)
}

func sendLog(conn *websocket.Conn, content string) {
	logItem := NewLogDataToMessage(content)
	conn.WriteJSON(logItem)
}
