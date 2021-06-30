package controller

import (
	"database/sql"
	"encoding/json"
	"gofarnay/model"
	"log"
	"net/http"
)

func Signin(w http.ResponseWriter, r *http.Request) {
	var c model.Credentials
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&c); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	if err := c.Signin(); err != nil {
		switch err {
		case sql.ErrNoRows:
			RespondWithError(w, http.StatusNotFound, "Invalid Email or Password")
		default:
			RespondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	token, err := c.GenerateToken()
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	RespondWithJSONAddToken(w, http.StatusOK, map[string]string{"result": "success"}, token)
}

func Auth(w http.ResponseWriter, r *http.Request) {
	log.Println("Auth")
	id := 1

	u := model.User{ID: id}

	//RespondWithError(w, http.StatusNotFound, "User not found")

	if err := u.GetUser(); err != nil {
		switch err {
		case sql.ErrNoRows:
			RespondWithError(w, http.StatusNotFound, "User not found")
		default:
			RespondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	RespondWithJSON(w, http.StatusOK, u)
}

