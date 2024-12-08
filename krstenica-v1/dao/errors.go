package dao

import "errors"

var (
	ErrHramNotFound = errors.New("Hram nije nadjen")
	ErrHramDeleted  = errors.New("Hram je obrisan")
)
