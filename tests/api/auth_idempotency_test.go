package api_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/nguyen1302/realtime-quiz/internal/realtime"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWSAuthAndIdempotency(t *testing.T) {
	// 1. Setup Environment
	_, _, server := setupTest(t)

	// 2. Create User & Quiz
	request(t, server, "POST", "/api/v1/auth/register", `{"username":"tester","password":"password","email":"tester@example.com"}`)
	loginResp := request(t, server, "POST", "/api/v1/auth/login", `{"email":"tester@example.com","password":"password"}`)
	token := getToken(t, loginResp)

	// Create Quiz
	quizResp := requestWithAuth(t, server, "POST", "/api/v1/quizzes", `{"title":"Auth Quiz","description":"Test"}`, token)
	var quizObj struct {
		Data struct {
			ID string `json:"id"`
		} `json:"data"`
	}
	json.Unmarshal(quizResp, &quizObj)
	quizID := quizObj.Data.ID

	// Add Question
	questionResp := requestWithAuth(t, server, "POST", fmt.Sprintf("/api/v1/quizzes/%s/questions", quizID), `{"text":"Q1","options":["A", "B"],"correct_answer":"A","points":100,"time_limit":10}`, token)
	t.Logf("Add Question Response: %s", string(questionResp))
	var questionObj struct {
		Data struct {
			ID string `json:"id"`
		} `json:"data"`
	}
	json.Unmarshal(questionResp, &questionObj)
	questionID := questionObj.Data.ID

	// 3. Test WS Auth with Query Param
	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/api/v1/ws?token=" + token
	dialer := websocket.Dialer{}
	conn, resp, err := dialer.Dial(wsURL, nil)
	require.NoError(t, err)
	defer conn.Close()
	assert.Equal(t, http.StatusSwitchingProtocols, resp.StatusCode)

	// Join Quiz
	joinMsg := map[string]interface{}{
		"type": "join_quiz",
		"payload": map[string]string{
			"quiz_id": quizID,
		},
	}
	err = conn.WriteJSON(joinMsg)
	require.NoError(t, err)
	time.Sleep(100 * time.Millisecond)

	// 4. Test Idempotency
	submitBody := fmt.Sprintf(`{"question_id":"%s","answer":"A"}`, questionID)

	// First Submission
	submitResp1 := requestWithAuth(t, server, "POST", fmt.Sprintf("/api/v1/quizzes/%s/submit", quizID), submitBody, token)
	assert.Contains(t, string(submitResp1), "points") // basic check for success

	// Second Submission (Should fail or return specific message)
	submitResp2 := requestWithAuth(t, server, "POST", fmt.Sprintf("/api/v1/quizzes/%s/submit", quizID), submitBody, token)
	// We expect 500 or 400 with our current error handling, ideally 409 Conflict, but let's check the error message
	// Since we return error in service, and handler likely wraps it.
	// We need to check response body for "already answered"
	t.Logf("Second Submission Response: %s", string(submitResp2))
	assert.Contains(t, string(submitResp2), "already answered")

	// Read WS message to verify One broadcast (optional, or check broadcast content)
	conn.SetReadDeadline(time.Now().Add(1 * time.Second))
	_, msg, err := conn.ReadMessage()
	require.NoError(t, err)
	var wsMsg realtime.WSMessage
	json.Unmarshal(msg, &wsMsg)
	assert.Equal(t, "leaderboard_update", wsMsg.Type)
}
