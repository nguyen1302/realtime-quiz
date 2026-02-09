package handler

import (
	"errors"
	"net/http"

	"github.com/nguyen1302/realtime-quiz/internal/service"
	"github.com/nguyen1302/realtime-quiz/pkg/response"

	"github.com/gin-gonic/gin"
)

type AuthHandler interface {
	Register(c *gin.Context)
	Login(c *gin.Context)
	GetMe(c *gin.Context)
}

type authHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) AuthHandler {
	return &authHandler{authService: authService}
}

// Register handles user registration
// POST /api/v1/auth/register
func (h *authHandler) Register(c *gin.Context) {
	var req service.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	user, err := h.authService.Register(c.Request.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrEmailExists):
			response.Error(c, http.StatusConflict, "Email already exists", nil)
		case errors.Is(err, service.ErrUsernameExists):
			response.Error(c, http.StatusConflict, "Username already exists", nil)
		default:
			response.Error(c, http.StatusInternalServerError, "Failed to register user", nil)
		}
		return
	}

	response.Success(c, http.StatusCreated, "User registered successfully", gin.H{
		"id":       user.ID,
		"username": user.Username,
		"email":    user.Email,
	})
}

// Login handles user authentication
// POST /api/v1/auth/login
func (h *authHandler) Login(c *gin.Context) {
	var req service.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	token, user, err := h.authService.Login(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			response.Error(c, http.StatusUnauthorized, "Invalid email or password", nil)
			return
		}
		response.Error(c, http.StatusInternalServerError, "Failed to login", nil)
		return
	}

	response.Success(c, http.StatusOK, "Login successful", gin.H{
		"token": token,
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
		},
	})
}

// GetMe returns current authenticated user info
// GET /api/v1/auth/me
func (h *authHandler) GetMe(c *gin.Context) {
	claims, exists := c.Get("claims")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	userClaims := claims.(*service.Claims)
	response.Success(c, http.StatusOK, "User info retrieved", gin.H{
		"id":       userClaims.UserID,
		"username": userClaims.Username,
		"email":    userClaims.Email,
	})
}
