package websock

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var Upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Hub struct {
	clients    []*Client
	nextId     int
	register   chan *Client
	unregister chan *Client
	mutex      *sync.Mutex
}

func NewHub() *Hub {
	return &Hub{
		clients:    make([]*Client, 0),
		nextId:     0,
		register:   make(chan *Client),
		unregister: make(chan *Client),
		mutex:      &sync.Mutex{},
	}
}

func (h *Hub) Register(client *Client) {
	h.register <- client
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.onConnect(client)
		case client := <-h.unregister:
			h.onDisconnect(client)
		}
	}
}

func (h *Hub) Broadcast(msg any, ignore *Client) {
	data, _ := json.Marshal(msg)
	for _, c := range h.clients {
		if c != ignore {
			c.outbound <- data
		}
	}
}

func (h *Hub) Send(msg any, client *Client) {
	data, _ := json.Marshal(msg)
	client.outbound <- data
}

func (h *Hub) onConnect(client *Client) {
	log.Println("client connected:", client.socket.RemoteAddr())

	h.mutex.Lock()
	defer h.mutex.Unlock()

	client.id = h.nextId
	h.nextId++
	h.clients = append(h.clients, client)
}

func (h *Hub) onDisconnect(client *Client) {
	log.Println("client disconnected:", client.socket.RemoteAddr())

	client.close()
	h.mutex.Lock()
	defer h.mutex.Unlock()

	var idx int
	for i := range h.clients {
		if h.clients[i].id == client.id {
			idx = i
			break
		}
	}

	copy(h.clients[idx:], h.clients[idx+1:])
	h.clients[len(h.clients)-1] = nil
	h.clients = h.clients[:len(h.clients)-1]
}
