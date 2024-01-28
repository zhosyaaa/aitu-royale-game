package main

import (
	"auth/internal/config"
	"auth/internal/db"
	"auth/internal/repository"
	"auth/internal/rest/handlers"
	"auth/internal/rest/routers"
	"auth/pkg/logger"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}
}

func initializeDB() config.Database {
	dbConfig := config.Database{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Sslmode:  os.Getenv("DB_SSLMODE"),
		Name:     os.Getenv("DB_NAME"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
	}
	return dbConfig
}

func initializeRedis() config.RedisConfig {
	redisDB, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		log.Fatalf("Error converting REDIS_DB to int: %s", err)
	}

	redisConfig := config.RedisConfig{
		Addr:     os.Getenv("REDIS_ADDRESS"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       redisDB,
	}

	return redisConfig
}

func initializeEmail() config.EmailConfig {
	emailConfig := config.EmailConfig{
		From:     os.Getenv("EMAIL_FROM"),
		Password: os.Getenv("EMAIL_PASSWORD"),
		SMTPHost: os.Getenv("EMAIL_SMTP_HOST"),
		SMTPPort: os.Getenv("EMAIL_SMTP_PORT"),
	}
	return emailConfig
}

var appConfig config.App

func main() {
	logger.InitLogger()

	appConfig = config.App{
		PORT:  os.Getenv("APP_PORT"),
		DB:    initializeDB(),
		Redis: initializeRedis(),
		Email: initializeEmail(),
	}

	db, err := db.GetDBInstance(appConfig.DB)
	if err != nil {
		logger.GetLogger().Fatal("Error initializing DB:", err)
	}

	userRepo := repository.NewUserRepository(db)
	gameRepo := repository.NewGameRepository(db)
	authHandlers := handlers.NewAuthHandlers(userRepo, appConfig.Redis, appConfig.Email)
	var gameHandlers = handlers.NewGameHandlers(userRepo, *gameRepo)

	r := gin.Default()
	router := routers.NewRouters(*authHandlers, *gameHandlers)
	router.SetupRoutes(r)
	r.Use(rateLimitMiddleware())

	server := &http.Server{
		Addr:    ":" + appConfig.PORT,
		Handler: r,
	}

	gracefulShutdown(server)
}

func rateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		rateLimit := time.Tick(time.Second)
		select {
		case <-rateLimit:
			c.Next()
		default:
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "Rate limit exceeded"})
			c.Abort()
		}
	}
}

func gracefulShutdown(server *http.Server) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-stop
		logger.GetLogger().Info("Server is shutting down...")

		timeout := 5 * time.Second
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			logger.GetLogger().Fatal("Server shutdown error:", err)
		}

		logger.GetLogger().Info("Server has gracefully stopped")
		os.Exit(0)
	}()

	logger.GetLogger().Info("Server is running on :" + appConfig.PORT)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.GetLogger().Fatal("Error starting server:", err)
	}
}
