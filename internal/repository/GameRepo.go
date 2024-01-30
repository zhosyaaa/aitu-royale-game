package repository

import (
	"auth/internal/rest/models"
	"database/sql"
	"errors"
	"fmt"
)

type GameRepo interface {
	GetSpellByID(id uint) (models.Spell, error)
	GetHeroByID(id uint) (models.Hero, error)
	GetDeckByID(id uint) (models.Deck, error)
	CreateSpell(spell *models.Spell) error
	CreateHero(hero *models.Hero) error
	CreateDeck(deck *models.Deck) error
	AddHeroToUser(userID, heroID uint) error
	AddSpellToUser(userID, spellID uint) error
	AddHeroToDeck(deckID, heroID uint) ([]models.Deck, error)
	DeleteHeroFromDeck(deckID, heroID uint) ([]models.Deck, error)
	AddSpellToDeck(deckID, spellID uint) ([]models.Deck, error)
	DeleteSpellFromDeck(deckID, spellID uint) ([]models.Deck, error)
	GetDecksForUser(userID uint) ([]models.Deck, error)

	GetAllSpells(sortBy, sortOrder, filterName string, page, pageSize int) ([]models.Spell, error)
	GetAllHeros(sortBy, sortOrder, filterName string, page, pageSize int) ([]models.Hero, error)
	GetMySpells(userID uint) ([]models.Spell, error)
	GetMyHeros(userID uint) ([]models.Hero, error)
	HasUserBoughtHero(userID, heroID uint) (bool, error)
}

type GameRepository struct {
	db *sql.DB
}

func NewGameRepository(db *sql.DB) *GameRepository {
	return &GameRepository{db}
}

func (repo *GameRepository) HasUserBoughtHero(userID, heroID uint) (bool, error) {
	query := `
        SELECT COUNT(*)
        FROM user_heros
        WHERE user_id = $1 AND hero_id = $2
    `

	var count int
	err := repo.db.QueryRow(query, userID, heroID).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check if user has bought hero: %v", err)
	}

	return count > 0, nil
}
func (repo *GameRepository) GetSpellByID(id uint) (models.Spell, error) {
	query := `
		SELECT * FROM spells WHERE id = $1
	`
	var spell models.Spell
	err := repo.db.QueryRow(query, id).Scan(
		&spell.ID,
		&spell.CreatedAt,
		&spell.UpdatedAt,
		&spell.DeletedAt,
		&spell.Name,
		&spell.Description,
		&spell.Area,
		&spell.DamageType,
		&spell.Damage,
		&spell.Duration,
		&spell.Effect,
		&spell.Price,
	)
	if err == sql.ErrNoRows {
		return models.Spell{}, errors.New("spell not found")
	} else if err != nil {
		return models.Spell{}, fmt.Errorf("failed to get spell by ID: %v", err)
	}
	return spell, nil
}
func (repo *GameRepository) GetHeroByID(id uint) (models.Hero, error) {
	query := `
		SELECT * FROM heros WHERE id = $1
	`
	var hero models.Hero
	err := repo.db.QueryRow(query, id).Scan(
		&hero.ID,
		&hero.CreatedAt,
		&hero.UpdatedAt,
		&hero.DeletedAt,
		&hero.Name,
		&hero.Description,
		&hero.Rarity,
		&hero.DamageType,
		&hero.Effect,
		&hero.Hitpoint,
		&hero.Damage,
		&hero.CostElixir,
		&hero.DamageTower,
		&hero.Speed,
		&hero.Price,
	)
	if err == sql.ErrNoRows {
		return models.Hero{}, errors.New("hero not found")
	} else if err != nil {
		return models.Hero{}, fmt.Errorf("failed to get hero by ID: %v", err)
	}
	return hero, nil
}
func (repo *GameRepository) GetDeckByID(id uint) (models.Deck, error) {
	query := `
		SELECT * FROM decks WHERE id = $1
	`

	var deck models.Deck
	err := repo.db.QueryRow(query, id).Scan(
		&deck.ID,
		&deck.CreatedAt,
		&deck.UpdatedAt,
		&deck.DeletedAt,
		&deck.UserID,
		&deck.Name,
		&deck.Description,
	)

	if err == sql.ErrNoRows {
		return models.Deck{}, errors.New("deck not found")
	} else if err != nil {
		return models.Deck{}, fmt.Errorf("failed to get deck by ID: %v", err)
	}
	return deck, nil
}
func (repo *GameRepository) CreateSpell(spell *models.Spell) error {
	query := `
		INSERT INTO spells (name, description, area, damage_type, damage, duration, effect, price)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id
	`

	err := repo.db.QueryRow(
		query,
		spell.Name,
		spell.Description,
		spell.Area,
		spell.DamageType,
		spell.Damage,
		spell.Duration,
		spell.Effect,
		spell.Price,
	).Scan(&spell.ID)
	if err != nil {
		return fmt.Errorf("failed to create spell: %v", err)
	}
	return nil
}

