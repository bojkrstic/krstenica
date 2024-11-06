package api

import (
	"krstenica/krstenica-v1/dao"
	"log"
	"net/http"
)

type HramGet struct {
	ID int `path:"id"`
}

func (h *HramGet) Handle(w http.ResponseWriter, r *http.Request) (interface{}, error) {

	//get hram from db, to create return json
	setup, err := db.GetHram(uint(h.ID))
	if err != nil {
		log.Println(err)
		if err == dao.ErrHramNotFound {
			return nil, ErrHramNotFound
		}
	}
	res := makeResultSysApplication(setup)

	// dlogger.Log(lrec)
	return res, nil

}
