package api_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/nguyen1302/realtime-quiz/internal/bootstrap"
	"github.com/nguyen1302/realtime-quiz/internal/config"
	"github.com/nguyen1302/realtime-quiz/internal/models"
	"github.com/nguyen1302/realtime-quiz/internal/realtime"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestRealTimeQuizFlow(t *testing.T) {
	// 1. Setup Environment
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	require.NoError(t, err)

	// Run migrations
	err = db.AutoMigrate(&models.User{}, &models.Quiz{}, &models.Question{}, &models.Answer{})
	require.NoError(t, err)

	// Setup Redis (Mock or Real? Using miniredis is better but for now assuming local redis or skip)
	// For integration test, we might need real redis.
	// Let's assume a real redis at localhost:6379 for this manual verification based on previous context
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		t.Skip("Redis not available, skipping integration test")
	}
	defer rdb.FlushAll(context.Background())

	cfg := &config.Config{
		JWT: config.JWTConfig{
			Secret:      "test-secret",
			ExpiryHours: 1,
		},
	}

	app := bootstrap.NewRouter(db, rdb, cfg)
	server := httptest.NewServer(app.Engine())
	defer server.Close()

	// 2. Create User & Quiz
	// Register User
	request(t, server, "POST", "/api/v1/auth/register", `{"username":"testuser","password":"password","email":"test@example.com"}`)

	// Login User
	loginResp := request(t, server, "POST", "/api/v1/auth/login", `{"email":"test@example.com","password":"password"}`)
	token := getToken(t, loginResp)

	// Create Quiz
	quizResp := requestWithAuth(t, server, "POST", "/api/v1/quizzes", `{"title":"RT Quiz","description":"Test"}`, token)
	var quizObj struct {
		Data struct {
			ID string `json:"id"`
		} `json:"data"`
	}
	json.Unmarshal(quizResp, &quizObj)
	quizID := quizObj.Data.ID

	// Add Question
	questionResp := requestWithAuth(t, server, "POST", fmt.Sprintf("/api/v1/quizzes/%s/questions", quizID), `{"text":"Q1","options":["A","B"],"correct_answer":"A","point":100,"time_limit":10}`, token)
	var questionObj struct {
		Data struct {
			ID string `json:"id"`
		} `json:"data"`
	}
	json.Unmarshal(questionResp, &questionObj)
	questionID := questionObj.Data.ID

	// 3. Connect WebSocket
	// Convert http url to ws url
	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/api/v1/ws"
	// Append token? Our ws handler checks "userID" from context, which implies middleware.
	// But /ws route might not be protected by auth middleware in router.go?
	// Checking router.go: 	api.GET("/ws", r.handlers.Realtime().HandleConnection) is OUTSIDE protected group.
	// But it checks c.Get("userID").
	// So we need to pass auth somehow? usually via header or query param.
	// Validating router.go again...
	// Ah, it's public. But it tries to get userID. If not logged in, it's anonymous?
	// Wait, without userID, we can still join quiz?
	// Let's try connecting.

	dialer := websocket.Dialer{}
	conn, _, err := dialer.Dial(wsURL, nil)
	require.NoError(t, err)
	defer conn.Close()

	// 4. Join Quiz Room
	joinMsg := map[string]interface{}{
		"type": "join_quiz",
		"payload": map[string]string{
			"quiz_id": quizID,
		},
	}
	err = conn.WriteJSON(joinMsg)
	require.NoError(t, err)

	// Allow time for subscription
	time.Sleep(100 * time.Millisecond)

	// 5. Submit Answer (triggering broadcast)
	submitBody := fmt.Sprintf(`{"question_id":"%s","answer":"A"}`, questionID)
	requestWithAuth(t, server, "POST", fmt.Sprintf("/api/v1/quizzes/%s/submit", quizID), submitBody, token)

	// 6. Verify Broadcast Message
	conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	_, msg, err := conn.ReadMessage()
	require.NoError(t, err)

	t.Logf("Received WS Message: %s", string(msg))

	var wsMsg realtime.WSMessage
	err = json.Unmarshal(msg, &wsMsg)
	require.NoError(t, err)

	assert.Equal(t, "leaderboard_update", wsMsg.Type)
}

// Helpers moved to helper_test.go
