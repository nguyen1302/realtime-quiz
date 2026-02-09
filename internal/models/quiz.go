package models

import (
	"time"

	"github.com/google/uuid"
)

type QuizStatus string

const (
	QuizStatusDraft    QuizStatus = "DRAFT"
	QuizStatusActive   QuizStatus = "ACTIVE"
	QuizStatusFinished QuizStatus = "FINISHED"
)

type Quiz struct {
	ID          uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Title       string     `gorm:"not null" json:"title"`
	Description string     `json:"description"`
	Code        string     `gorm:"uniqueIndex;not null" json:"code"`
	Status      QuizStatus `gorm:"type:varchar(20);default:'DRAFT'" json:"status"`
	OwnerID     uuid.UUID  `gorm:"type:uuid;not null;index" json:"owner_id"`
	Owner       User       `gorm:"foreignKey:OwnerID" json:"owner,omitempty"`
	Questions   []Question `gorm:"foreignKey:QuizID" json:"questions,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}
