package handlers

import (
	"net/http"
	"vms_go/internal/ws"
)

type WSHandler struct {
	Hub *ws.Hub
}

func (h *WSHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c, err := ws.Upgrade(w, r)
	if err != nil {
		http.Error(w, "upgrade failed", http.StatusBadRequest)
		return
	}
	h.Hub.Add(c)
}