func (repo *GameRepository) CreateHero(hero *models.Hero) error {
	query := `
		INSERT INTO heros (name, description, rarity, damage_type, effect, hitpoint, damage, cost_elixir, damage_tower, speed, price)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id
	`
	err := repo.db.QueryRow(
		query,
		hero.Name,
		hero.Description,
		hero.Rarity,
		hero.DamageType,
		hero.Effect,
		hero.Hitpoint,
		hero.Damage,
		hero.CostElixir,
		hero.DamageTower,
		hero.Speed,
		hero.Price,
	).Scan(&hero.ID)
	if err != nil {
		return fmt.Errorf("failed to create hero: %v", err)
	}
	return nil
}
func (repo *GameRepository) AddHeroToUser(userID, heroID uint) error {
	query := `
		INSERT INTO user_heros (user_id, hero_id)
		VALUES ($1, $2)
	`
	_, err := repo.db.Exec(query, userID, heroID)
	if err != nil {
		return fmt.Errorf("failed to add hero to user: %v", err)
	}
	return nil
}
func (repo *GameRepository) AddSpellToUser(userID, spellID uint) error {
	query := `
		INSERT INTO user_spells (user_id, spell_id)
		VALUES ($1, $2)
	`

	_, err := repo.db.Exec(query, userID, spellID)
	if err != nil {
		return fmt.Errorf("failed to add spell to user: %v", err)
	}

	return nil
}
func (repo *GameRepository) AddHeroToDeck(deckID, heroID uint) ([]models.Deck, error) {
	existingDeck, err := repo.GetDeckByID(deckID)
	if err != nil {
		return nil, fmt.Errorf("failed to get deck by ID: %v", err)
	}

	_, err = repo.GetHeroByID(heroID)
	if err != nil {
		return nil, fmt.Errorf("failed to get hero by ID: %v", err)
	}

	query := `
		INSERT INTO deck_heros (deck_id, hero_id)
		VALUES ($1, $2)
		RETURNING id
	`

	var deck models.Deck
	err = repo.db.QueryRow(query, deckID, heroID).Scan(&deck.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to add hero to deck: %v", err)
	}

	updatedDecks, err := repo.GetDecksForUser(existingDeck.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get updated decks: %v", err)
	}

	return updatedDecks, nil
}

