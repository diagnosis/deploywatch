package events

import (
	"context"

	"github.com/google/uuid"
)

type Hub struct {
	clients         map[*Client]bool
	register        chan *Client
	unregister      chan *Client
	broadcast       chan Event
	broadcastToUser chan UserMessage
}

func NewHub() *Hub {
	return &Hub{
		clients:         make(map[*Client]bool),
		register:        make(chan *Client),
		unregister:      make(chan *Client),
		broadcast:       make(chan Event, 100),
		broadcastToUser: make(chan UserMessage, 100),
	}
}

type Event struct {
	Type string
	Data string
}

type UserMessage struct {
	UserID  uuid.UUID
	Message Event
}

type Client struct {
	UserID uuid.UUID  `json:"user_id"`
	Send   chan Event `json:"send"`
}

func (h *Hub) Run(ctx context.Context) {
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
		case msg := <-h.broadcast:
			for c := range h.clients {
				select {
				case c.Send <- Event{
					Type: msg.Type,
					Data: msg.Data,
				}:
				default:
				}
			}
		case userMsg := <-h.broadcastToUser:
			for c := range h.clients {
				if c.UserID == userMsg.UserID {
					select {
					case c.Send <- userMsg.Message:
					default:
					}
				}
			}
		}

	}
}

// public methods
func (h *Hub) Register(c *Client) {
	h.register <- c
}

func (h *Hub) Unregister(c *Client) {
	h.unregister <- c
}
func (h *Hub) Broadcast(msg Event) {
	h.broadcast <- msg
}

func (h *Hub) BroadcastToUser(userID uuid.UUID, msg Event) {
	h.broadcastToUser <- UserMessage{
		UserID:  userID,
		Message: msg,
	}
}
