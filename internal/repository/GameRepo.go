package repository

import (
	"auth/internal/rest/models"
	"database/sql"
)

type GameRepo interface {
	GetSpellByID(id uint) (models.Spell, error)
	GetHeroByID(id uint) (models.Hero, error)
	GetDeckByID(id uint) (models.Deck, error)
	CreateSpell(spell *models.Spell) error
	CreateHero(hero *models.Hero) error
	AddHeroToUser(userID, heroID uint) error
	AddSpellToUser(userID, spellID uint) error
	AddHeroToDeck(deckID, heroID uint) ([]models.Deck, error)
	DeleteHeroFromDeck(deckID, heroID uint) ([]models.Deck, error)
	AddSpellToDeck(deckID, spellID uint) ([]models.Deck, error)
	DeleteSpellFromDeck(deckID, spellID uint) ([]models.Deck, error)
	GetRandomDeck(userID uint) ([]models.Deck, error)

	GetAllSpells() ([]models.Spell, error)
	GetAllHeros() ([]models.Hero, error)
	GetMySpells(userID uint) ([]models.Spell, error)
	GetMyHeros(userID uint) ([]models.Hero, error)
}

type GameRepository struct {
	db *sql.DB
}

func NewGameRepository(db *sql.DB) *GameRepository {
	return &GameRepository{db}
}

func (r *GameRepository) GetSpellByID(id uint) (models.Spell, error) {
	var spell models.Spell
	err := r.db.QueryRow("SELECT * FROM spells WHERE id = $1", id).
		Scan(&spell.ID, &spell.CreatedAt, &spell.UpdatedAt, &spell.DeletedAt, &spell.Name, &spell.Description, &spell.Area, &spell.DamageType, &spell.Damage, &spell.Duration, &spell.Effect)
	if err != nil {
		return models.Spell{}, err
	}
	return spell, nil
}

func (r *GameRepository) GetHeroByID(id uint) (models.Hero, error) {
	var hero models.Hero
	err := r.db.QueryRow("SELECT * FROM heros WHERE id = $1", id).
		Scan(&hero.ID, &hero.CreatedAt, &hero.UpdatedAt, &hero.DeletedAt, &hero.Name, &hero.Description, &hero.Rarity, &hero.DamageType, &hero.Effect, &hero.Hitpoint, &hero.Damage, &hero.CostElixir, &hero.DamageTower, &hero.Speed, &hero.Price)
	if err != nil {
		return models.Hero{}, err
	}
	return hero, nil
}

func (r *GameRepository) GetDeckByID(id uint) (models.Deck, error) {
	var deck models.Deck
	err := r.db.QueryRow("SELECT * FROM decks WHERE id = $1", id).
		Scan(&deck.ID, &deck.CreatedAt, &deck.UpdatedAt, &deck.DeletedAt, &deck.UserID, &deck.Name, &deck.Description)
	if err != nil {
		return models.Deck{}, err
	}
	return deck, nil
}

func (r *GameRepository) CreateSpell(spell *models.Spell) error {
	_, err := r.db.Exec("INSERT INTO spells (name, description, area, damage_type, damage, duration, effect) VALUES ($1, $2, $3, $4, $5, $6, $7)",
		spell.Name, spell.Description, spell.Area, spell.DamageType, spell.Damage, spell.Duration, spell.Effect)
	if err != nil {
		return err
	}
	return nil
}

func (r *GameRepository) CreateHero(hero *models.Hero) error {
	_, err := r.db.Exec("INSERT INTO heros (name, description, rarity, damage_type, effect, hitpoint, damage, cost_elixir, damage_tower, speed, price) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)",
		hero.Name, hero.Description, hero.Rarity, hero.DamageType, hero.Effect, hero.Hitpoint, hero.Damage, hero.CostElixir, hero.DamageTower, hero.Speed, hero.Price)
	if err != nil {
		return err
	}
	return nil
}

func (r *GameRepository) AddHeroToUser(userID, heroID uint) error {
	_, err := r.db.Exec("INSERT INTO user_heros (user_id, hero_id) VALUES ($1, $2)", userID, heroID)
	if err != nil {
		return err
	}
	return nil
}

func (r *GameRepository) AddSpellToUser(userID, spellID uint) error {
	_, err := r.db.Exec("INSERT INTO user_spells (user_id, spell_id) VALUES ($1, $2)", userID, spellID)
	if err != nil {
		return err
	}
	return nil
}

