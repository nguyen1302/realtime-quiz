package handler

import (
	"github.com/nguyen1302/realtime-quiz/internal/service"
)

type Handlers struct {
	Auth *AuthHandler
	// Add other handlers here
}

func NewHandlers(services *service.Services) *Handlers {
	return &Handlers{
		Auth: NewAuthHandler(services.Auth),
	}
}
