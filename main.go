package main

import (
	"fmt"
	"net/http"
	"tracker/config"
	"tracker/database"
	"tracker/handler"
	"tracker/repo"
	"tracker/routes"
	"tracker/service"
)


func main(){

	// load up variables
	config.LoadEnv()

	// connect to database
	database.ConnectDB()

	//initialise repos
	userRepo := &repo.UserRepo{}
	// transactionRepo := &repo.TransactionRepo{}
	// budgetRepo := &repo.BudgetRepo{}

	//initialise services
	userService := &service.UserService{Repo: userRepo}
	// transactionService := &service.TransactionService{Repo: transactionRepo}
	// budgetService := &service.BudgetService{Repo: budgetRepo}

	//initialise handlers

	userHandler := &handler.UserHandler{Service: userService}
	// transactionHandler := &handler.TransactionHandler{Service: transactionService}
	// budgetHandler := &handler.BudgetHandler{Service: budgetService}

	router := routes.SetupRouter(userHandler)

	fmt.Println("server is running on localhost:8080")
	http.ListenAndServe(":8080", router)
}