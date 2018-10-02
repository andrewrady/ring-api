package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func helloWorld(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World")
}

func handleRequest() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", helloWorld).Methods("GET")
	myRouter.HandleFunc("/rings", AllRings).Methods("GET")
	myRouter.HandleFunc("/rings", NewRing).Methods("POST")
	myRouter.HandleFunc("/ring/{id}", DeleteRing).Methods("DELETE")
	myRouter.HandleFunc("/ring/{id}", UpdateRing).Methods("PUT")
	log.Fatal(http.ListenAndServe(":3000", myRouter))
}

func main() {
	fmt.Println("Server is running")
	InitialMigration()
	handleRequest()
}
