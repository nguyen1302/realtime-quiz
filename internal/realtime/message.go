package realtime

import (
	"encoding/json"
)

// Event types
const (
	EventError        = "error"
	EventQuizState    = "quiz_state"
	EventUserJoined   = "user_joined"
	EventUserLeft     = "user_left"
	EventChatMessage  = "chat_message"
	EventLeaderboard  = "leaderboard_update"
	EventQuestion     = "new_question"
	EventAnswerResult = "answer_result"
)

// Message represents a WebSocket message
type Message struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload,omitempty"`
}

// WSMessage is the raw message passed around in the Hub
type WSMessage struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload,omitempty"`
	RoomID  string      `json:"-"` // Optional: for room-based broadcasting
	UserID  string      `json:"-"` // Optional: for direct messaging
}
