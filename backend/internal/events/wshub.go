package events

import (
	"context"

	"github.com/google/uuid"
)

type WSHub struct {
	clients         map[*WSClient]bool
	register        chan *WSClient
	unregister      chan *WSClient
	broadcastToUser chan WSUserMessage
}

type WSClient struct {
	UserID uuid.UUID
	Send   chan Event
}

type WSUserMessage struct {
	UserID  uuid.UUID
	Message Event
}

func NewWSHub() *WSHub {
	return &WSHub{
		clients:         make(map[*WSClient]bool),
		register:        make(chan *WSClient),
		unregister:      make(chan *WSClient),
		broadcastToUser: make(chan WSUserMessage, 100),
	}
}

func (h *WSHub) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case c := <-h.register:
			h.clients[c] = true
		case c := <-h.unregister:
			if _, ok := h.clients[c]; ok {
				delete(h.clients, c)
				close(c.Send)
			}
		case msg := <-h.broadcastToUser:
			for c := range h.clients {
				if c.UserID == msg.UserID {
					select {
					case c.Send <- msg.Message:
					default:
					}
				}
			}
		}
	}
}

func (h *WSHub) Register(c *WSClient) {
	h.register <- c
}

func (h *WSHub) Unregister(c *WSClient) {
	h.unregister <- c
}

func (h *WSHub) BroadcastToUser(userID uuid.UUID, msg Event) {
	h.broadcastToUser <- WSUserMessage{
		UserID:  userID,
		Message: msg,
	}
}
