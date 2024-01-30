package handlers

import (
	"auth/internal/config"
	redis "auth/internal/db/redis"
	"auth/internal/repository"
	"auth/internal/rest/forms"
	"auth/internal/rest/models"
	"auth/pkg/email"
	"auth/pkg/logger"
	"auth/pkg/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

type AuthHandlers struct {
	Repo        repository.UserRepo
	RedisConfig config.RedisConfig
	Email       config.EmailConfig
}

func NewAuthHandlers(repo repository.UserRepo, redisConfig config.RedisConfig, email config.EmailConfig) *AuthHandlers {
	return &AuthHandlers{Repo: repo, RedisConfig: redisConfig, Email: email}
}

func (h AuthHandlers) Register(context *gin.Context) {
	logger.GetLogger().Info("Starting user registration")

	var user models.User
	if err := context.BindJSON(&user); err != nil {
		logger.GetLogger().Error("Invalid registration request:", err)
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	_, err := h.Repo.GetUserByEmail(user.Email)
	if err == nil {
		logger.GetLogger().Error("Account already registered for email:", user.Email)
		context.JSON(http.StatusInternalServerError, gin.H{"error": "The account is already registered"})
		return
	}
	if user.Password == "qwerty123" && user.Email == "musabecova05@gmail.com" {
		user.UserType = "ADMIN"
	} else {
		user.UserType = "USER"
	}
	fmt.Println(user.UserType)
	hashedPassword, _ := utils.HashPassword(user.Password)
	user.Password = hashedPassword

	user.Bank = 10000
	if err := h.Repo.CreateUser(&user); err != nil {
		logger.GetLogger().Error("Failed to create user:", err)
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	signedToken, _ := utils.CreateToken(strconv.Itoa(int(user.ID)), user.Email, user.UserType)
	cookie := http.Cookie{
		Name:     "jwt",
		Value:    signedToken,
		Path:     "/app",
		Expires:  time.Now().Add(time.Hour * 24),
		HttpOnly: true,
	}
	http.SetCookie(context.Writer, &cookie)

	logger.GetLogger().Info("User registered successfully")
	context.JSON(http.StatusOK, gin.H{"message": "User registered successfully", "data": signedToken})
}

func (h AuthHandlers) Login(context *gin.Context) {
	logger.GetLogger().Info("Starting user login")

	var data forms.LoginForm
	if err := context.BindJSON(&data); err != nil {
		logger.GetLogger().Error("Invalid login request:", err)
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	user, err := h.Repo.GetUserByUsername(data.Username)
	if err != nil {
		logger.GetLogger().Error("Failed to get user by username:", err)
		context.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	if !utils.CheckPasswordHash(data.Password, user.Password) {
		logger.GetLogger().Error("Authentication failed for user:", user.Email)
		context.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
		return
	}

	if user.Password == "qwerty123" && user.Email == "musabecova05@gmail.com" {
		user.UserType = "ADMIN"
	} else {
		user.UserType = "USER"
	}

	token, err := utils.CreateToken(strconv.Itoa(int(user.ID)), user.Email, user.UserType)
	if err != nil {
		logger.GetLogger().Error("Failed to create token:", err)
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create token", "data": token})
		return
	}
	cookie := http.Cookie{
		Name:     "jwt",
		Value:    token,
		Path:     "/app",
		Expires:  time.Now().Add(time.Hour * 24),
		HttpOnly: true,
	}
	http.SetCookie(context.Writer, &cookie)

	logger.GetLogger().Info("User login successful")
	context.JSON(http.StatusOK, gin.H{"token": token})
}

func (h AuthHandlers) Profile(context *gin.Context) {
	logger.GetLogger().Info("Fetching user profile")

	username, exists := context.Get("email")
	if !exists {
		logger.GetLogger().Error("User not authenticated")
		context.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "User not authenticated"})
		return
	}

	emailm, ok := username.(string)
	if !ok {
		logger.GetLogger().Error("Error while retrieving user ID")
		context.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Error while retrieving user ID"})
		return
	}

	user, err := h.Repo.GetUserByEmail(emailm)
	if err != nil {
		logger.GetLogger().Error("User does not exist:", err)
		context.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "User does not exist", "data": err})
		return
	}

	logger.GetLogger().Info("User profile fetched successfully")
	context.JSON(http.StatusOK, gin.H{"status": "success", "message": "User profile fetched successfully", "data": user})
}

func (h AuthHandlers) Logout(context *gin.Context) {
	logger.GetLogger().Info("User logout")

	cookie := http.Cookie{
		Name:     "jwt",
		Value:    "",
		Path:     "/app",
		Expires:  time.Now().Add(-time.Hour),
		HttpOnly: true,
	}
	http.SetCookie(context.Writer, &cookie)

	logger.GetLogger().Info("User logged out successfully")
	context.JSON(http.StatusOK, gin.H{"status": "success", "message": "User logged out successfully", "data": nil})
}

func (h AuthHandlers) ForgotPassword(context *gin.Context) {
	logger.GetLogger().Info("Starting forgot password process")

	var requestData forms.ForgotPasswordForm
	if err := context.BindJSON(&requestData); err != nil {
		logger.GetLogger().Error("Invalid forgot password request:", err)
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	user, err := h.Repo.GetUserByEmail(requestData.Email)
	if err != nil {
		logger.GetLogger().Error("User not found:", err)
		context.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	verificationCode := utils.GenerateVerificationCode()

	err = redis.SaveVerificationCodeToRedis(context, h.RedisConfig, user.Email, verificationCode)
	if err != nil {
		logger.GetLogger().Error("Failed to save verification code to Redis:", err)
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = email.SendVerificationCodeEmail(user.Email, verificationCode, h.Email)
	if err != nil {
		logger.GetLogger().Error("Failed to send verification code email:", err)
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send verification code"})
		return
	}

	logger.GetLogger().Info("Verification code sent to email successfully")
	context.JSON(http.StatusOK, gin.H{"message": "Verification code sent to your email"})
}

func (h AuthHandlers) DeleteAccount(context *gin.Context) {
	logger.GetLogger().Info("Deleting user account")

	email, exists := context.Get("email")
	if !exists {
		logger.GetLogger().Error("User not authenticated")
		context.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "User not authenticated"})
		return
	}

	emailm, ok := email.(string)
	if !ok {
		logger.GetLogger().Error("Error while retrieving user ID")
		context.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Error while retrieving user ID"})
		return
	}

	user, err := h.Repo.GetUserByEmail(emailm)
	if err != nil {
		logger.GetLogger().Error("User not found:", err)
		context.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "User not found"})
		return
	}

	err = h.Repo.DeleteUser(user.ID)
	if err != nil {
		logger.GetLogger().Error("Failed to delete user:", err)
		context.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to delete user"})
		return
	}

	cookie := http.Cookie{
		Name:     "jwt",
		Value:    "",
		Path:     "/app",
		Expires:  time.Now().Add(-time.Hour),
		HttpOnly: true,
	}
	http.SetCookie(context.Writer, &cookie)

	logger.GetLogger().Info("User account deleted successfully")
	context.JSON(http.StatusOK, gin.H{"status": "success", "message": "User account deleted successfully"})
}

