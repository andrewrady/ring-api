package users

import (
	"encoding/json"
	"net/http"
	"time"

	Config "ring-api/config"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"golang.org/x/crypto/bcrypt"
)

var db *gorm.DB
var err error

type User struct {
	gorm.Model
	Email    string
	Password string
}

//NewUser creates new user record
func NewUser(w http.ResponseWriter, r *http.Request) {
	db, err = gorm.Open("postgres", Config.DbConnectionString)
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

//GetUsers returns all users
func GetUsers(w http.ResponseWriter, r *http.Request) {
	db, err = gorm.Open("postgres", Config.DbConnectionString)
	if err != nil {
		panic("Could not connect to the database")
	}
	defer db.Close()

	var users []User
	db.Order("ID asc").Find(&users)
	json.NewEncoder(w).Encode(users)
}

//UserLogin compares user data and compares to credential and returns a jwt if everything is correct
func UserLogin(w http.ResponseWriter, r *http.Request) {
	db, err = gorm.Open("postgres", Config.DbConnectionString)
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
		if compareHashedPassword(user.Password, []byte(requestUser.Password)) {
			jwt, err := generateJWT(user.Email)
			if err != nil {
				panic("Error creating JWT")
			}
			json.NewEncoder(w).Encode(jwt)
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

func compareHashedPassword(hashedPassword string, plainPassword []byte) bool {
	byteHash := []byte(hashedPassword)

	err := bcrypt.CompareHashAndPassword(byteHash, plainPassword)
	if err != nil {
		return false
	} else {
		return true
	}
}

func generateJWT(user string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["user"] = user
	claims["exp"] = time.Now().Add(time.Minute * 60).Unix()

	tokenString, err := token.SignedString(Config.MySigningKey)

	if err != nil {
		panic(err)
	}

	return tokenString, nil
}
