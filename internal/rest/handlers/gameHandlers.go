package handlers

import (
	"auth/internal/repository"
	"auth/internal/rest/models"
	"auth/pkg/logger"
	"auth/pkg/rest/helper"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type GameHandlers struct {
	UserRepo repository.UserRepo
	GameRepo repository.GameRepository
}

func NewGameHandlers(userRepo repository.UserRepo, gameRepo repository.GameRepository) *GameHandlers {
	return &GameHandlers{UserRepo: userRepo, GameRepo: gameRepo}
}

func (h GameHandlers) AddHeroToDeck(context *gin.Context) {
	var input struct {
		DeckID uint `json:"deck_id"`
		HeroID uint `json:"hero_id"`
	}

	if err := context.ShouldBindJSON(&input); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
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

	user, err := h.UserRepo.GetUserByEmail(emailm)
	if err != nil {
		logger.GetLogger().Error("User not found:", err)
		context.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "User not found"})
		return
	}

	heroBought, err := h.GameRepo.HasUserBoughtHero(user.ID, input.HeroID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check if user has bought hero"})
		return
	}

	if !heroBought {
		context.JSON(http.StatusBadRequest, gin.H{"error": "User has not bought this hero"})
		return
	}
	deck, err := h.GameRepo.GetDeckByID(input.DeckID)
	if err != nil {
		userID := user.ID
		newDeck := models.Deck{
			UserID: userID,
		}

		err := h.GameRepo.CreateDeck(&newDeck)
		if err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create deck"})
			return
		}

		deck, err = h.GameRepo.GetDeckByID(newDeck.ID)
		if err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get created deck"})
			return
		}
	}

	updatedDecks, err := h.GameRepo.AddHeroToDeck(deck.ID, input.HeroID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add hero to deck"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"decks": updatedDecks})
}

func (h GameHandlers) DeleteHeroToDeсk(context *gin.Context) {
	deckID, err := strconv.ParseUint(context.Param("deckID"), 10, 64)
	if err != nil {
		logger.GetLogger().Error("Invalid deck ID format:", err)

		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid deck ID format"})
		return
	}
	heroID, err := strconv.ParseUint(context.Param("id"), 10, 64)
	if err != nil {
		logger.GetLogger().Error("Invalid spell ID format:", err)

		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid spell ID format"})
		return
	}
	updatedDecks, err := h.GameRepo.DeleteHeroFromDeck(uint(deckID), uint(heroID))
	if err != nil {
		logger.GetLogger().Error("Failed to delete hero from deck:", err)
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete hero from deck"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Hero successfully delete from deck", "decks": updatedDecks})
}

func (h GameHandlers) GetMyHeros(context *gin.Context) {
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

	user, err := h.UserRepo.GetUserByEmail(emailm)
	if err != nil {
		logger.GetLogger().Error("User not found:", err)
		context.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "User not found"})
		return
	}

	heros, err := h.GameRepo.GetMyHeros(user.ID)
	if err != nil {
		logger.GetLogger().Error("Failed to get user heros:", err)
		context.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to get user heros",
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{"heros": heros})
}

func (h GameHandlers) AddSpellToDeck(context *gin.Context) {
	deckID, err := strconv.ParseUint(context.Param("deckID"), 10, 64)
	if err != nil {
		logger.GetLogger().Error("Invalid deck ID format:", err)

		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid deck ID format"})
		return
	}
	spellID, err := strconv.ParseUint(context.Param("id"), 10, 64)
	if err != nil {
		logger.GetLogger().Error("Invalid spell ID format:", err)

		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid spell ID format"})
		return
	}
	updatedDecks, err := h.GameRepo.AddSpellToDeck(uint(deckID), uint(spellID))
	if err != nil {
		logger.GetLogger().Error("Failed to adding spell from deck:", err)

		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to adding spell from deck"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Spell successfully added from deck", "decks": updatedDecks})
}

func (h GameHandlers) DeleteSpellToDeсk(context *gin.Context) {
	deckID, err := strconv.ParseUint(context.Param("deckID"), 10, 64)
	if err != nil {
		logger.GetLogger().Error("Invalid deck ID format:", err)

		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid deck ID format"})
		return
	}
	spellID, err := strconv.ParseUint(context.Param("id"), 10, 64)
	if err != nil {
		logger.GetLogger().Error("Invalid spell ID format:", err)

		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid spell ID format"})
		return
	}
	updatedDecks, err := h.GameRepo.DeleteSpellFromDeck(uint(deckID), uint(spellID))
	if err != nil {
		logger.GetLogger().Error("Failed to delete spell from deck:", err)

		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete spell from deck"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Spell successfully deleted from deck", "decks": updatedDecks})
}

func (h GameHandlers) GetMySpell(context *gin.Context) {
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

	user, err := h.UserRepo.GetUserByEmail(emailm)
	if err != nil {
		logger.GetLogger().Error("User not found:", err)
		context.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "User not found"})
		return
	}

	spells, err := h.GameRepo.GetMySpells(user.ID)
	if err != nil {
		logger.GetLogger().Error("Failed to get user spells:", err)

		context.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to get user spells",
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{"spells": spells})
}

