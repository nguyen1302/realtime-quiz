package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/nguyen1302/realtime-quiz/internal/models"
	"gorm.io/gorm"
)

type QuizRepository interface {
	Create(ctx context.Context, quiz *models.Quiz) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Quiz, error)
	GetByCode(ctx context.Context, code string) (*models.Quiz, error)
}

type quizRepository struct {
	db *gorm.DB
}

func NewQuizRepository(db *gorm.DB) QuizRepository {
	return &quizRepository{db: db}
}

func (r *quizRepository) Create(ctx context.Context, quiz *models.Quiz) error {
	return r.db.WithContext(ctx).Create(quiz).Error
}

func (r *quizRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Quiz, error) {
	var quiz models.Quiz
	if err := r.db.WithContext(ctx).Preload("Questions").Where("id = ?", id).First(&quiz).Error; err != nil {
		return nil, err
	}
	return &quiz, nil
}

func (r *quizRepository) GetByCode(ctx context.Context, code string) (*models.Quiz, error) {
	var quiz models.Quiz
	if err := r.db.WithContext(ctx).Where("code = ?", code).First(&quiz).Error; err != nil {
		return nil, err
	}
	return &quiz, nil
}
