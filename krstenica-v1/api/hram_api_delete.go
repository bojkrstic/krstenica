package api

import (
	"krstenica/krstenica-v1/dao"
	"log"
	"net/http"

	"krstenica/pkg/apiutil"
)

type HramDelete struct {
	apiutil.PathRegistry
	ID uint `path:"id"`
}

func (h *HramDelete) Handle(w http.ResponseWriter, r *http.Request) (interface{}, error) {
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
	//promena statusa
	del, err := db.DeleteHram(uint(h.ID))
	if err != nil {
		log.Println(err)
		if err == dao.ErrHramDeleted {
			return nil, ErrHramAlreadyDeleted
		}
	}

	res := makeResultSysApplication(del)

	// dlogger.Log(lrec)
	return res, nil

}
