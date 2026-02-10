package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Answer struct {
	ID         uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	QuizID     uuid.UUID      `gorm:"type:uuid;not null;index" json:"quiz_id"`
	QuestionID uuid.UUID      `gorm:"type:uuid;not null;index" json:"question_id"`
	UserID     uuid.UUID      `gorm:"type:uuid;not null;index" json:"user_id"`
	Answer     string         `gorm:"type:text;not null" json:"answer"`
	IsCorrect  bool           `gorm:"default:false" json:"is_correct"`
	Points     int            `gorm:"default:0" json:"points"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

func (a *Answer) BeforeCreate(tx *gorm.DB) (err error) {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	return
}
