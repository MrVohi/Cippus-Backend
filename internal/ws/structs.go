package ws

import "github.com/gorilla/websocket"

type Client struct {
	Hub    *Hub
	UserID uint
	Conn   *websocket.Conn
	Send   chan []byte
}

type OutboundMsg struct {
	RecipientID uint
	Payload     []byte
}

type Hub struct {
	Clients    map[uint][]*Client
	Register   chan *Client
	Unregister chan *Client
	Send       chan OutboundMsg
}
