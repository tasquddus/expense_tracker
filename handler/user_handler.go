package handler

import (
	"encoding/json"
	"net/http"
	"tracker/models"
	"tracker/service"
)

type UserHandler struct {
	Service *service.UserService
}

func (h *UserHandler) RegisterUser (w http.ResponseWriter, r *http.Request){
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	err = h.Service.RegisterUser(&user)
	if err != nil {
		http.Error(w, "could not register user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)

}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	 var loginRequest models.LoginUser
	 err := json.NewDecoder(r.Body).Decode(&loginRequest)
	 if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	 }
	 
	 token, err := h.Service.LoginUser(loginRequest)
	 if err != nil {
		http.Error(w, "invalid credentials", http.StatusInternalServerError)
		return
	 }

	 w.WriteHeader(http.StatusOK)
	 json.NewEncoder(w).Encode(token)

}