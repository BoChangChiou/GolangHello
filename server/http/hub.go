package http

import (
	"encoding/json"
)

var wsHub = hub{
	connections: make(map[*connection]bool),
	disconnect:  make(chan *connection),
	broadcast:   make(chan []byte),
	newConnect:  make(chan *connection),
}

type hub struct {
	connections map[*connection]bool
	broadcast   chan []byte
	newConnect  chan *connection
	disconnect  chan *connection
}

func (h *hub) run() {
	for {
		select {
		case c := <-h.newConnect:
			h.connections[c] = true
			c.data.Ip = c.ws.RemoteAddr().String()
			c.data.Type = "handshake"
			c.data.UserList = user_list
			data_b, _ := json.Marshal(c.data)
			c.sc <- data_b
		case c := <-h.disconnect:
			if _, ok := h.connections[c]; ok {
				delete(h.connections, c)
				close(c.sc)
			}
		case data := <-h.broadcast:
			for c := range h.connections {
				select {
				case c.sc <- data:
				default:
					delete(h.connections, c)
					close(c.sc)
				}
			}
		}
	}
}
