package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gaming-leaderboard/internal/database"
	"github.com/gaming-leaderboard/internal/handlers"
	"github.com/gaming-leaderboard/internal/middleware"
	"github.com/gin-gonic/gin"
	nrgin "github.com/newrelic/go-agent/v3/integrations/nrgin"
	"github.com/newrelic/go-agent/v3/newrelic"
)

func setupRouter(handler *handlers.Handler) *gin.Engine {
	router := gin.Default()

	app, err := setupNewRelic()
	if err != nil {
		log.Fatalf("Failed to set up New Relic: %v", err)
	}
	// New Relic middleware
	router.Use(nrgin.Middleware(app))

	// Add CORS middleware
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Timestamp, X-Nonce, X-Signature")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Serve static files
	router.Static("/leaderboard", "./static")

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Initialize security middleware
	secretKey := os.Getenv("API_SECRET_KEY")
	if secretKey == "" {
		log.Fatal("API_SECRET_KEY environment variable is required")
	}
	securityMiddleware := middleware.NewSecurityMiddleware(secretKey)

	// Configure rate limiters for different endpoints
	submitLimiter := middleware.NewRateLimiter(5, 10, 1*time.Hour)
	readLimiter := middleware.NewRateLimiter(10, 20, 1*time.Hour)

	// API routes with rate limiting and security
	api := router.Group("/api")
	{
		leaderboard := api.Group("/leaderboard")
		{
			// Apply stricter rate limit and security for submission endpoint
			leaderboard.POST("/submit",
				submitLimiter.Middleware(),
				securityMiddleware.ValidateRequest(),
				handler.SubmitScore,
			)

			// Apply more lenient rate limit for read endpoints
			leaderboard.GET("/top", readLimiter.Middleware(), handler.GetLeaderboard)
			leaderboard.GET("/rank/:user_id", readLimiter.Middleware(), handler.GetPlayerRank)
		}
	}
	go cleanUpNonces(securityMiddleware)
	return router
}

func setupNewRelic() (*newrelic.Application, error) {
	app, err := newrelic.NewApplication(
		newrelic.ConfigAppName("gaming-leaderboard"),
		newrelic.ConfigLicense(os.Getenv("NEW_RELIC_LICENSE_KEY")),
		newrelic.ConfigAppLogForwardingEnabled(true),
	)
	return app, err
}

func cleanUpNonces(m *middleware.SecurityMiddleware) {
	for {
		m.CleanupNonces()
		time.Sleep(10 * time.Minute)
	}
}

func main() {
	// Parse command line flags
	migrateDown := flag.Bool("down", false, "Run down migrations instead of up migrations")
	version := flag.Int("version", -1, "Migrate to a specific version (-1 for all migrations)")
	flag.Parse()

	// Set Gin mode based on environment
	ginMode := os.Getenv("GIN_MODE")
	if ginMode == "" {
		ginMode = "debug"
	}
	gin.SetMode(ginMode)

	// Database configuration
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "localhost"
	}
	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		dbUser = "postgres"
	}
	dbPass := os.Getenv("DB_PASSWORD")
	if dbPass == "" {
		dbPass = "postgres"
	}
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "leaderboard"
	}

	// Initialize database connection
	db, err := database.NewDB(
		dbHost,
		dbUser,
		dbPass,
		dbName,
		5432,
	)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Run database migrations
	if err := db.MigrateDB(*migrateDown, *version); err != nil {
		log.Fatalf("Failed to run database migrations: %v", err)
	}

	if *migrateDown || *version >= 0 {
		log.Println("Migration completed successfully")
		return
	}
	log.Println("Database migrations completed successfully")

	// Start background worker for leaderboard updates
	workerQuit := make(chan struct{})

	// Initialize handler and router
	handler := handlers.NewHandler(db)
	router := setupRouter(handler)

	// Configure server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Server starting on port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Stop the background worker
	close(workerQuit)

	// Give outstanding operations up to 5 seconds to complete
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}
