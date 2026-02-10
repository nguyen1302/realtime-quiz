package service

import (
	"github.com/nguyen1302/realtime-quiz/internal/realtime"
)

type RealtimeService interface {
	BroadcastToUser(userID string, messageType string, payload interface{})
	BroadcastToQuiz(quizID string, messageType string, payload interface{})
	BroadcastToAll(messageType string, payload interface{})
	GetManager() *realtime.Manager
}

type realtimeService struct {
	manager *realtime.Manager
}

func NewRealtimeService() RealtimeService {
	return &realtimeService{
		manager: realtime.NewManager(),
	}
}

func (s *realtimeService) BroadcastToUser(userID string, messageType string, payload interface{}) {
	msg := &realtime.WSMessage{
		Type:    messageType,
		Payload: payload,
		UserID:  userID,
	}
	s.manager.SendToUser(userID, msg)
}

func (s *realtimeService) BroadcastToQuiz(quizID string, messageType string, payload interface{}) {
	msg := &realtime.WSMessage{
		Type:    messageType,
		Payload: payload,
		RoomID:  quizID,
	}
	s.manager.BroadcastToQuiz(quizID, msg)
}

func (s *realtimeService) BroadcastToAll(messageType string, payload interface{}) {
	msg := &realtime.WSMessage{
		Type:    messageType,
		Payload: payload,
	}
	s.manager.Broadcast(msg)
}

func (s *realtimeService) GetManager() *realtime.Manager {
	return s.manager
}
