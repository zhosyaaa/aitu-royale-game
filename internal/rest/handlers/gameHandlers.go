package handlers

import (
	"auth/internal/repository"
	"github.com/gin-gonic/gin"
	"net/http"
)

type GameHandlers struct {
	GameRepo repository.GameRepository
}

func NewGameHandlers(gameRepo repository.GameRepository) *GameHandlers {
	return &GameHandlers{GameRepo: gameRepo}
}

func (h GameHandlers) AddHeroToDeck(context *gin.Context) {

}

func (h GameHandlers) DeleteHeroToDeсk(context *gin.Context) {

}

func (h GameHandlers) GetMyHeros(context *gin.Context) {

}

func (h GameHandlers) AddSpellToDeck(context *gin.Context) {

}

func (h GameHandlers) DeleteSpellToDeсk(context *gin.Context) {

}

func (h GameHandlers) GetMySpell(context *gin.Context) {

}

func (h GameHandlers) BuyHero(context *gin.Context) {

}

func (h GameHandlers) BuySpell(context *gin.Context) {

}

func (h GameHandlers) CreateHero(context *gin.Context) {
	userType, exists := context.Get("userType")
	if !exists {
		context.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "User not authenticated",
		})
		return
	}

	userTypem, ok := userType.(string)
	if !ok {
		context.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Error while retrieving user ID",
		})
		return
	}
	if userTypem != "ADMIN" {
		context.JSON(http.StatusForbidden, gin.H{
			"status":  "error",
			"message": "Not enough rights to act",
		})
		return
	}
}

func (h GameHandlers) CreateSpell(context *gin.Context) {
	userType, exists := context.Get("userType")
	if !exists {
		context.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "User not authenticated",
		})
		return
	}

	userTypem, ok := userType.(string)
	if !ok {
		context.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Error while retrieving user ID",
		})
		return
	}
	if userTypem != "ADMIN" {
		context.JSON(http.StatusForbidden, gin.H{
			"status":  "error",
			"message": "Not enough rights to act",
		})
		return
	}
}

func (h GameHandlers) GetAllHeros(context *gin.Context) {

}

func (h GameHandlers) GetAllSpell(context *gin.Context) {

}

func (h GameHandlers) GetMyDeck(context *gin.Context) {

}

func (h GameHandlers) GetHero(context *gin.Context) {

}

func (h GameHandlers) GetSpell(context *gin.Context) {

}

func (h GameHandlers) GetRandomDeck(context *gin.Context) {

}
