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

	router.HandleFunc("/signin", controller.Signin).Methods("POST")
	router.HandleFunc("/auth", controller.Auth).Methods("GET")

	router.HandleFunc("/users", controller.GetUsers).Methods("GET")

	router.HandleFunc("/order-form", controller.OrderForm).Methods("POST")
	router.HandleFunc("/orders", controller.GetOrders).Methods("GET")

	router.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))
	router.PathPrefix("/uploads/").Handler(http.StripPrefix("/uploads/", http.FileServer(http.Dir("uploads"))))

	log.Fatal(http.ListenAndServe(":8010", router))

}
