package service

import (
	"github.com/nguyen1302/realtime-quiz/internal/config"
	"github.com/nguyen1302/realtime-quiz/internal/repository"
)

// Service is the interface for the service manager
type Service interface {
	Auth() AuthService
	Quiz() QuizService
	Realtime() RealtimeService
}

// serviceImpl is the concrete implementation of Service
type serviceImpl struct {
	auth     AuthService
	quiz     QuizService
	realtime RealtimeService
}

// NewService creates a new instance of Service
func NewService(repo repository.Repository, cfg *config.Config) Service {
	realtimeSvc := NewRealtimeService()
	return &serviceImpl{
		auth:     NewAuthService(repo.User(), cfg.JWT),
		quiz:     NewQuizService(repo.Quiz(), repo.Question(), repo.Leaderboard(), repo.Answer(), realtimeSvc),
		realtime: realtimeSvc,
	}
}

func (s *serviceImpl) Auth() AuthService {
	return s.auth
}

func (s *serviceImpl) Quiz() QuizService {
	return s.quiz
}

func (s *serviceImpl) Realtime() RealtimeService {
	return s.realtime
}
