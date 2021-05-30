package controller

import (
	"database/sql"
	"encoding/json"
	"github.com/gorilla/mux"
	"gofarnay/model"
	"net/http"
	"strconv"
)

func GetUsers(w http.ResponseWriter, r *http.Request) {
	count, _ := strconv.Atoi(r.FormValue("count"))
	start, _ := strconv.Atoi(r.FormValue("start"))

	if count > 10 || count < 1 {
		count = 10
	}
	if start < 0 {
		start = 0
	}

	users, err := model.GetUsers(start, count)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	RespondWithJSON(w, http.StatusOK, users)
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var u model.User
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&u); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	if err := u.CreateUser(); err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	RespondWithJSON(w, http.StatusCreated, u)
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid User ID")
		return
	}

	u := model.User{ID: id}

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

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid User ID")
		return
	}

	var u model.User
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&u); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	defer r.Body.Close()
	u.ID = id

	if err := u.UpdateUser(); err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	RespondWithJSON(w, http.StatusOK, u)
}

func DeleteUser(w http.ResponseWriter, r *http.Request)  {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid User ID")
		return
	}

	u := model.User{ID: id}
	if err := u.DeleteUser(); err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	RespondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}