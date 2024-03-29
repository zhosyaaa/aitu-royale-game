package models

import (
	"github.com/lib/pq"
	"time"
)

type User struct {
	ID        uint      `json:"ID,omitempty"`
	CreatedAt time.Time `json:"CreatedAt"`
	UpdatedAt time.Time `json:"UpdatedAt"`
	DeletedAt time.Time `json:"DeletedAt"`
	Username  string    `json:"username" json:"Username,omitempty"`
	Email     string    `json:"email" json:"Email,omitempty"`
	Password  string    `json:"password" json:"Password,omitempty"`
	Heros     []Hero    `json:"heros" json:"Heros,omitempty"`
	Spells    []Spell   `json:"spells" json:"Spells,omitempty"`
	Deck      []Deck    `json:"deck" json:"Deck,omitempty"`
	Bank      int64     `json:"bank" json:"Bank,omitempty"`
	Awards    int32     `json:"awards" json:"Awards,omitempty"`
	UserType  string    `json:"userType"`
}

type Deck struct {
	ID          uint        `json:"ID,omitempty"`
	CreatedAt   time.Time   `json:"CreatedAt"`
	UpdatedAt   time.Time   `json:"UpdatedAt"`
	DeletedAt   pq.NullTime `json:"DeletedAt"`
	Name        string      `json:"Name,omitempty"`
	Description string      `json:"Description,omitempty"`
	Heroes      []Hero      `json:"heroes" json:"Heroes,omitempty"`
	Spells      []Spell     `json:"spells" json:"Spells,omitempty"`
	UserID      uint        `json:"UserID"`
}

type Hero struct {
	ID          uint        `json:"ID,omitempty"`
	CreatedAt   time.Time   `json:"CreatedAt"`
	UpdatedAt   time.Time   `json:"UpdatedAt"`
	DeletedAt   pq.NullTime `json:"DeletedAt"`
	Name        string      `json:"Name,omitempty"`
	Description string      `json:"Description,omitempty"`
	Rarity      string      `json:"Rarity,omitempty"`
	DamageType  string      `json:"DamageType,omitempty"`
	Effect      string      `json:"Effect,omitempty"`
	Hitpoint    int32       `json:"Hitpoint,omitempty"`
	Damage      int32       `json:"Damage,omitempty"`
	CostElixir  int32       `json:"CostElixir,omitempty"`
	DamageTower int32       `json:"DamageTower,omitempty"`
	Speed       int32       `json:"Speed,omitempty"`
	Price       int64       `json:"Price,omitempty"`
}

type Spell struct {
	ID          uint        `json:"ID"`
	CreatedAt   time.Time   `json:"CreatedAt"`
	UpdatedAt   time.Time   `json:"UpdatedAt"`
	DeletedAt   pq.NullTime `json:"DeletedAt"`
	Name        string      `json:"Name"`
	Description string      `json:"Description"`
	Area        int32       `json:"Area"`
	DamageType  string      `json:"DamageType"`
	Damage      int32       `json:"Damage"`
	Duration    int64       `json:"Duration"`
	Effect      string      `json:"Effect"`
	Price       int64       `json:"Price"`
}
