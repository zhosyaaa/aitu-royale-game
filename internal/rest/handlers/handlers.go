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
	Repo        repository.UserRepo
	RedisConfig config.RedisConfig
	Email       config.EmailConfig
}

func (h Handlers) Register(context *gin.Context) {
	var user models.User
	if err := context.BindJSON(&user); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	_, err := h.Repo.GetUserByEmail(user.Email)
	if err == nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "The account is already registered"})
		return
	}
	hashedPassword, _ := utils.HashPassword(user.Password)
	user.Password = hashedPassword
	if err := h.Repo.CreateUser(&user); err != nil {
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
	user, err := h.Repo.GetUserByUsername(data.Username)
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
	username, exists := context.Get("email")
	if !exists {
		context.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "User not authenticated",
		})
		return
	}

	emailm, ok := username.(string)
	if !ok {
		context.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Error while retrieving user ID",
		})
		return
	}

	user, err := h.Repo.GetUserByEmail(emailm)
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
	user, err := h.Repo.GetUserByEmail(requestData.Email)
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	verificationCode := utils.GenerateVerificationCode()

	err = redis2.SaveVerificationCodeToRedis(context, h.RedisConfig, user.Email, verificationCode)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	err = email.SendVerificationCodeEmail(user.Email, verificationCode, h.Email)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send verification code"})
		return
	}
	context.JSON(http.StatusOK, gin.H{"message": "Verification code sent to your email"})
}

func (h Handlers) DeleteAccount(context *gin.Context) {
	email, exists := context.Get("email")
	if !exists {
		context.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "User not authenticated",
		})
		return
	}

	emailm, ok := email.(string)
	if !ok {
		context.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Error while retrieving user ID",
		})
		return
	}
	user, err := h.Repo.GetUserByEmail(emailm)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "User not found",
		})
		return
	}

	err = h.Repo.DeleteUser(user.ID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to delete user",
		})
		return
	}

	cookie := http.Cookie{
		Name:     "jwt",
		Value:    "",
		Path:     "/app/v1",
		Expires:  time.Now().Add(-time.Hour),
		HttpOnly: true,
	}
	http.SetCookie(context.Writer, &cookie)
	context.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "User account deleted successfully",
	})
}

func (h Handlers) ResetPassword(context *gin.Context) {
	var requestData forms.ResetPasswordForm
	if err := context.BindJSON(&requestData); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	user, err := h.Repo.GetUserByEmail(requestData.Email)
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	if requestData.Password != requestData.PasswordConfirm {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Passwords don't match"})
	}
	hashedPassword, _ := utils.HashPassword(requestData.Password)
	user.Password = hashedPassword
	if err = h.Repo.UpdateUser(user); err != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	context.JSON(http.StatusOK, gin.H{"message": "Password reset successful"})
}

func (h Handlers) CheckCode(context *gin.Context) {
	var code forms.CheckCode
	if err := context.BindJSON(&code); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	user, err := h.Repo.GetUserByEmail(code.Email)
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	validCode, err := redis2.CheckVerificationCode(context, h.RedisConfig, user.Email, code.Code)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify code"})
		return
	}

	if !validCode {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save user information"})
		return
	}
	context.JSON(http.StatusOK, gin.H{"message": "Password reset successful"})
}