func (h GameHandlers) BuyHero(context *gin.Context) {
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

	user, err := h.UserRepo.GetUserByEmail(emailm)
	if err != nil {
		logger.GetLogger().Error("User not found:", err)
		context.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "User not found"})
		return
	}
	heroID, err := strconv.ParseUint(context.Param("id"), 10, 64)
	if err != nil {
		logger.GetLogger().Error("Invalid hero ID format:", err)
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid hero ID format"})
		return
	}
	hero, err := h.GameRepo.GetHeroByID(uint(heroID))
	if err != nil {
		logger.GetLogger().Error("Failed to get hero by ID:", err)
		context.JSON(http.StatusNotFound, gin.H{"error": "Hero not found"})
		return
	}
	if user.Bank < hero.Price {
		logger.GetLogger().Error("Insufficient balance to buy the hero")
		context.JSON(http.StatusForbidden, gin.H{"status": "error", "message": "Insufficient balance to buy the hero"})
		return
	}
	user.Bank -= hero.Price
	if err := h.GameRepo.AddHeroToUser(user.ID, hero.ID); err != nil {
		logger.GetLogger().Error("Failed to add hero to user's collection:", err)
		context.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to add hero to user's collection"})
		return
	}
	if err := h.UserRepo.UpdateUserBalance(user.ID, user.Bank); err != nil {
		logger.GetLogger().Error("Failed to update user's balance:", err)
		context.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to update user's balance"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"status": "success", "message": "Hero bought successfully"})
}

