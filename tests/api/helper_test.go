package api_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/nguyen1302/realtime-quiz/internal/bootstrap"
	"github.com/nguyen1302/realtime-quiz/internal/config"
	"github.com/nguyen1302/realtime-quiz/internal/models"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Common setup for tests
func setupTest(t *testing.T) (*gorm.DB, *redis.Client, *httptest.Server) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(&models.User{}, &models.Quiz{}, &models.Question{}, &models.Answer{})
	require.NoError(t, err)

	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		t.Skip("Redis not available, skipping test")
	}
	rdb.FlushAll(context.Background())

	cfg := &config.Config{
		JWT: config.JWTConfig{
			Secret:      "test-secret",
			ExpiryHours: 1,
		},
	}

	app := bootstrap.NewRouter(db, rdb, cfg)
	server := httptest.NewServer(app.Engine())

	// Ensure cleanup
	t.Cleanup(func() {
		server.Close()
		rdb.FlushAll(context.Background())
	})

	return db, rdb, server
}

// Helpers
func request(t *testing.T, server *httptest.Server, method, path, body string) []byte {
	req, err := http.NewRequest(method, server.URL+path, strings.NewReader(body))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	return respBody
}

func requestWithAuth(t *testing.T, server *httptest.Server, method, path, body, token string) []byte {
	req, err := http.NewRequest(method, server.URL+path, strings.NewReader(body))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	return respBody
}

func getToken(t *testing.T, respBody []byte) string {
	var resp struct {
		Data struct {
			Token string `json:"token"`
		} `json:"data"`
	}
	err := json.Unmarshal(respBody, &resp)
	require.NoError(t, err)
	require.NotEmpty(t, resp.Data.Token, "Token should not be empty")
	return resp.Data.Token
}
