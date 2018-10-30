package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var db *gorm.DB
var err error

type Ring struct {
	gorm.Model
	RingNumber int
	Age        string
	Rank       string
	Division   string
	Gender     string
}

func InitialMigration() {
	db, err := gorm.Open("postgres", "host=localhost port=5432 user=postgres dbname=ring_tracker password=postgres sslmode=disable")
	if err != nil {
		fmt.Println(err.Error())
		panic("Failed to connect to database")
	}
	defer db.Close()

	db.AutoMigrate(&Ring{})
}

func AllRings(w http.ResponseWriter, r *http.Request) {
	db, err = gorm.Open("postgres", "host=localhost port=5432 user=postgres dbname=ring_tracker password=postgres sslmode=disable")
	if err != nil {
		panic("Could not connect to the database")
	}
	defer db.Close()

	var rings []Ring
	db.Find(&rings)
	json.NewEncoder(w).Encode(rings)
}

func NewRing(w http.ResponseWriter, r *http.Request) {
	db, err = gorm.Open("postgres", "host=localhost port=5432 user=postgres dbname=ring_tracker password=postgres sslmode=disable")
	if err != nil {
		panic("Could not connect to the database")
	}
	defer db.Close()

	decoder := json.NewDecoder(r.Body)

	var newRing Ring
	err := decoder.Decode(&newRing)

	if err != nil {
		panic(err)
	}

	db.Create(&Ring{RingNumber: newRing.RingNumber, Age: newRing.Age, Rank: newRing.Rank, Division: newRing.Division, Gender: newRing.Gender})

	var rings []Ring
	db.Last(&rings)
	json.NewEncoder(w).Encode(rings)
}

func ShowRing(w http.ResponseWriter, r *http.Request) {
	db, err = gorm.Open("postgres", "host=localhost port=5432 user=postgres dbname=ring_tracker password=postgres sslmode=disable")
	if err != nil {
		panic("Could not connect to the database")
	}
	defer db.Close()

	vars := mux.Vars(r)
	ringID := vars["id"]

	var ring Ring
	db.Where("ID = ?", ringID).Find(&ring)
	if ring.ID > 0 {
		json.NewEncoder(w).Encode(ring)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func DeleteRing(w http.ResponseWriter, r *http.Request) {
	db, err = gorm.Open("postgres", "host=localhost port=5432 user=postgres dbname=ring_tracker password=postgres sslmode=disable")
	if err != nil {
		panic("Could not connect to the database")
	}
	defer db.Close()

	vars := mux.Vars(r)
	ringID := vars["id"]

	var ring Ring
	db.Where("ID = ?", ringID).Find(&ring)
	db.Delete(&ring)
}

func UpdateRing(w http.ResponseWriter, r *http.Request) {
	db, err = gorm.Open("postgres", "host=localhost port=5432 user=postgres dbname=ring_tracker password=postgres sslmode=disable")
	if err != nil {
		panic("Could not connect to the database")
	}
	defer db.Close()

	vars := mux.Vars(r)
	ringID := vars["id"]

	var updatedRing Ring
	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&updatedRing)

	if err != nil {
		panic(err)
	}

	var ring Ring
	db.Where("ID = ?", ringID).Find(&ring)

	ring.RingNumber = updatedRing.RingNumber
	ring.Age = updatedRing.Age
	ring.Rank = updatedRing.Rank
	ring.Division = updatedRing.Division
	ring.Gender = updatedRing.Division

	db.Save(&ring)

	json.NewEncoder(w).Encode(&ring)
}
