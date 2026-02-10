package service

import (
	"context"
	"crypto/rand"
	"errors"
	"math"
	"math/big"

	"github.com/google/uuid"
	"github.com/nguyen1302/realtime-quiz/internal/models"
	"github.com/nguyen1302/realtime-quiz/internal/realtime"
	"github.com/nguyen1302/realtime-quiz/internal/repository"
)

var ErrAlreadyAnswered = errors.New("you have already answered this question")

type QuizService interface {
	CreateQuiz(ctx context.Context, title, description string, ownerID uuid.UUID) (*models.Quiz, error)
	AddQuestion(ctx context.Context, input AddQuestionInput) (*models.Question, error)
	GetQuiz(ctx context.Context, id uuid.UUID) (*models.Quiz, error)
	JoinQuiz(ctx context.Context, code string) (*models.Quiz, error)
	SubmitAnswer(ctx context.Context, input SubmitAnswerInput) (*models.Answer, error)
	GetLeaderboard(ctx context.Context, quizID uuid.UUID) ([]models.LeaderboardEntry, error)
}

type AddQuestionInput struct {
	QuizID        uuid.UUID
	Text          string
	Options       []string
	CorrectAnswer string
	TimeLimit     int
	Points        int
	Order         int
}

type SubmitAnswerInput struct {
	QuizID     uuid.UUID
	QuestionID uuid.UUID
	UserID     uuid.UUID
	Answer     string
}

type quizService struct {
	quizRepo        repository.QuizRepository
	questionRepo    repository.QuestionRepository
	leaderboardRepo repository.LeaderboardRepository
	answerRepo      repository.AnswerRepository
	realtimeService RealtimeService
}

func NewQuizService(quizRepo repository.QuizRepository, questionRepo repository.QuestionRepository, leaderboardRepo repository.LeaderboardRepository, answerRepo repository.AnswerRepository, realtimeService RealtimeService) QuizService {
	return &quizService{
		quizRepo:        quizRepo,
		questionRepo:    questionRepo,
		leaderboardRepo: leaderboardRepo,
		answerRepo:      answerRepo,
		realtimeService: realtimeService,
	}
}

func (s *quizService) CreateQuiz(ctx context.Context, title, description string, ownerID uuid.UUID) (*models.Quiz, error) {
	code, err := generateQuizCode()
	if err != nil {
		return nil, err
	}

	quiz := &models.Quiz{
		Title:       title,
		Description: description,
		Code:        code,
		Status:      models.QuizStatusDraft,
		OwnerID:     ownerID,
	}

	if err := s.quizRepo.Create(ctx, quiz); err != nil {
		return nil, err
	}

	return quiz, nil
}

func (s *quizService) AddQuestion(ctx context.Context, input AddQuestionInput) (*models.Question, error) {
	if input.TimeLimit == 0 {
		input.TimeLimit = 30 // Default
	}
	if input.Points == 0 {
		input.Points = 100 // Default
	}

	question := &models.Question{
		QuizID:        input.QuizID,
		Text:          input.Text,
		Options:       input.Options,
		CorrectAnswer: input.CorrectAnswer,
		TimeLimit:     input.TimeLimit,
		Points:        input.Points,
		Order:         input.Order,
	}

	if err := s.questionRepo.Create(ctx, question); err != nil {
		return nil, err
	}
	return question, nil
}

func (s *quizService) GetQuiz(ctx context.Context, id uuid.UUID) (*models.Quiz, error) {
	// Get quiz info
	quiz, err := s.quizRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return quiz, nil
}

func (s *quizService) JoinQuiz(ctx context.Context, code string) (*models.Quiz, error) {
	quiz, err := s.quizRepo.GetByCode(ctx, code)
	if err != nil {
		return nil, err
	}
	return quiz, nil
}

func (s *quizService) SubmitAnswer(ctx context.Context, input SubmitAnswerInput) (*models.Answer, error) {
	// 1. Fetch Question to check Answer
	question, err := s.questionRepo.GetByID(ctx, input.QuestionID)
	if err != nil {
		return nil, err
	}

	// 2. Check correctness
	isCorrect := question.CorrectAnswer == input.Answer
	points := 0

	if isCorrect {
		// 3. Get rank from Redis (Atomic) if correct
		rank, err := s.leaderboardRepo.GetSubmissionRank(ctx, input.QuizID, input.QuestionID)
		if err != nil {
			return nil, err
		}

		// 4. Calculate points
		// Logic: MaxPoints * (0.9 ^ (rank-1))
		// Example (Max 100):
		// Rank 1: 100
		// Rank 2: 90
		// Rank 3: 81
		maxPoints := float64(question.Points)
		if maxPoints == 0 {
			maxPoints = 1000 // Fallback
		}

		decayFactor := math.Pow(0.9, float64(rank-1))
		calculatedPoints := maxPoints * decayFactor

		// Minimum floor: 10% of maxPoints or 1 (whichever is higher)
		minPoints := maxPoints * 0.1
		if minPoints < 1 {
			minPoints = 1
		}

		if calculatedPoints < minPoints {
			calculatedPoints = minPoints
		}
		points = int(calculatedPoints)
	}

	answer := &models.Answer{
		QuizID:     input.QuizID,
		QuestionID: input.QuestionID,
		UserID:     input.UserID,
		Answer:     input.Answer,
		IsCorrect:  isCorrect,
		Points:     points,
	}

	// 5. Update DB (Transaction?)
	// Check idempotency
	hasAnswered, err := s.answerRepo.HasAnswered(ctx, input.QuizID, input.QuestionID, input.UserID)
	if err != nil {
		return nil, err
	}
	if hasAnswered {
		return nil, ErrAlreadyAnswered
	}

	if err := s.answerRepo.Create(ctx, answer); err != nil {
		return nil, err
	}

	// 6. Update Leaderboard (Real-time)
	if err := s.leaderboardRepo.UpdateScore(ctx, input.QuizID, input.UserID, float64(points)); err != nil {
		return nil, err
	}

	// 7. Broadcast Leaderboard Update
	leaderboard, _ := s.leaderboardRepo.GetLeaderboard(ctx, input.QuizID, 10)
	s.realtimeService.BroadcastToQuiz(input.QuizID.String(), realtime.EventLeaderboard, leaderboard)

	return answer, nil
}

func generateQuizCode() (string, error) {
	const charset = "0123456789"
	code := make([]byte, 6)
	for i := range code {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		code[i] = charset[num.Int64()]
	}
	return string(code), nil
}

func (s *quizService) GetLeaderboard(ctx context.Context, quizID uuid.UUID) ([]models.LeaderboardEntry, error) {
	// Get top 10
	return s.leaderboardRepo.GetLeaderboard(ctx, quizID, 10)
}
