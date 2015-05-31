package room

import (
	"encoding/json"

	"gopkg.in/igm/sockjs-go.v2/sockjs"
)

type Room struct {
	// Room name - Useful for debugging.
	name string

	// Registered connections.
	connections map[sockjs.Session]bool

	// Inbound messages from the connections.
	broadcast chan string

	// Register requests from the connections.
	register chan sockjs.Session

	// Unregister requests from connections.
	unregister chan sockjs.Session
}

func New(name string) Room {
	return Room{
		name:        name,
		broadcast:   make(chan string),
		register:    make(chan sockjs.Session),
		unregister:  make(chan sockjs.Session),
		connections: make(map[sockjs.Session]bool),
	}
}

func (h *Room) Run() {
	for {
		select {
		case c := <-h.register:
			h.connections[c] = true
		case c := <-h.unregister:
			delete(h.connections, c)
		case m := <-h.broadcast:
			for c := range h.connections {
				go c.Send(m)
			}
		}
	}
}

func (h *Room) Join(ses sockjs.Session) {
	h.register <- ses
}

func (h *Room) Leave(ses sockjs.Session) {
	h.unregister <- ses
}

func (h *Room) Broadcast(m interface{}) error {
	b, err := json.Marshal(m)
	if err != nil {
		return err
	}
	h.broadcast <- string(b)
	return nil
}