//func (r *GameRepository) AddHeroToDeck(deckID, heroID uint) ([]models.Deck, error) {
//	deck, err := r.GetDeckByID(deckID)
//	if err != nil {
//		return nil, err
//	}
//
//	hero, err := r.GetHeroByID(heroID)
//	if err != nil {
//		return nil, err
//	}
//
//	userID, err := r.getUserIDByHeroID(heroID)
//	if err != nil || userID != deck.UserID {
//		return nil, errors.New("hero does not belong to the user")
//	}
//
//	_, err = r.db.Exec("INSERT INTO deck_heros (deck_id, hero_id) VALUES ($1, $2)", deckID, heroID)
//	if err != nil {
//		return nil, err
//	}
//
//	return r.getUserDecks(userID)
//}

//// DeleteHeroFromDeck deletes a hero from a deck.
//func (r *GameRepository) DeleteHeroFromDeck(deckID, heroID uint) ([]models.Deck, error) {
//	// Проверяем, существует ли колода
//	deck, err := r.GetDeckByID(deckID)
//	if err != nil {
//		return nil, err
//	}
//
//	// Проверяем, существует ли герой
//	hero, err := r.GetHeroByID(heroID)
//	if err != nil {
//		return nil, err
//	}
//
//	// Проверяем, принадлежит ли герой тому же пользователю
//	userID, err := r.getUserIDByHeroID(heroID)
//	if err != nil || userID != deck.UserID {
//		return nil, errors.New("hero does not belong to the user")
//	}
//
//	// Удаляем героя из колоды
//	_, err = r.db.Exec("DELETE FROM deck_heros WHERE deck_id = $1 AND hero_id = $2", deckID, heroID)
//	if err != nil {
//		return nil, err
//	}
//
//	// Возвращаем обновленный список колод пользователя
//	return r.getUserDecks(userID)
//}

// AddSpellToDeck adds a spell to a deck.
//func (r *GameRepository) AddSpellToDeck(deckID, spellID uint) ([]models.Deck, error) {
//	deck, err := r.GetDeckByID(deckID)
//	if err != nil {
//		return nil, err
//	}
//
//	spell, err := r.GetSpellByID(spellID)
//	if err != nil {
//		return nil, err
//	}
//
//	userID, err := r.getUserIDBySpellID(spellID)
//	if err != nil || userID != deck.UserID {
//		return nil, errors.New("spell does not belong to the user")
//	}
//
//	_, err = r.db.Exec("INSERT INTO deck_spells (deck_id, spell_id) VALUES ($1, $2)", deckID, spellID)
//	if err != nil {
//		return nil, err
//	}
//
//	// Возвращаем обновленный список колод пользователя
//	return r.getUserDecks(userID)
//}

// DeleteSpellFromDeck deletes a spell from a deck.
//func (r *GameRepository) DeleteSpellFromDeck(deckID, spellID uint) ([]models.Deck, error) {
//	// Проверяем, существует ли колода
//	deck, err := r.GetDeckByID(deckID)
//	if err != nil {
//		return nil, err
//	}
//
//	// Проверяем, существует ли заклинание
//	spell, err := r.GetSpellByID(spellID)
//	if err != nil {
//		return nil, err
//	}
//
//	// Проверяем, принадлежит ли заклинание тому же пользователю
//	userID, err := r.getUserIDBySpellID(spellID)
//	if err != nil || userID != deck.UserID {
//		return nil, errors.New("spell does not belong to the user")
//	}
//
//	// Удаляем заклинание из колоды
//	_, err = r.db.Exec("DELETE FROM deck_spells WHERE deck_id = $1 AND spell_id = $2", deckID, spellID)
//	if err != nil {
//		return nil, err
//	}
//
//	// Возвращаем обновленный список колод пользователя
//	return r.getUserDecks(userID)
//}
//
//// GetRandomDeck retrieves a random deck for a user.
//func (r *GameRepository) GetRandomDeck(userID uint) ([]models.Deck, error) {
//	// Получаем все колоды пользователя
//	decks, err := r.getUserDecks(userID)
//	if err != nil {
//		return nil, err
//	}
//
//	// Если у пользователя нет колод, возвращаем пустой список
//	if len(decks) == 0 {
//		return []models.Deck{}, nil
//	}
//
//	// Получаем случайную колоду из списка
//	randomDeck := decks[0] // Пока что просто берем первую колоду, но можно реализовать случайный выбор
//
//	return []models.Deck{randomDeck}, nil
//}

// GetAllSpells retrieves all spells.
func (r *GameRepository) GetAllSpells() ([]models.Spell, error) {
	rows, err := r.db.Query("SELECT * FROM spells")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var spells []models.Spell
	for rows.Next() {
		var spell models.Spell
		err := rows.Scan(
			&spell.ID, &spell.CreatedAt, &spell.UpdatedAt, &spell.DeletedAt,
			&spell.Name, &spell.Description, &spell.Area, &spell.DamageType, &spell.Damage, &spell.Duration, &spell.Effect,
		)
		if err != nil {
			return nil, err
		}
		spells = append(spells, spell)
	}

	return spells, nil
}
