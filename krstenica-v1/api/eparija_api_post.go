package api

// import (
// 	"krstenica/krstenica-v1/dao"
// 	"krstenica/pkg/apiutil"
// 	"log"
// 	"net/http"
// 	"time"
// )

// // HramAdd handlers POST HTTP request
// type EparhijaAdd struct {
// }

// func (ac EparhijaAdd) Handle(w http.ResponseWriter, r *http.Request) (interface{}, error) {
// 	//input
// 	var reqData EparhijaWo
// 	err := apiutil.GetRequestBody(r, &reqData)
// 	if err != nil {
// 		log.Println(err)
// 		return nil, apiutil.NewIntError(err)
// 	}

// 	if !apiutil.IsEmail(reqData.NazivEparhije) {
// 		return nil, ErrInvalidEmailAdr
// 	}
// 	//check if user exist with this email
// 	user, err := db.GetUserByEmail(reqData.Email)
// 	if err != nil {
// 		if err != dao.ErrUserNotFound {
// 			log.Println(err)
// 			return nil, err
// 		}
// 	}

// 	if user != nil {
// 		return nil, ErrUserExistWithThisEmail
// 	}
// 	//create new user with this email
// 	newUser := &dao.UserDo{
// 		Email:     reqData.Email,
// 		CreatedAt: time.Now(),
// 	}
// 	//create new user
// 	newUserID, err := db.CreateUser(newUser)
// 	if err != nil {
// 		log.Println(err)
// 		return nil, apiutil.NewIntError(err)
// 	}

// 	//get user from db, to create return json
// 	mainUser, err := db.GetUser(newUserID)
// 	if err != nil {
// 		log.Println(err)
// 		if err == dao.ErrUserNotFound {
// 			return nil, ErrUserNotFound
// 		}
// 	}
// 	resWo := &UserCrtResWo{
// 		UserID:    mainUser.UserID,
// 		Email:     mainUser.Email,
// 		CreatedAt: mainUser.CreatedAt.Format("2006-01-02 15:15:05"),
// 	}

// 	return resWo, nil
// }