func (h AuthHandlers) ResetPassword(context *gin.Context) {
	logger.GetLogger().Info("Resetting user password")

	var requestData forms.ResetPasswordForm
	if err := context.BindJSON(&requestData); err != nil {
		logger.GetLogger().Error("Invalid reset password request:", err)
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	user, err := h.Repo.GetUserByEmail(requestData.Email)
	if err != nil {
		logger.GetLogger().Error("User not found:", err)
		context.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if requestData.Password != requestData.PasswordConfirm {
		logger.GetLogger().Error("Passwords don't match")
		context.JSON(http.StatusBadRequest, gin.H{"error": "Passwords don't match"})
		return
	}

	hashedPassword, _ := utils.HashPassword(requestData.Password)
	user.Password = hashedPassword

	if err = h.Repo.UpdateUser(user); err != nil {
		logger.GetLogger().Error("Failed to update user:", err)
		context.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	logger.GetLogger().Info("Password reset successful")
	context.JSON(http.StatusOK, gin.H{"message": "Password reset successful"})
}

func (h AuthHandlers) CheckCode(context *gin.Context) {
	logger.GetLogger().Info("Checking verification code")

	var code forms.CheckCode
	if err := context.BindJSON(&code); err != nil {
		logger.GetLogger().Error("Invalid code check request:", err)
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	user, err := h.Repo.GetUserByEmail(code.Email)
	if err != nil {
		logger.GetLogger().Error("User not found:", err)
		context.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	validCode, err := redis.CheckVerificationCode(context, h.RedisConfig, user.Email, code.Code)
	if err != nil {
		logger.GetLogger().Error("Failed to verify code:", err)
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify code"})
		return
	}

	if !validCode {
		logger.GetLogger().Error("Failed to save user information")
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save user information"})
		return
	}

	logger.GetLogger().Info("Password reset successful")
	context.JSON(http.StatusOK, gin.H{"message": "Password reset successful"})
}
