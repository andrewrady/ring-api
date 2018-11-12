package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	jwt "github.com/dgrijalva/jwt-go"

	"github.com/gorilla/mux"
)

var mySigningKey = []byte(os.Getenv("mySigningKey"))
var dbConnectionString = "host=" + os.Getenv("dbHost") + " port=" + os.Getenv("dbPort") + " user=" + os.Getenv("dbUser") + " dbname=" + os.Getenv("dbName") + " password=" + os.Getenv("dbPassword") + " sslmode=disable"

func handleRequest() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/rings", AllRings).Methods("GET")
	myRouter.Handle("/rings", isAuthorized(NewRing)).Methods("POST")
	myRouter.Handle("/rings/{id}", isAuthorized(ShowRing)).Methods("GET")
	myRouter.Handle("/rings/{id}", isAuthorized(DeleteRing)).Methods("DELETE")
	myRouter.Handle("/rings/{id}", isAuthorized(UpdateRing)).Methods("PUT")
	//User Routes
	myRouter.Handle("/users", isAuthorized(GetUsers)).Methods("GET")
	myRouter.Handle("/users", isAuthorized(NewUser)).Methods("POST")
	myRouter.HandleFunc("/users/login", UserLogin).Methods("POST")
	log.Fatal(http.ListenAndServe(":3000", myRouter))
}

func isAuthorized(endpoint func(http.ResponseWriter, *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.Header["Authorization"] != nil {
			token, err := jwt.Parse(r.Header["Authorization"][0], func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("There was an error")
				}
				return mySigningKey, nil
			})

			if err != nil {
				panic(err)
			}

			if token.Valid {
				endpoint(w, r)
			}
		} else {
			w.WriteHeader(http.StatusUnauthorized)
		}
	})
}

func main() {
	fmt.Println("Server is running")
	InitialMigration()
	handleRequest()
}
