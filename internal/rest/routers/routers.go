package routers

import (
	"auth/internal/rest/handlers"
	"github.com/gin-gonic/gin"
)

type Routers struct {
	handlers handlers.Handlers
}

func (r *Routers) SetupRoutes(app *gin.Engine) {
	authRouter := app.Group("/auth")
	{
		authRouter.POST("/register", r.handlers.Register)
		authRouter.POST("/login", r.handlers.Login)
		authRouter.POST("/forgotPassword", r.handlers.ForgotPassword)
		authRouter.POST("/logout", r.handlers.Logout)
		authRouter.DELETE("/deleteAccount", r.handlers.DeleteAccount)
		authRouter.GET("/profile", r.handlers.Profile)
	}
}
