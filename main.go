package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	jwt "github.com/dgrijalva/jwt-go"

	Config "ring-api/config"
	Rings "ring-api/rings"
	Users "ring-api/users"

	"github.com/gorilla/mux"
)

func handleRequest() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/rings", Rings.AllRings).Methods("GET")
	myRouter.Handle("/rings", isAuthorized(Rings.NewRing)).Methods("POST")
	myRouter.Handle("/rings/{id}", isAuthorized(Rings.ShowRing)).Methods("GET")
	myRouter.Handle("/rings/{id}", isAuthorized(Rings.DeleteRing)).Methods("DELETE")
	myRouter.Handle("/rings/{id}", isAuthorized(Rings.UpdateRing)).Methods("PUT")
	//User Routes
	myRouter.Handle("/users", isAuthorized(Users.GetUsers)).Methods("GET")
	myRouter.Handle("/users", isAuthorized(Users.NewUser)).Methods("POST")
	myRouter.HandleFunc("/users/login", Users.UserLogin).Methods("POST")
	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), myRouter))
}

func isAuthorized(endpoint func(http.ResponseWriter, *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.Header["Authorization"] != nil {
			token, err := jwt.Parse(r.Header["Authorization"][0], func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("There was an error")
				}
				return Config.MySigningKey, nil
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
	fmt.Println("Server is running on port " + os.Getenv("PORT"))
	Rings.InitialMigration()
	handleRequest()
}
