package realtime

import (
	"log"
	"sync"
)

// Hub maintains the set of active clients and broadcasts messages to clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	// Inbound messages from the clients.
	broadcast chan *WSMessage

	// Map QuizID to Clients for targeted messaging
	quizClients map[string]map[*Client]bool

	// Optional: Map UserID to Clients for targeted messaging
	userClients map[string]map[*Client]bool
	mu          sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		broadcast:   make(chan *WSMessage),
		register:    make(chan *Client),
		unregister:  make(chan *Client),
		clients:     make(map[*Client]bool),
		userClients: make(map[string]map[*Client]bool),
		quizClients: make(map[string]map[*Client]bool),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			if client.userID != "" {
				if h.userClients[client.userID] == nil {
					h.userClients[client.userID] = make(map[*Client]bool)
				}
				h.userClients[client.userID][client] = true
			}
			h.mu.Unlock()
			log.Printf("Client registered: %s", client.userID)

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				if client.userID != "" && h.userClients[client.userID] != nil {
					delete(h.userClients[client.userID], client)
					if len(h.userClients[client.userID]) == 0 {
						delete(h.userClients, client.userID)
					}
				}
				// Remove from all quizzes
				for quizID, clients := range h.quizClients {
					if _, ok := clients[client]; ok {
						delete(clients, client)
						if len(clients) == 0 {
							delete(h.quizClients, quizID)
						}
					}
				}

				close(client.send)
				log.Printf("Client unregistered: %s", client.userID)
			}
			h.mu.Unlock()

		case message := <-h.broadcast:
			// Logic handled by helper methods now, but this channel can still be used for global broadcast
			h.broadcastMessage(message)
		}
	}
}

func (h *Hub) broadcastMessage(message *WSMessage) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for client := range h.clients {
		select {
		case client.send <- message:
		default:
			close(client.send)
			delete(h.clients, client)
		}
	}
}

func (h *Hub) BroadcastToUser(userID string, message *WSMessage) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if clients, ok := h.userClients[userID]; ok {
		for client := range clients {
			select {
			case client.send <- message:
			default:
				// Handle slow client
			}
		}
	}
}

func (h *Hub) SubscribeToQuiz(client *Client, quizID string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.quizClients[quizID] == nil {
		h.quizClients[quizID] = make(map[*Client]bool)
	}
	h.quizClients[quizID][client] = true
	log.Printf("Client %s subscribed to quiz %s", client.userID, quizID)
}

func (h *Hub) UnsubscribeFromQuiz(client *Client, quizID string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if clients, ok := h.quizClients[quizID]; ok {
		delete(clients, client)
		if len(clients) == 0 {
			delete(h.quizClients, quizID)
		}
	}
}

func (h *Hub) BroadcastToQuiz(quizID string, message *WSMessage) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if clients, ok := h.quizClients[quizID]; ok {
		for client := range clients {
			select {
			case client.send <- message:
			default:
				// Handle slow client
			}
		}
	}
}
