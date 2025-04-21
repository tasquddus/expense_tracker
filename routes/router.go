package routes

import (
	"tracker/handler"
	"tracker/middleware"

	"github.com/gorilla/mux"
)

// public routes

func SetupRouter(userHandler *handler.UserHandler) *mux.Router{
	// budgetHandler *handler.BudgetHandler, transacionHandler *handler.TransactionHandler){
	r := mux.NewRouter()

	// public routes
	r.HandleFunc("/register", userHandler.RegisterUser).Methods("POST")
	r.HandleFunc("/login", userHandler.Login).Methods("POST")

	protected := r.PathPrefix("/").Subrouter()
	protected.Use(middleware.AuthMiddleware)

	// Editor only roites


	// Authenticated user routes


	return r
}