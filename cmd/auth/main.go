package main

import (
	"auth/internal/config"
	"auth/internal/db"
	"auth/internal/repository"
	"auth/internal/rest/handlers"
	"auth/internal/rest/routers"
	"auth/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
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

func main() {
	logger.InitLogger()

	appConfig := config.App{
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
	var gameHandlers = handlers.NewGameHandlers(*gameRepo)

	r := gin.Default()
	router := routers.NewRouters(*authHandlers, *gameHandlers)
	router.SetupRoutes(r)

	if err := r.Run(":" + appConfig.PORT); err != nil {
		logger.GetLogger().Fatal("Error starting server:", err)
	}
}
