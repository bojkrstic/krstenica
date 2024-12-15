package dao

import "errors"

var (
	ErrHramNotFound   = errors.New("Hram nije nadjen")
	ErrHramDeleted    = errors.New("Hram je obrisan")
	ErrHramDubleValue = errors.New("Hram sa istim imenom vec postoji")
	ErrBadPageNumber  = errors.New("Pogresan broj stranice")
)
