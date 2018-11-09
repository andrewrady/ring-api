package main

import (
	"encoding/json"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type User struct {
	gorm.Model
	Email    string
	Password string
}

func NewUser(w http.ResponseWriter, r *http.Request) {
	db, err = gorm.Open("postgres", "host=localhost port=5432 user=postgres dbname=ring_tracker password=postgres sslmode=disable")
	if err != nil {
		panic("Could not connect to the database")
	}
	defer db.Close()

	decoder := json.NewDecoder(r.Body)

	var newUser User
	err := decoder.Decode(&newUser)

	if err != nil {
		panic(err)
	}

	hashedPassword := hashAndSalt([]byte(newUser.Password))

	db.Create(&User{Email: newUser.Email, Password: hashedPassword})
	w.WriteHeader(http.StatusCreated)
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	db, err = gorm.Open("postgres", "host=localhost port=5432 user=postgres dbname=ring_tracker password=postgres sslmode=disable")
	if err != nil {
		panic("Could not connect to the database")
	}
	defer db.Close()

	var users []User
	db.Order("ID asc").Find(&users)
	json.NewEncoder(w).Encode(users)
}

func UserLogin(w http.ResponseWriter, r *http.Request) {
	db, err = gorm.Open("postgres", "host=localhost port=5432 user=postgres dbname=ring_tracker password=postgres sslmode=disable")
	if err != nil {
		panic("Could not connect to the database")
	}
	defer db.Close()

	decoder := json.NewDecoder(r.Body)

	var requestUser User
	err := decoder.Decode(&requestUser)
	if err != nil {
		panic(err)
	}

	var user User
	db.Where("EMAIL = ?", requestUser.Email).Find(&user)
	if user.ID > 0 {
		if CompareHashedPassword(user.Password, []byte(requestUser.Password)) {
			json.NewEncoder(w).Encode(user)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	} else {
		w.WriteHeader(http.StatusNotFound)
	}

}

func hashAndSalt(pwd []byte) string {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		panic(err)
	}
	return string(hash)
}

func CompareHashedPassword(hashedPassword string, plainPassword []byte) bool {
	byteHash := []byte(hashedPassword)

	err := bcrypt.CompareHashAndPassword(byteHash, plainPassword)
	if err != nil {
		return false
	} else {
		return true
	}
}