func (h GameHandlers) BuySpell(context *gin.Context) {
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

	user, err := h.UserRepo.GetUserByEmail(emailm)
	if err != nil {
		logger.GetLogger().Error("User not found:", err)
		context.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "User not found"})
		return
	}
	spellID, err := strconv.ParseUint(context.Param("id"), 10, 64)
	if err != nil {
		logger.GetLogger().Error("Invalid spell ID format:", err)
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid spell ID format"})
		return
	}

	spell, err := h.GameRepo.GetSpellByID(uint(spellID))
	if err != nil {
		logger.GetLogger().Error("Failed to get spell by ID:", err)
		context.JSON(http.StatusNotFound, gin.H{"error": "Spell not found"})
		return
	}
	if user.Bank < spell.Price {
		logger.GetLogger().Info("Insufficient balance to buy the spell")
		context.JSON(http.StatusForbidden, gin.H{"status": "error", "message": "Insufficient balance to buy the spell"})
		return
	}
	user.Bank -= spell.Price
	if err := h.UserRepo.UpdateUserBalance(user.ID, user.Bank); err != nil {
		logger.GetLogger().Error("Failed to update user's balance:", err)
		context.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to update user's balance"})
		return
	}
	if err := h.GameRepo.AddSpellToUser(user.ID, spell.ID); err != nil {
		logger.GetLogger().Error("Failed to add spell to user's collection:", err)
		context.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to add spell to user's collection"})
		return
	}
	context.JSON(http.StatusOK, gin.H{"status": "success", "message": "Spell purchased successfully"})
}
func (h GameHandlers) CreateHero(context *gin.Context) {
	logger.GetLogger().Info("Creating hero")

	email, exists := context.Get("email")
	if !exists {
		logger.GetLogger().Error("User not authenticated")
		context.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "User not authenticated"})
		return
	}

	emailm, ok := email.(string)
	if !ok {
		logger.GetLogger().Error("Failed to get user ID from context")
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	user, err := h.UserRepo.GetUserByEmail(emailm)
	if err != nil {
		logger.GetLogger().Error("Failed to get user ID from context")
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	if user.UserType != "ADMIN" {
		logger.GetLogger().Warn("Not enough rights to create hero")
		context.JSON(http.StatusForbidden, gin.H{"error": "Not enough rights to act"})
		return
	}

	var hero models.Hero
	if err := context.BindJSON(&hero); err != nil {
		logger.GetLogger().Error("Failed to bind JSON for hero creation:", err)
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}

	if err := h.GameRepo.CreateHero(&hero); err != nil {
		logger.GetLogger().Error("Failed to create hero:", err)
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	context.JSON(http.StatusCreated, gin.H{"hero": hero})
}

func (h GameHandlers) CreateSpell(context *gin.Context) {
	logger.GetLogger().Info("Creating spell")

	email, exists := context.Get("email")
	if !exists {
		logger.GetLogger().Error("User not authenticated")
		context.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "User not authenticated"})
		return
	}

	emailm, ok := email.(string)
	if !ok {
		logger.GetLogger().Error("Failed to get user ID from context")
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	user, err := h.UserRepo.GetUserByEmail(emailm)
	if err != nil {
		logger.GetLogger().Error("Failed to get user ID from context")
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	if user.UserType != "ADMIN" {
		logger.GetLogger().Warn("Not enough rights to create spell")
		context.JSON(http.StatusForbidden, gin.H{"error": "Not enough rights to act"})
		return
	}

	var spell models.Spell
	if err := context.BindJSON(&spell); err != nil {
		logger.GetLogger().Error("Failed to bind JSON for spell creation:", err)
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}

	if err := h.GameRepo.CreateSpell(&spell); err != nil {
		logger.GetLogger().Error("Failed to create spell:", err)
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	context.JSON(http.StatusCreated, gin.H{"spell": spell})
}
func (h GameHandlers) GetAllHeros(context *gin.Context) {
	logger.GetLogger().Info("Fetching all heroes")

	sortBy := context.DefaultQuery("sortBy", "")
	sortOrder := context.DefaultQuery("sortOrder", "")
	filterName := context.DefaultQuery("filterName", "")
	page, err := strconv.Atoi(context.DefaultQuery("page", "1"))
	if err != nil {
		logger.GetLogger().Error("Invalid page parameter:", err)
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page parameter"})
		return
	}
	pageSize, err := strconv.Atoi(context.DefaultQuery("pageSize", "10"))
	if err != nil {
		logger.GetLogger().Error("Invalid pageSize parameter:", err)
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pageSize parameter"})
		return
	}

	heros, err := h.GameRepo.GetAllHeros(sortBy, sortOrder, filterName, page, pageSize)
	if err != nil {
		logger.GetLogger().Error("Failed to get all heroes:", err)
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"heroes": heros})
}

func (h GameHandlers) GetAllSpell(context *gin.Context) {
	logger.GetLogger().Info("Fetching all spells")

	sortBy := context.DefaultQuery("sortBy", "")
	sortOrder := context.DefaultQuery("sortOrder", "")
	filterName := context.DefaultQuery("filterName", "")
	page, err := strconv.Atoi(context.DefaultQuery("page", "1"))
	if err != nil {
		logger.GetLogger().Error("Invalid page parameter:", err)
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page parameter"})
		return
	}
	pageSize, err := strconv.Atoi(context.DefaultQuery("pageSize", "10"))
	if err != nil {
		logger.GetLogger().Error("Invalid pageSize parameter:", err)
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pageSize parameter"})
		return
	}

	spells, err := h.GameRepo.GetAllSpells(sortBy, sortOrder, filterName, page, pageSize)
	if err != nil {
		logger.GetLogger().Error("Failed to get all spells:", err)
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"spells": spells})
}

func (h GameHandlers) GetMyDeck(context *gin.Context) {
	logger.GetLogger().Info("Fetching user's deck")

	email, exists := context.Get("email")
	if !exists {
		logger.GetLogger().Error("User not authenticated")
		context.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "User not authenticated"})
		return
	}

	emailm, ok := email.(string)
	if !ok {
		logger.GetLogger().Error("Failed to get user ID from context")
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	user, err := h.UserRepo.GetUserByEmail(emailm)
	if err != nil {
		logger.GetLogger().Error("Failed to get user ID from context")
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	deckID, err := strconv.ParseUint(context.Param("id"), 10, 64)
	if err != nil {
		logger.GetLogger().Error("Invalid deck ID format:", err)
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid deck ID format"})
		return
	}

	deck, err := h.GameRepo.GetDeckByID(uint(deckID))
	if err != nil {
		if errors.Is(err, helper.ErrDeckNotFound) {
			logger.GetLogger().Warn("Deck not found:", err)
			context.JSON(http.StatusNotFound, gin.H{"error": "Deck not found"})
			return
		}
		logger.GetLogger().Error("Failed to get deck by ID:", err)
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	if deck.UserID != user.ID {
		logger.GetLogger().Warn("Unauthorized access to deck:", user.ID)
		context.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized access"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"deck": deck})
}

func (h GameHandlers) GetHero(context *gin.Context) {
	logger.GetLogger().Info("Fetching hero information")

	heroID, err := strconv.ParseUint(context.Param("id"), 10, 64)
	if err != nil {
		logger.GetLogger().Error("Invalid hero ID format:", err)
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid hero ID format"})
		return
	}

	hero, err := h.GameRepo.GetHeroByID(uint(heroID))
	if err != nil {
		if errors.Is(err, helper.ErrHeroNotFound) {
			logger.GetLogger().Error("Hero not found:", err)
			context.JSON(http.StatusNotFound, gin.H{"error": "Hero not found"})
		} else {
			logger.GetLogger().Error("Failed to get hero by ID:", err)
			context.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		}
		return
	}

	context.JSON(http.StatusOK, gin.H{"hero": hero})
}

func (h GameHandlers) GetSpell(context *gin.Context) {
	logger.GetLogger().Info("Fetching spell information")

	spellID, err := strconv.ParseUint(context.Param("id"), 10, 64)
	if err != nil {
		logger.GetLogger().Error("Invalid spell ID format:", err)
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid spell ID format"})
		return
	}

	spell, err := h.GameRepo.GetSpellByID(uint(spellID))
	if err != nil {
		if errors.Is(err, helper.ErrSpellNotFound) {
			logger.GetLogger().Error("Spell not found:", err)
			context.JSON(http.StatusNotFound, gin.H{"error": "Spell not found"})
		} else {
			logger.GetLogger().Error("Failed to get spell by ID:", err)
			context.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		}
		return
	}
	context.JSON(http.StatusOK, gin.H{"spell": spell})
}
