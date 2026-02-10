package handler

import (
	"github.com/nguyen1302/realtime-quiz/internal/service"
)

// Handler is the interface for the handler manager
type Handler interface {
	Auth() AuthHandler
	Quiz() QuizHandler
	Realtime() WebSocketHandler
}

// handlerImpl is the concrete implementation of Handler
type handlerImpl struct {
	auth     AuthHandler
	quiz     QuizHandler
	realtime WebSocketHandler
}

// NewHandler creates a new instance of Handler
func NewHandler(svc service.Service) Handler {
	return &handlerImpl{
		auth:     NewAuthHandler(svc.Auth()),
		quiz:     NewQuizHandler(svc.Quiz()),
		realtime: NewWebSocketHandler(svc.Realtime()),
	}
}

func (h *handlerImpl) Auth() AuthHandler {
	return h.auth
}

func (h *handlerImpl) Quiz() QuizHandler {
	return h.quiz
}

func (h *handlerImpl) Realtime() WebSocketHandler {
	return h.realtime
}
