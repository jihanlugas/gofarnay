package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"gofarnay/controller"
	"log"
	"net/http"
)


func main()  {
	log.Println("Listening server at http://localhost:8010")
	router := mux.NewRouter()

	router.HandleFunc("/users", controller.GetUsers).Methods("GET")


	log.Fatal(http.ListenAndServe(":8010", router))

}
