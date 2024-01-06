package routers

import (
	"auth/internal/rest/handlers"
	"auth/pkg/middleware"
	"github.com/gin-gonic/gin"
)

type Routers struct {
	Handlers handlers.Handlers
}

func (r *Routers) SetupRoutes(app *gin.Engine) {
	authRouter := app.Group("/auth")
	{
		authRouter.POST("/register", r.Handlers.Register)                                               //+
		authRouter.POST("/login", middleware.RequireAuthMiddleware, r.Handlers.Login)                   //+
		authRouter.POST("/forgotPassword", r.Handlers.ForgotPassword)                                   //+
		authRouter.POST("/resetPassword", r.Handlers.ResetPassword)                                     //
		authRouter.POST("/checkVerificationCode", r.Handlers.CheckCode)                                 //+
		authRouter.POST("/logout", middleware.RequireAuthMiddleware, r.Handlers.Logout)                 //+
		authRouter.DELETE("/deleteAccount", middleware.RequireAuthMiddleware, r.Handlers.DeleteAccount) //+
		authRouter.GET("/profile", middleware.RequireAuthMiddleware, r.Handlers.Profile)                //+
	}
}
