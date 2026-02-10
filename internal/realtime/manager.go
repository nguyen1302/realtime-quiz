package realtime

// Manager is the facade for the realtime package
type Manager struct {
	Hub *Hub
}

func NewManager() *Manager {
	hub := NewHub()
	go hub.Run()
	return &Manager{
		Hub: hub,
	}
}

func (m *Manager) Broadcast(message *WSMessage) {
	m.Hub.broadcast <- message
}

func (m *Manager) SendToUser(userID string, message *WSMessage) {
	m.Hub.BroadcastToUser(userID, message)
}

func (m *Manager) BroadcastToQuiz(quizID string, message *WSMessage) {
	m.Hub.BroadcastToQuiz(quizID, message)
}
