package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"go-contacts/app"
	"go-contacts/controllers"
	"net/http"
	"os"
	httpSwagger "github.com/swaggo/http-swagger"
)

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/api/register", controllers.SignUp).Methods("POST")
	router.HandleFunc("/api/login", controllers.SignIn).Methods("POST")

	router.PathPrefix("/swagger").Handler(httpSwagger.WrapHandler)

	router.Use(app.JwtAuthentication)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	err := http.ListenAndServe(":" + port, router) //localhost:8000/api
	if err != nil {
		fmt.Print(err)
	}
}