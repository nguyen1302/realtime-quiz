package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nguyen1302/realtime-quiz/internal/service"
	"github.com/nguyen1302/realtime-quiz/pkg/response"
)

type QuizHandler interface {
	CreateQuiz(c *gin.Context)
	AddQuestion(c *gin.Context)
	GetQuiz(c *gin.Context)
	JoinQuiz(c *gin.Context)
}

type quizHandler struct {
	quizService service.QuizService
}

func NewQuizHandler(quizService service.QuizService) QuizHandler {
	return &quizHandler{quizService: quizService}
}

type CreateQuizRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
}

type AddQuestionRequest struct {
	Text          string   `json:"text" binding:"required"`
	Options       []string `json:"options" binding:"required,min=2"`
	CorrectAnswer string   `json:"correct_answer" binding:"required"`
	TimeLimit     int      `json:"time_limit"`
	Points        int      `json:"points"`
	Order         int      `json:"order"`
}

type JoinQuizRequest struct {
	Code string `json:"code" binding:"required,len=6"`
}

// POST /api/v1/quizzes
func (h *quizHandler) CreateQuiz(c *gin.Context) {
	var req CreateQuizRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	userID := c.MustGet("userID").(uuid.UUID)
	quiz, err := h.quizService.CreateQuiz(c.Request.Context(), req.Title, req.Description, userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to create quiz", nil)
		return
	}

	response.Success(c, http.StatusCreated, "Quiz created", quiz)
}

// POST /api/v1/quizzes/:id/questions
func (h *quizHandler) AddQuestion(c *gin.Context) {
	quizID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid quiz ID", nil)
		return
	}

	var req AddQuestionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	input := service.AddQuestionInput{
		QuizID:        quizID,
		Text:          req.Text,
		Options:       req.Options,
		CorrectAnswer: req.CorrectAnswer,
		TimeLimit:     req.TimeLimit,
		Points:        req.Points,
		Order:         req.Order,
	}

	question, err := h.quizService.AddQuestion(c.Request.Context(), input)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to add question", nil)
		return
	}

	response.Success(c, http.StatusCreated, "Question added", question)
}

// GET /api/v1/quizzes/:id
func (h *quizHandler) GetQuiz(c *gin.Context) {
	quizID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid quiz ID", nil)
		return
	}

	quiz, err := h.quizService.GetQuiz(c.Request.Context(), quizID)
	if err != nil {
		response.Error(c, http.StatusNotFound, "Quiz not found", nil)
		return
	}

	response.Success(c, http.StatusOK, "Quiz info retrieved", quiz)
}

// POST /api/v1/quizzes/join
func (h *quizHandler) JoinQuiz(c *gin.Context) {
	var req JoinQuizRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	quiz, err := h.quizService.JoinQuiz(c.Request.Context(), req.Code)
	if err != nil {
		response.Error(c, http.StatusNotFound, "Quiz not found", nil)
		return
	}

	response.Success(c, http.StatusOK, "Joined quiz successfully", quiz)
}
