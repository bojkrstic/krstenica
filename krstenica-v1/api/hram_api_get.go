package api

import (
	"krstenica/krstenica-v1/dao"
	"log"
	"net/http"

	"krstenica/pkg/apiutil"
)

type HramGet struct {
	apiutil.PathRegistry
	ID uint `path:"id"`
}

func (h *HramGet) Handle(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	log.Println("Handler called with ID:", h.ID)
	//get hram from db, to create return json
	setup, err := db.GetHram(uint(h.ID))
	if err != nil {
		log.Println(err)
		if err == dao.ErrHramNotFound {
			return nil, ErrHramNotFound
		}
	}
	if setup.Status == "deleted" {
		return nil, ErrHramAlreadyDeleted
	}

	res := makeResultSysApplication(setup)

	// dlogger.Log(lrec)
	return res, nil

}
