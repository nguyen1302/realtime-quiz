package api_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const baseURL = "http://localhost:8080/api/v1"

type LoginResponse struct {
	Data struct {
		Token string `json:"token"`
		User  struct {
			ID string `json:"id"`
		} `json:"user"`
	} `json:"data"`
}

type QuizResponse struct {
	Data struct {
		ID   string `json:"id"`
		Code string `json:"code"`
	} `json:"data"`
}

type QuestionResponse struct {
	Data struct {
		ID string `json:"id"`
	} `json:"data"`
}

type LeaderboardEntry struct {
	UserID string  `json:"user_id"`
	Score  float64 `json:"score"`
	Rank   int     `json:"rank"`
}

type LeaderboardResponse struct {
	Data []LeaderboardEntry `json:"data"`
}

func TestLeaderboard(t *testing.T) {
	// Ensure server is running (This is a black-box integration test)
	// In a real CI/CD, we might spin up a test container here.

	// Unique suffix for this run
	runID := time.Now().UnixNano()

	// 1. Owner Auth
	ownerEmail := fmt.Sprintf("owner_%d@example.com", runID)
	ownerToken, _ := auth(t, ownerEmail, "password123")
	require.NotEmpty(t, ownerToken, "Owner login failed")

	// 2. Create Quiz
	quizID, _ := createQuiz(t, ownerToken)
	require.NotEmpty(t, quizID, "Quiz creation failed")
	t.Logf("Quiz created: %s", quizID)

	// 3. Add Question
	questionID := addQuestion(t, ownerToken, quizID)
	require.NotEmpty(t, questionID, "Question addition failed")

	// 4. Player 1 Auth
	p1Email := fmt.Sprintf("p1_%d@example.com", runID)
	p1Token, p1ID := auth(t, p1Email, "password123")
	require.NotEmpty(t, p1Token)

	// 5. Player 2 Auth
	p2Email := fmt.Sprintf("p2_%d@example.com", runID)
	p2Token, p2ID := auth(t, p2Email, "password123")
	require.NotEmpty(t, p2Token)

	// 6. Player 1 Submits (Fast)
	submitAnswer(t, p1Token, quizID, questionID, "Paris")

	// 7. Player 2 Submits (Slow)
	time.Sleep(100 * time.Millisecond)
	submitAnswer(t, p2Token, quizID, questionID, "Paris")

	// 8. Player 3 Auth & Incorrect Submission
	p3Email := fmt.Sprintf("p3_%d@example.com", runID)
	p3Token, p3ID := auth(t, p3Email, "password123")
	require.NotEmpty(t, p3Token)

	submitAnswer(t, p3Token, quizID, questionID, "London") // Incorrect

	// 9. Get Leaderboard
	entries := getLeaderboard(t, p1Token, quizID)

	// Verification
	require.Len(t, entries, 3, "Leaderboard should have 3 entries (including Player 3 with 0 points)")

	// Player 1 should be Rank 1
	assert.Equal(t, p1ID, entries[0].UserID)
	assert.Equal(t, 1, entries[0].Rank)
	assert.Equal(t, 1000.0, entries[0].Score)

	// Player 2 should be Rank 2
	assert.Equal(t, p2ID, entries[1].UserID)
	assert.Equal(t, 2, entries[1].Rank)
	assert.Equal(t, 900.0, entries[1].Score) // Assuming exponential decay (1000 * 0.9)

	// Player 3 should be in the list
	foundP3 := false
	for _, e := range entries {
		if e.UserID == p3ID {
			foundP3 = true
			assert.Equal(t, 0.0, e.Score, "Player 3 should have 0 points")
			break
		}
	}
	assert.True(t, foundP3, "Player 3 should be found in leaderboard")
}

// Helpers

func auth(t *testing.T, email, password string) (string, string) {
	username := email
	registerPayload := map[string]string{"email": email, "password": password, "username": username}
	post(t, "/auth/register", registerPayload, "")

	loginPayload := map[string]string{"email": email, "password": password}
	resp := post(t, "/auth/login", loginPayload, "")

	var res LoginResponse
	err := json.Unmarshal([]byte(resp), &res)
	require.NoError(t, err)
	return res.Data.Token, res.Data.User.ID
}

func createQuiz(t *testing.T, token string) (string, string) {
	payload := map[string]string{"title": "Test Quiz"}
	resp := post(t, "/quizzes", payload, token)
	var res QuizResponse
	err := json.Unmarshal([]byte(resp), &res)
	require.NoError(t, err)
	return res.Data.ID, res.Data.Code
}

func addQuestion(t *testing.T, token, quizID string) string {
	payload := map[string]interface{}{
		"text":           "Capital of France?",
		"options":        []string{"Paris", "London"},
		"correct_answer": "Paris",
		"points":         1000,
	}
	resp := post(t, fmt.Sprintf("/quizzes/%s/questions", quizID), payload, token)
	var res QuestionResponse
	require.NoError(t, json.Unmarshal([]byte(resp), &res))
	return res.Data.ID
}

func submitAnswer(t *testing.T, token, quizID, questionID, answer string) {
	payload := map[string]string{
		"question_id": questionID,
		"answer":      answer,
	}
	post(t, fmt.Sprintf("/quizzes/%s/submit", quizID), payload, token)
}

func getLeaderboard(t *testing.T, token, quizID string) []LeaderboardEntry {
	resp := get(t, fmt.Sprintf("/quizzes/%s/leaderboard", quizID), token)
	var res LeaderboardResponse
	require.NoError(t, json.Unmarshal([]byte(resp), &res))
	return res.Data
}

func post(t *testing.T, endpoint string, data interface{}, token string) string {
	jsonData, _ := json.Marshal(data)
	req, _ := http.NewRequest("POST", baseURL+endpoint, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	// Check for non-2xx status code if needed, but for now just returning body
	return string(body)
}

func get(t *testing.T, endpoint, token string) string {
	req, _ := http.NewRequest("GET", baseURL+endpoint, nil)
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	return string(body)
}
