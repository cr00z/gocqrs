package websock

import "github.com/gorilla/websocket"

type Client struct {
	hub      *Hub
	id       int
	socket   *websocket.Conn
	outbound chan []byte
}

func NewClient(hub *Hub, socket *websocket.Conn) *Client {
	return &Client{
		hub:      hub,
		socket:   socket,
		outbound: make(chan []byte),
	}
}

func (c *Client) Write() {
	for {
		data, ok := <-c.outbound
		if !ok {
			_ = c.socket.WriteMessage(websocket.CloseMessage, []byte{})
			return
		}
		_ = c.socket.WriteMessage(websocket.TextMessage, data)
	}
}

func (c *Client) close() {
	c.socket.Close()
	close(c.outbound)
}