func (repo *GameRepository) GetDecksForUser(userID uint) ([]models.Deck, error) {
	query := `
		SELECT * FROM decks WHERE user_id = $1
	`

	rows, err := repo.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get decks for user: %v", err)
	}
	defer rows.Close()

	var decks []models.Deck
	for rows.Next() {
		var deck models.Deck
		err := rows.Scan(
			&deck.ID,
			&deck.CreatedAt,
			&deck.UpdatedAt,
			&deck.DeletedAt,
			&deck.UserID,
			&deck.Name,
			&deck.Description,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan deck: %v", err)
		}
		decks = append(decks, deck)
	}

	return decks, nil
}
func (repo *GameRepository) DeleteHeroFromDeck(deckID, heroID uint) ([]models.Deck, error) {
	existingDeck, err := repo.GetDeckByID(deckID)
	if err != nil {
		return nil, fmt.Errorf("failed to get deck by ID: %v", err)
	}

	_, err = repo.GetHeroByID(heroID)
	if err != nil {
		return nil, fmt.Errorf("failed to get hero by ID: %v", err)
	}

	query := `
		DELETE FROM deck_heros
		WHERE deck_id = $1 AND hero_id = $2
	`

	_, err = repo.db.Exec(query, deckID, heroID)
	if err != nil {
		return nil, fmt.Errorf("failed to delete hero from deck: %v", err)
	}

	updatedDecks, err := repo.GetDecksForUser(existingDeck.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get updated decks: %v", err)
	}

	return updatedDecks, nil
}
func (repo *GameRepository) AddSpellToDeck(deckID, spellID uint) ([]models.Deck, error) {
	existingDeck, err := repo.GetDeckByID(deckID)
	if err != nil {
		return nil, fmt.Errorf("failed to get deck by ID: %v", err)
	}

	_, err = repo.GetSpellByID(spellID)
	if err != nil {
		return nil, fmt.Errorf("failed to get spell by ID: %v", err)
	}

	query := `
		INSERT INTO deck_spells (deck_id, spell_id)
		VALUES ($1, $2)
		RETURNING id
	`
	var deck models.Deck
	err = repo.db.QueryRow(query, deckID, spellID).Scan(&deck.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to add spell to deck: %v", err)
	}
	updatedDecks, err := repo.GetDecksForUser(existingDeck.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get updated decks: %v", err)
	}

	return updatedDecks, nil
}

func (repo *GameRepository) DeleteSpellFromDeck(deckID, spellID uint) ([]models.Deck, error) {
	existingDeck, err := repo.GetDeckByID(deckID)
	if err != nil {
		return nil, fmt.Errorf("failed to get deck by ID: %v", err)
	}

	_, err = repo.GetSpellByID(spellID)
	if err != nil {
		return nil, fmt.Errorf("failed to get spell by ID: %v", err)
	}

	query := `
		DELETE FROM deck_spells
		WHERE deck_id = $1 AND spell_id = $2
	`

	_, err = repo.db.Exec(query, deckID, spellID)
	if err != nil {
		return nil, fmt.Errorf("failed to delete spell from deck: %v", err)
	}

	updatedDecks, err := repo.GetDecksForUser(existingDeck.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get updated decks: %v", err)
	}

	return updatedDecks, nil
}
func (repo *GameRepository) GetAllSpells(sortBy, sortOrder, filterName string, page, pageSize int) ([]models.Spell, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	query := `
        SELECT * FROM spells
        WHERE ($1 = '' OR name ILIKE $1)
        ORDER BY %s %s
        LIMIT $2 OFFSET $3
    `

	if sortBy == "" || (sortBy != "name" && sortBy != "damage") {
		sortBy = "id"
	}
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "asc"
	}

	query = fmt.Sprintf(query, sortBy, sortOrder)

	rows, err := repo.db.Query(query, "%"+filterName+"%", pageSize, (page-1)*pageSize)
	if err != nil {
		return nil, fmt.Errorf("failed to get all spells: %v", err)
	}
	defer rows.Close()

	var spells []models.Spell
	for rows.Next() {
		var spell models.Spell
		err := rows.Scan(
			&spell.ID,
			&spell.CreatedAt,
			&spell.UpdatedAt,
			&spell.DeletedAt,
			&spell.Name,
			&spell.Description,
			&spell.Area,
			&spell.DamageType,
			&spell.Damage,
			&spell.Duration,
			&spell.Effect,
			&spell.Price, // Add this line for the new 'price' field
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan spell: %v", err)
		}
		spells = append(spells, spell)
	}

	return spells, nil
}

