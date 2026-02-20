package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"sink.io/m/src/store"
	"sink.io/m/src/ws"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type Handler struct {
	Store *store.Store
	Hub   *ws.Hub
}

func (h *Handler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/api/mailbox/{address}", h.getEmails).Methods("GET")
	r.HandleFunc("/api/mailbox/{address}", h.deleteMailbox).Methods("DELETE")
	r.HandleFunc("/ws/{address}", h.handleWS)
}

func (h *Handler) getEmails(w http.ResponseWriter, r *http.Request) {
	addr := strings.ToLower(mux.Vars(r)["address"])
	emails := h.Store.GetByAddress(addr)
	if emails == nil {
		emails = []store.Email{}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(emails)
}

func (h *Handler) deleteMailbox(w http.ResponseWriter, r *http.Request) {
	addr := strings.ToLower(mux.Vars(r)["address"])
	h.Store.DeleteByAddress(addr)
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) handleWS(w http.ResponseWriter, r *http.Request) {
	addr := strings.ToLower(mux.Vars(r)["address"])
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	client := h.Hub.Register(addr, conn)
	defer h.Hub.Unregister(client)
	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			break
		}
	}
}
