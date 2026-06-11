package ws

import (
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func NewHub() *Hub {
	return &Hub{
		Clients:    make(map[uint][]*Client),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Send:       make(chan OutboundMsg, 256),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.Clients[client.UserID] = append(h.Clients[client.UserID], client)
		case client := <-h.Unregister:
			list := h.Clients[client.UserID]

			for i, c := range list {
				if c == client {
					list = append(list[:i], list[i+1:]...)
					h.Clients[client.UserID] = list

					if len(h.Clients[client.UserID]) == 0 {
						delete(h.Clients, client.UserID)
					}

					close(client.Send)
					break
				}
			}
		case msg := <-h.Send:
			list := h.Clients[msg.RecipientID]

			for _, c := range list {
				select {
				case c.Send <- msg.Payload:
				default:

				}
			}
		}
	}
}
