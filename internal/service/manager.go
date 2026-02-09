package service

import (
	"github.com/nguyen1302/realtime-quiz/internal/config"
	"github.com/nguyen1302/realtime-quiz/internal/repository"
)

type Services struct {
	Auth AuthService
	// Add other services here
}

func NewServices(repos *repository.Repositories, cfg *config.Config) *Services {
	return &Services{
		Auth: NewAuthService(repos.User, cfg.JWT),
	}
}
