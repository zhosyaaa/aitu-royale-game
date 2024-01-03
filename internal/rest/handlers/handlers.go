package handlers

import (
	"auth/internal/config"
	redis2 "auth/internal/db/redis"
	"auth/internal/repository"
	"auth/internal/rest/forms"
	"auth/internal/rest/models"
	"auth/pkg/email"
	"auth/pkg/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

type Handlers struct {
	repo        repository.UserRepo
	redisConfig config.RedisConfig
	email       config.EmailConfig
}

func (h Handlers) Register(context *gin.Context) {
	var user models.User
	if err := context.BindJSON(&user); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	hashedPassword, _ := utils.HashPassword(user.Password)
	user.Password = hashedPassword
	if err := h.repo.CreateUser(&user); err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	signedToken, _ := utils.CreateToken(strconv.Itoa(int(user.ID)), user.Email)
	cookie := http.Cookie{
		Name:     "jwt",
		Value:    signedToken,
		Path:     "/auth",
		Expires:  time.Now().Add(time.Hour * 24),
		HttpOnly: true,
	}
	http.SetCookie(context.Writer, &cookie)
	context.JSON(http.StatusOK, gin.H{"message": "User registered successfully", "data": signedToken})
}

func (h Handlers) Login(context *gin.Context) {
	var data forms.LoginForm
	if err := context.BindJSON(&data); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	user, err := h.repo.GetUserByUsername(data.Username)
	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	if !utils.CheckPasswordHash(data.Password, user.Password) {
		context.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
		return
	}
	token, err := utils.CreateToken(strconv.Itoa(int(user.ID)), user.Email)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create token", "data": token})
		return
	}
	context.JSON(http.StatusOK, gin.H{"token": token})
}

func (h Handlers) Profile(context *gin.Context) {
	idInterface, exists := context.Get("id")
	if !exists {
		context.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "User not authenticated",
		})
		return
	}

	id, ok := idInterface.(uint)
	if !ok {
		context.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Error while retrieving user ID",
		})
		return
	}

	user, err := h.repo.GetUserByID(id)
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "user does not exist",
			"data":    err,
		})
		return
	}
	context.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "User successfully found",
		"data":    user,
	})
}
func (h Handlers) Logout(context *gin.Context) {
	cookie := http.Cookie{
		Name:     "jwt",
		Value:    "",
		Path:     "/auth",
		Expires:  time.Now().Add(-time.Hour),
		HttpOnly: true,
	}
	http.SetCookie(context.Writer, &cookie)
	context.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "User logged out successfully",
		"data":    nil,
	})
}
func (h Handlers) ForgotPassword(context *gin.Context) {
	var requestData forms.ForgotPasswordForm
	if err := context.BindJSON(&requestData); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	user, err := h.repo.GetUserByEmail(requestData.Email)
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	verificationCode := utils.GenerateVerificationCode()

	err = redis2.SaveVerificationCodeToRedis(context, h.redisConfig, user.Email, verificationCode)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save verification code"})
		return
	}
	err = email.SendVerificationCodeEmail(user.Email, verificationCode, h.email)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send verification code"})
		return
	}
	context.JSON(http.StatusOK, gin.H{"message": "Verification code sent to your email"})
}

func (h Handlers) DeleteAccount(context *gin.Context) {

}
