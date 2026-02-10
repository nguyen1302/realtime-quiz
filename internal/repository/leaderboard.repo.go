package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/nguyen1302/realtime-quiz/internal/models"
	"github.com/redis/go-redis/v9"
)

type LeaderboardRepository interface {
	GetSubmissionRank(ctx context.Context, quizID, questionID uuid.UUID) (int64, error)
	UpdateScore(ctx context.Context, quizID uuid.UUID, userID uuid.UUID, points float64) error
	GetLeaderboard(ctx context.Context, quizID uuid.UUID, limit int64) ([]models.LeaderboardEntry, error)
}

type leaderboardRepository struct {
	rdb *redis.Client
}

func NewLeaderboardRepository(rdb *redis.Client) LeaderboardRepository {
	return &leaderboardRepository{rdb: rdb}
}

func (r *leaderboardRepository) GetSubmissionRank(ctx context.Context, quizID, questionID uuid.UUID) (int64, error) {
	key := fmt.Sprintf("quiz:%s:question:%s:submissions", quizID, questionID)
	// INCR returns the new value. 1st submission gets 1, 2nd gets 2, etc.
	// We might want to set expiry on this key if it doesn't exist?
	// But simple INCR is enough for logic.
	rank, err := r.rdb.Incr(ctx, key).Result()
	if err != nil {
		return 0, err
	}
	// Set expiry for cleanup (e.g. 24h)
	if rank == 1 {
		r.rdb.Expire(ctx, key, 24*time.Hour)
	}
	return rank, nil
}

func (r *leaderboardRepository) UpdateScore(ctx context.Context, quizID uuid.UUID, userID uuid.UUID, points float64) error {
	key := fmt.Sprintf("quiz:%s:leaderboard", quizID)
	// ZINCRBY updates the score
	err := r.rdb.ZIncrBy(ctx, key, points, userID.String()).Err()
	if err != nil {
		return err
	}
	// Ensure Leaderboard expires eventually
	r.rdb.Expire(ctx, key, 24*time.Hour)
	return nil
}

func (r *leaderboardRepository) GetLeaderboard(ctx context.Context, quizID uuid.UUID, limit int64) ([]models.LeaderboardEntry, error) {
	key := fmt.Sprintf("quiz:%s:leaderboard", quizID)
	// ZREVRANGE to get top scores (highest first). WithScores to get score.
	results, err := r.rdb.ZRevRangeWithScores(ctx, key, 0, limit-1).Result()
	if err != nil {
		return nil, err
	}

	entries := make([]models.LeaderboardEntry, len(results))
	for i, z := range results {
		uid, err := uuid.Parse(z.Member.(string))
		if err != nil {
			continue // Should not happen if we store UUID strings
		}
		entries[i] = models.LeaderboardEntry{
			UserID: uid,
			Score:  z.Score,
			Rank:   i + 1,
		}
	}
	return entries, nil
}
