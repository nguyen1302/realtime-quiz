package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// JSONB definition for GORM
type JSONB []string

func (j JSONB) Value() (driver.Value, error) {
	return json.Marshal(j)
}

func (j *JSONB) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, j)
}

type Question struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	QuizID        uuid.UUID `gorm:"type:uuid;not null;index" json:"quiz_id"`
	Text          string    `gorm:"not null" json:"text"`
	Options       JSONB     `gorm:"type:jsonb" json:"options"`
	CorrectAnswer string    `gorm:"not null" json:"correct_answer"`
	TimeLimit     int       `gorm:"default:30" json:"time_limit"`
	Points        int       `gorm:"default:100" json:"points"`
	Order         int       `gorm:"column:item_order;default:0" json:"order"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (q *Question) BeforeCreate(tx *gorm.DB) (err error) {
	if q.ID == uuid.Nil {
		q.ID = uuid.New()
	}
	return
}
