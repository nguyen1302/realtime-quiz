package service

import (
	"context"
	"crypto/rand"
	"math/big"

	"github.com/google/uuid"
	"github.com/nguyen1302/realtime-quiz/internal/models"
	"github.com/nguyen1302/realtime-quiz/internal/repository"
)

type QuizService interface {
	CreateQuiz(ctx context.Context, title, description string, ownerID uuid.UUID) (*models.Quiz, error)
	AddQuestion(ctx context.Context, input AddQuestionInput) (*models.Question, error)
	GetQuiz(ctx context.Context, id uuid.UUID) (*models.Quiz, error)
	JoinQuiz(ctx context.Context, code string) (*models.Quiz, error)
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

type quizService struct {
	quizRepo     repository.QuizRepository
	questionRepo repository.QuestionRepository
}

func NewQuizService(quizRepo repository.QuizRepository, questionRepo repository.QuestionRepository) QuizService {
	return &quizService{
		quizRepo:     quizRepo,
		questionRepo: questionRepo,
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
