package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/nguyen1302/realtime-quiz/internal/models"
	"gorm.io/gorm"
)

type AnswerRepository interface {
	Create(ctx context.Context, answer *models.Answer) error
	HasAnswered(ctx context.Context, quizID, questionID, userID uuid.UUID) (bool, error)
}

type answerRepository struct {
	db *gorm.DB
}

func NewAnswerRepository(db *gorm.DB) AnswerRepository {
	return &answerRepository{db: db}
}

func (r *answerRepository) Create(ctx context.Context, answer *models.Answer) error {
	return r.db.WithContext(ctx).Create(answer).Error
}

func (r *answerRepository) HasAnswered(ctx context.Context, quizID, questionID, userID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Answer{}).
		Where("quiz_id = ? AND question_id = ? AND user_id = ?", quizID, questionID, userID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
