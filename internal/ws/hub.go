package ws

import (
	"sync"
)

type Hub struct {
	mu      sync.RWMutex
	clients map[*Conn]struct{}
}

func NewHub() *Hub {
	return &Hub{clients: make(map[*Conn]struct{})}
}

func (h *Hub) Add(c *Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.clients[c] = struct{}{}
}

func (h *Hub) Remove(c *Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.clients, c)
	_ = c.conn.Close()
}

func (h *Hub) Broadcast(msg string) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	for c := range h.clients {
		_ = c.WriteText(msg)
	}
}
