package handler

import (
	"github.com/nguyen1302/realtime-quiz/internal/service"
)

// Handler is the interface for the handler manager
type Handler interface {
	Auth() AuthHandler
	Quiz() QuizHandler
}

// handlerImpl is the concrete implementation of Handler
type handlerImpl struct {
	auth AuthHandler
	quiz QuizHandler
}

// NewHandler creates a new instance of Handler
func NewHandler(svc service.Service) Handler {
	return &handlerImpl{
		auth: NewAuthHandler(svc.Auth()),
		quiz: NewQuizHandler(svc.Quiz()),
	}
}

func (h *handlerImpl) Auth() AuthHandler {
	return h.auth
}

func (h *handlerImpl) Quiz() QuizHandler {
	return h.quiz
}
