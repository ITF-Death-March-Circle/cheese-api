package main

import (
	"log"
	"main/redis"

	"github.com/pkg/errors"
)

var errorSocketResponse = []byte(`{"action":"ERROR_MESSAGE","status":"NG","error": true}`)

var h = hub{
	broadcast:  make(chan message),
	notify:     make(chan message),
	register:   make(chan subscription),
	unregister: make(chan subscription),
	rooms:      make(map[string]map[*connection]bool),
}

//

func (h *hub) run() {
	for {
		select {
		case s := <-h.register:
			err := redis.AddValue(COUNT_USER)
			if err != nil {
				log.Fatalln(errors.WithStack(err))
			}
			connections := h.rooms[s.room]
			if connections == nil {
				connections = make(map[*connection]bool)
				h.rooms[s.room] = connections
			}
			h.rooms[s.room][s.conn] = true
		case s := <-h.unregister:
			_, err := redis.DeclValue(COUNT_USER)
			if err != nil {
				log.Fatalln(errors.WithStack(err))
			}
			connections := h.rooms[s.room]
			if connections != nil {
				if _, ok := connections[s.conn]; ok {
					delete(connections, s.conn)
					close(s.conn.send)
					if len(connections) == 0 {
						delete(h.rooms, s.room)
					}
				}
			}
		case m := <-h.broadcast:
			connections := h.rooms[m.room]
			for c := range connections {
				select {
				case c.send <- m.data:
				default:
					close(c.send)
					delete(connections, c)
					if len(connections) == 0 {
						delete(h.rooms, m.room)
					}
				}
			}
		case p := <-h.notify:
			connections := h.rooms["maid"]
			for c := range connections {
				select {
				case c.send <- p.data:
				default:
					close(c.send)
					delete(connections, c)
					if len(connections) == 0 {
						delete(h.rooms, p.room)
					}
				}
			}
		}
	}
}
