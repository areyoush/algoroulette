package main

import (
	"log"
	"os"
	"time"

	"github.com/areyoush/algoroulette/internal/db"
	"github.com/areyoush/algoroulette/internal/handler"
	"github.com/areyoush/algoroulette/internal/middleware"
	"github.com/areyoush/algoroulette/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/ulule/limiter/v3"
	ginlimiter "github.com/ulule/limiter/v3/drivers/middleware/gin"
	memory "github.com/ulule/limiter/v3/drivers/store/memory"
)

func main() {
	_ = godotenv.Load()
	
	database, err := db.Connect()
	if err != nil {
		log.Fatal("Could not connect to DB:", err)
	}
	defer database.Close()
	
	database.SetMaxOpenConns(25)
	database.SetMaxIdleConns(25)
	database.SetConnMaxLifetime(5 * time.Minute)

	questionRepo := repository.NewQuestionRepository(database)
	userRepo := repository.NewUserRepository(database)
	
	go func() {
		log.Printf("Initializing token denylist cleanup...")
		if err := userRepo.CleanupDenylist(); err != nil {
        	log.Printf("Initial cleanup error: %v", err)
    	}
		ticker := time.NewTicker(24 * time.Hour)
		defer ticker.Stop()
		
		for range ticker.C {
			log.Println("Running token denylist garbage collection...")
			if err := userRepo.CleanupDenylist(); err != nil {
				log.Printf("Background worker error: %v", err)
			}
		}
	}()

	questionHandler := handler.NewQuestionHandler(questionRepo)
	authHandler := handler.NewAuthHandler(userRepo)

	r := gin.Default()

	rate, err := limiter.NewRateFromFormatted("5-M")
	if err != nil {
		log.Fatal("Could not create rate limiter:", err)
	}

	store := memory.NewStore()
	instance := limiter.New(store, rate)
	rateLimiter := ginlimiter.NewMiddleware(instance, ginlimiter.WithLimitReachedHandler(func(c *gin.Context) {
		c.JSON(429, gin.H{"error": "too many attempts, please wait"})
		c.Abort()
	}), ginlimiter.WithKeyGetter(func(c *gin.Context) string {
		return c.ClientIP()
	}))

	r.Use(func(c *gin.Context) {
		// Set your allowed origin via env variable, fallback to localhost
		allowedOrigin := os.Getenv("FRONTEND_URL")
		if allowedOrigin == "" {
			allowedOrigin = "http://localhost:3000" // Adjust port if your frontend uses a different one
		}

		origin := c.Request.Header.Get("Origin")
		if origin == allowedOrigin {
			c.Header("Access-Control-Allow-Origin", origin)
		}

		c.Header("Access-Control-Allow-Methods", "GET, POST, DELETE, PATCH, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		
		// Required for cross-origin cookies
		c.Header("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	r.Static("/css", "./static/css")
	r.Static("/js", "./static/js")
	r.StaticFile("/", "./static/index.html")


	r.POST("/auth/register", rateLimiter, authHandler.Register)
	r.POST("/auth/login", rateLimiter, authHandler.Login)
	r.POST("/auth/logout", authHandler.Logout)


	protected := r.Group("/")
	protected.Use(middleware.AuthRequired(userRepo))
	protected.GET("/question", questionHandler.GetRandom)
	protected.POST("/question", questionHandler.Create)
	protected.POST("/questions/import", questionHandler.Import)
	protected.DELETE("/questions", questionHandler.ClearAll)
	protected.PATCH("/question/:id/status", questionHandler.UpdateStatus)
	protected.PATCH("/question/:id/bookmark", questionHandler.UpdateBookmark)
	protected.PATCH("/question/:id/notes", questionHandler.UpdateNotes)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	
	r.Run(":" + port)
}