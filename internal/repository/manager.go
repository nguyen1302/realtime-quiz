package repository

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// Repository is the interface for the repository manager
type Repository interface {
	User() UserRepository
	Quiz() QuizRepository
	Question() QuestionRepository
	Leaderboard() LeaderboardRepository
	Answer() AnswerRepository
}

// repositoryImpl is the concrete implementation of Repository
type repositoryImpl struct {
	user        UserRepository
	quiz        QuizRepository
	question    QuestionRepository
	leaderboard LeaderboardRepository
	answer      AnswerRepository
}

// NewRepository creates a new instance of Repository
func NewRepository(db *gorm.DB, rdb *redis.Client) Repository {
	return &repositoryImpl{
		user:        NewUserRepository(db),
		quiz:        NewQuizRepository(db),
		question:    NewQuestionRepository(db),
		leaderboard: NewLeaderboardRepository(rdb),
		answer:      NewAnswerRepository(db),
	}
}

func (r *repositoryImpl) User() UserRepository {
	return r.user
}

func (r *repositoryImpl) Quiz() QuizRepository {
	return r.quiz
}

func (r *repositoryImpl) Question() QuestionRepository {
	return r.question
}

func (r *repositoryImpl) Leaderboard() LeaderboardRepository {
	return r.leaderboard
}

func (r *repositoryImpl) Answer() AnswerRepository {
	return r.answer
}
