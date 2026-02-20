package ws

import (
	"encoding/json"
	"sync"

	"github.com/gorilla/websocket"
)

type Client struct {
	conn    *websocket.Conn
	address string
	send    chan []byte
}

type Hub struct {
	mu      sync.RWMutex
	clients map[string][]*Client
}

func NewHub() *Hub {
	return &Hub{clients: make(map[string][]*Client)}
}

func (h *Hub) Register(addr string, conn *websocket.Conn) *Client {
	c := &Client{conn: conn, address: addr, send: make(chan []byte, 32)}
	h.mu.Lock()
	h.clients[addr] = append(h.clients[addr], c)
	h.mu.Unlock()
	go c.writePump()
	return c
}

func (h *Hub) Unregister(c *Client) {
	h.mu.Lock()
	clients := h.clients[c.address]
	for i, cl := range clients {
		if cl == c {
			h.clients[c.address] = append(clients[:i], clients[i+1:]...)
			break
		}
	}
	h.mu.Unlock()
	close(c.send)
}

func (h *Hub) Broadcast(addr string, v any) {
	data, _ := json.Marshal(v)
	h.mu.RLock()
	for _, c := range h.clients[addr] {
		select {
		case c.send <- data:
		default:
		}
	}
	h.mu.RUnlock()
}

func (c *Client) writePump() {
	for msg := range c.send {
		if err := c.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			break
		}
	}
	c.conn.Close()
}
