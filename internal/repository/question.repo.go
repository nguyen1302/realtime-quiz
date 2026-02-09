package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/nguyen1302/realtime-quiz/internal/models"
	"gorm.io/gorm"
)

type QuestionRepository interface {
	Create(ctx context.Context, question *models.Question) error
	GetByQuizID(ctx context.Context, quizID uuid.UUID) ([]models.Question, error)
}

type questionRepository struct {
	db *gorm.DB
}

func NewQuestionRepository(db *gorm.DB) QuestionRepository {
	return &questionRepository{db: db}
}

func (r *questionRepository) Create(ctx context.Context, question *models.Question) error {
	return r.db.WithContext(ctx).Create(question).Error
}

func (r *questionRepository) GetByQuizID(ctx context.Context, quizID uuid.UUID) ([]models.Question, error) {
	var questions []models.Question
	if err := r.db.WithContext(ctx).Where("quiz_id = ?", quizID).Order("item_order asc").Find(&questions).Error; err != nil {
		return nil, err
	}
	return questions, nil
}