func (repo *GameRepository) GetAllHeros(sortBy, sortOrder, filterName string, page, pageSize int) ([]models.Hero, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	query := `
        SELECT * FROM heros
        WHERE ($1 = '' OR name ILIKE $1)
        ORDER BY %s %s
        LIMIT $2 OFFSET $3
    `

	if sortBy == "" || (sortBy != "name" && sortBy != "damage" && sortBy != "speed") {
		sortBy = "id"
	}
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "asc"
	}

	query = fmt.Sprintf(query, sortBy, sortOrder)

	rows, err := repo.db.Query(query, "%"+filterName+"%", pageSize, (page-1)*pageSize)
	if err != nil {
		return nil, fmt.Errorf("failed to get all heros: %v", err)
	}
	defer rows.Close()

	var heros []models.Hero
	for rows.Next() {
		var hero models.Hero
		err := rows.Scan(
			&hero.ID,
			&hero.CreatedAt,
			&hero.UpdatedAt,
			&hero.DeletedAt,
			&hero.Name,
			&hero.Description,
			&hero.Rarity,
			&hero.DamageType,
			&hero.Effect,
			&hero.Hitpoint,
			&hero.Damage,
			&hero.CostElixir,
			&hero.DamageTower,
			&hero.Speed,
			&hero.Price,
		)
		fmt.Println("test:", hero.ID, hero.DeletedAt, hero.Name)
		fmt.Errorf("%e,errrorrrrrr", err)
		if err != nil {
			return nil, fmt.Errorf("failed to scan hero: %v", err)
		}
		heros = append(heros, hero)
	}

	return heros, nil
}

func (repo *GameRepository) CreateDeck(deck *models.Deck) error {
	query := `
        INSERT INTO decks (user_id, name, description)
        VALUES ($1, $2, $3)
        RETURNING id
    `

	err := repo.db.QueryRow(
		query,
		deck.UserID,
		deck.Name,
		deck.Description,
	).Scan(&deck.ID)
	if err != nil {
		return fmt.Errorf("failed to create deck: %v", err)
	}

	return nil
}

func (repo *GameRepository) GetMySpells(userID uint) ([]models.Spell, error) {
	query := `
		SELECT s.id, s.created_at, s.updated_at, s.deleted_at, s.name, s.description, s.area, s.damage_type, s.damage, s.duration, s.effect
		FROM spells s
		JOIN user_spells us ON s.id = us.spell_id
		WHERE us.user_id = $1
	`

	rows, err := repo.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get my spells: %v", err)
	}
	defer rows.Close()

	var spells []models.Spell
	for rows.Next() {
		var spell models.Spell
		err := rows.Scan(
			&spell.ID,
			&spell.CreatedAt,
			&spell.UpdatedAt,
			&spell.DeletedAt,
			&spell.Name,
			&spell.Description,
			&spell.Area,
			&spell.DamageType,
			&spell.Damage,
			&spell.Duration,
			&spell.Effect,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan spell: %v", err)
		}
		spells = append(spells, spell)
	}

	return spells, nil
}

func (repo *GameRepository) GetMyHeros(userID uint) ([]models.Hero, error) {
	query := `
		SELECT h.* FROM heros h
		JOIN user_heros uh ON h.id = uh.hero_id
		WHERE uh.user_id = $1
	`

	rows, err := repo.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get my heros: %v", err)
	}
	defer rows.Close()

	var myHeros []models.Hero
	for rows.Next() {
		var hero models.Hero
		err := rows.Scan(
			&hero.ID,
			&hero.CreatedAt,
			&hero.UpdatedAt,
			&hero.DeletedAt,
			&hero.Name,
			&hero.Description,
			&hero.Rarity,
			&hero.DamageType,
			&hero.Effect,
			&hero.Hitpoint,
			&hero.Damage,
			&hero.CostElixir,
			&hero.DamageTower,
			&hero.Speed,
			&hero.Price,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan hero: %v", err)
		}
		myHeros = append(myHeros, hero)
	}

	return myHeros, nil
}
