package routers

import (
	"auth/internal/rest/handlers"
	"auth/pkg/middleware"
	"github.com/gin-gonic/gin"
)

type Routers struct {
	authHandlers handlers.AuthHandlers
	gameHandlers handlers.GameHandlers
}

func NewRouters(authHandlers handlers.AuthHandlers, gameHandlers handlers.GameHandlers) *Routers {
	return &Routers{authHandlers: authHandlers, gameHandlers: gameHandlers}
}

func (r *Routers) SetupRoutes(app *gin.Engine) {
	authRouter := app.Group("/auth")
	{
		authRouter.POST("/register", r.authHandlers.Register)
		authRouter.POST("/login", middleware.RequireAuthMiddleware, r.authHandlers.Login)
		authRouter.POST("/forgotPassword", r.authHandlers.ForgotPassword)
		authRouter.POST("/resetPassword", r.authHandlers.ResetPassword)
		authRouter.POST("/checkVerificationCode", r.authHandlers.CheckCode)
		authRouter.POST("/logout", middleware.RequireAuthMiddleware, r.authHandlers.Logout)
		authRouter.DELETE("/deleteAccount", middleware.RequireAuthMiddleware, r.authHandlers.DeleteAccount)
		authRouter.GET("/profile", middleware.RequireAuthMiddleware, r.authHandlers.Profile)
	}
	gameRouter := app.Group("/game", middleware.RequireAuthMiddleware)
	{
		gameRouter.POST("/add-hero-to-deck", r.gameHandlers.AddHeroToDeck)
		gameRouter.POST("/delete-hero-to-deсk", r.gameHandlers.DeleteHeroToDeсk)
		gameRouter.GET("/get-my-heros", r.gameHandlers.GetMyHeros)
		gameRouter.POST("/add-spell-to-deck", r.gameHandlers.AddSpellToDeck)
		gameRouter.POST("/delete-spell-to-deсk", r.gameHandlers.DeleteSpellToDeсk)
		gameRouter.GET("/get-my-spells", r.gameHandlers.GetMySpell)
		gameRouter.POST("/hero/:id", r.gameHandlers.BuyHero)
		gameRouter.POST("/spell/:id", r.gameHandlers.BuySpell)
		gameRouter.POST("/create-hero", r.gameHandlers.CreateHero)
		gameRouter.POST("/create-spell", r.gameHandlers.CreateSpell)
		gameRouter.GET("/get-all-heros", r.gameHandlers.GetAllHeros)
		gameRouter.GET("/get-all-spells", r.gameHandlers.GetAllSpell)
		gameRouter.GET("/get-my-deсk/:id", r.gameHandlers.GetMyDeck)
		gameRouter.GET("/hero/:id", r.gameHandlers.GetHero)
		gameRouter.GET("/spell/:id", r.gameHandlers.GetSpell)
	}
}
