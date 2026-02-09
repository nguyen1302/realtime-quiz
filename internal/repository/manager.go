package repository

import (
	"gorm.io/gorm"
)

// Repository is the interface for the repository manager
type Repository interface {
	User() UserRepository
	Quiz() QuizRepository
	Question() QuestionRepository
}

// repositoryImpl is the concrete implementation of Repository
type repositoryImpl struct {
	user     UserRepository
	quiz     QuizRepository
	question QuestionRepository
}

// NewRepository creates a new instance of Repository
func NewRepository(db *gorm.DB) Repository {
	return &repositoryImpl{
		user:     NewUserRepository(db),
		quiz:     NewQuizRepository(db),
		question: NewQuestionRepository(db),
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
