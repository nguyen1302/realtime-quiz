package models

import "github.com/google/uuid"

type LeaderboardEntry struct {
	UserID   uuid.UUID `json:"user_id"`
	Username string    `json:"username,omitempty"` // Enriched later
	Score    float64   `json:"score"`
	Rank     int       `json:"rank"`
}
