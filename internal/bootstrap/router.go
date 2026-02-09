package bootstrap

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nguyen1302/realtime-quiz/internal/config"
	"github.com/nguyen1302/realtime-quiz/internal/handler"
	"github.com/nguyen1302/realtime-quiz/internal/middleware"
	"github.com/nguyen1302/realtime-quiz/internal/repository"
	"github.com/nguyen1302/realtime-quiz/internal/service"
)

type Router struct {
	engine   *gin.Engine
	handlers *handler.Handlers
	services *service.Services
}

func NewRouter(db *pgxpool.Pool, cfg *config.Config) *Router {
	// Initialize managers
	repos := repository.NewRepositories(db)
	services := service.NewServices(repos, cfg)
	handlers := handler.NewHandlers(services)

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
		auth.POST("/register", r.handlers.Auth.Register)
		auth.POST("/login", r.handlers.Auth.Login)
	}

	// Protected routes
	protected := api.Group("")
	protected.Use(middleware.AuthMiddleware(r.services.Auth))
	{
		protected.GET("/auth/me", r.handlers.Auth.GetMe)
	}

	// Health check
	r.engine.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
}

func (r *Router) Engine() *gin.Engine {
	return r.engine
}
