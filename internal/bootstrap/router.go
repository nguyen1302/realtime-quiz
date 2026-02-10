package bootstrap

import (
	"github.com/gin-gonic/gin"
	"github.com/nguyen1302/realtime-quiz/internal/config"
	"github.com/nguyen1302/realtime-quiz/internal/handler"
	"github.com/nguyen1302/realtime-quiz/internal/middleware"
	"github.com/nguyen1302/realtime-quiz/internal/repository"
	"github.com/nguyen1302/realtime-quiz/internal/service"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Router struct {
	engine   *gin.Engine
	handlers handler.Handler
	services service.Service
}

func NewRouter(db *gorm.DB, rdb *redis.Client, cfg *config.Config) *Router {
	// Initialize managers
	repos := repository.NewRepository(db, rdb)
	services := service.NewService(repos, cfg)
	handlers := handler.NewHandler(services)

	// Setup Gin
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	engine.Use(gin.Recovery())
	engine.Use(middleware.LoggerMiddleware())
	engine.Use(middleware.CORSMiddleware())

	router := &Router{
		engine:   engine,
		handlers: handlers,
		services: services,
	}

	router.setupRoutes()

	return router
}

func (r *Router) setupRoutes() {
	api := r.engine.Group("/api/v1")

	// Public auth routes
	auth := api.Group("/auth")
	{
		auth.POST("/register", r.handlers.Auth().Register)
		auth.POST("/login", r.handlers.Auth().Login)
	}

	// WebSocket route
	api.GET("/ws", r.handlers.Realtime().HandleConnection)

	// Protected routes
	protected := api.Group("")
	protected.Use(middleware.AuthMiddleware(r.services.Auth()))
	{
		protected.GET("/auth/me", r.handlers.Auth().GetMe)

		// Quiz routes
		quizzes := protected.Group("/quizzes")
		{
			quizzes.POST("", r.handlers.Quiz().CreateQuiz)
			quizzes.GET("/:id", r.handlers.Quiz().GetQuiz)
			quizzes.POST("/:id/questions", r.handlers.Quiz().AddQuestion)
			quizzes.POST("/:id/submit", r.handlers.Quiz().SubmitAnswer)
			quizzes.GET("/:id/leaderboard", r.handlers.Quiz().GetLeaderboard)
			quizzes.POST("/join", r.handlers.Quiz().JoinQuiz)
		}
	}

	// Health check
	r.engine.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
}

func (r *Router) Engine() *gin.Engine {
	return r.engine
}
