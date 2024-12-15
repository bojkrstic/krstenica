package api

import (
	"krstenica/krstenica-v1/dao"
	"krstenica/pkg/apiutil"
	"log"
	"net/http"
)

type HramUpdate struct {
	ID uint `path:"id"`
}

// Handle Update
func (ac *HramUpdate) Handle(w http.ResponseWriter, r *http.Request) (interface{}, error) {

	var updateData HramUpdateData
	err := apiutil.GetRequestBody(r, &updateData)
	if err != nil {
		log.Println(err)
		return nil, apiutil.NewIntError(err)
	}

	//check if hram exist with this email
	_, err = db.GetHramByName(*updateData.NazivHrama)
	if err != nil {
		if err != dao.ErrHramNotFound {
			log.Println(err)
			return nil, err
		}
	}

	//validate

	update, err := updateData.Validate(ac.ID)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	//update
	err = db.UpdateHram(ac.ID, update)
	if err != nil {
		log.Println(err)
		return nil, apiutil.NewIntError(err)
	}

	// get hram from db, to create return json
	mainUser, err := db.GetHram(ac.ID)
	if err != nil {
		log.Println(err)
		if err == dao.ErrHramNotFound {
			return nil, ErrHramNotFound
		}
	}
	resWo := &HramCrtResWo{
		HramID:     mainUser.HramID,
		NazivHrama: mainUser.HramName,
		Status:     mainUser.Status,
		CreatedAt:  mainUser.CreatedAt.Format("2006-01-02 15:15:05"),
	}

	return resWo, nil
}
