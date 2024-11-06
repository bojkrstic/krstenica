package api

import (
	"krstenica/krstenica-v1/dao"
	"krstenica/pkg/apiutil"
	"log"
	"net/http"
	"time"
)

// HramAdd handlers POST HTTP request
type HramAdd struct {
}

func (ac HramAdd) Handle(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	//input
	var reqData HramWo
	err := apiutil.GetRequestBody(r, &reqData)
	if err != nil {
		log.Println(err)
		return nil, apiutil.NewIntError(err)
	}

	// if !apiutil.IsEmail(reqData.NazivHrama) {
	// 	return nil, ErrInvalidNazivHrama
	// }
	//check if hram exist with this email
	hram, err := db.GetHramByName(reqData.NazivHrama)
	if err != nil {
		if err != dao.ErrHramNotFound {
			log.Println(err)
			return nil, err
		}
	}

	if hram != nil {
		return nil, ErrHramExistWithThisName
	}
	//create new hram with this email
	newHram := &dao.HramDo{
		HramName:  reqData.NazivHrama,
		CreatedAt: time.Now(),
	}
	//create new hram
	newUserID, err := db.CreateHram(newHram)
	if err != nil {
		log.Println(err)
		return nil, apiutil.NewIntError(err)
	}

	//get hram from db, to create return json
	mainUser, err := db.GetHram(newUserID)
	if err != nil {
		log.Println(err)
		if err == dao.ErrHramNotFound {
			return nil, ErrHramNotFound
		}
	}
	resWo := &HramCrtResWo{
		HramID:     mainUser.HramID,
		NazivHrama: mainUser.HramName,
		CreatedAt:  mainUser.CreatedAt.Format("2006-01-02 15:15:05"),
	}

	return resWo, nil
}
