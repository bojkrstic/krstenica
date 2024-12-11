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
func (h *HramUpdate) Handle(w http.ResponseWriter, r *http.Request) (interface{}, error) {

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
	return nil, nil
}
