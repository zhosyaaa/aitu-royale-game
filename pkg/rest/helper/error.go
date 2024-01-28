package helper

import "errors"

var (
	ErrHeroNotFound  = errors.New("hero not found")
	ErrSpellNotFound = errors.New("spell not found")
	ErrDeckNotFound  = errors.New("deck not found")
)
