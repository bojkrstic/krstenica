package api

import (
	"krstenica/pkg/apiutil"
	"net/http"
)

var (
	ErrInvalidNazivHrama         = apiutil.NewError(http.StatusBadRequest, "NAZIV_HRAMA_IS_INVALED", "Pogresano unesen naziv Hrama")
	ErrInvalidNazivEparhije      = apiutil.NewError(http.StatusBadRequest, "NAZIV_EPARHIJE_IS_INVALED", "Pogresano unesen naziv Eparhije")
	ErrHramExistWithThisName     = apiutil.NewError(http.StatusBadRequest, "HRAM_EXIST", "Hram sa ovim imenom postoji")
	ErrEparhijaExistWithThisName = apiutil.NewError(http.StatusBadRequest, "EPARHIJA_EXIST", "Eparhija sa ovim imenom postoji")
	ErrHramNotFound              = apiutil.NewError(http.StatusNotFound, "HRAM_NOT_FOUND", "Hram nije nadjen")
	ErrEparhijaNotFound          = apiutil.NewError(http.StatusNotFound, "EPARHIJA_NOT_FOUND", "Eparhija nije nadjena")
)
