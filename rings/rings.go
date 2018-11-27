package rings

import (
	"encoding/json"
	"fmt"
	"net/http"

	Config "rings-api/config"

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
	Status     bool
}

//InitialMigration initial connection to database
func InitialMigration() {
	db, err := gorm.Open("postgres", Config.DbConnectionString)
	if err != nil {
		fmt.Println(err.Error())
		panic("Failed to connect to database")
	}
	defer db.Close()

	db.AutoMigrate(&Ring{})
}

//AllRings return json for every ring
func AllRings(w http.ResponseWriter, r *http.Request) {
	db, err = gorm.Open("postgres", Config.DbConnectionString)
	if err != nil {
		panic("Could not connect to the database")
	}
	defer db.Close()

	var rings []Ring
	db.Order("ID asc").Find(&rings)
	json.NewEncoder(w).Encode(rings)
}

//NewRing creates new record
func NewRing(w http.ResponseWriter, r *http.Request) {
	db, err = gorm.Open("postgres", Config.DbConnectionString)
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

	db.Create(&Ring{RingNumber: newRing.RingNumber, Age: newRing.Age, Rank: newRing.Rank, Division: newRing.Division, Gender: newRing.Gender, Status: newRing.Status})

	var rings []Ring
	db.Last(&rings)
	json.NewEncoder(w).Encode(rings)
}

//ShowRing shows one ring based on id passed in params
func ShowRing(w http.ResponseWriter, r *http.Request) {
	db, err = gorm.Open("postgres", Config.DbConnectionString)
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

//DeleteRing deletes ring
func DeleteRing(w http.ResponseWriter, r *http.Request) {
	db, err = gorm.Open("postgres", Config.DbConnectionString)
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

//UpdateRing updates the ring record
func UpdateRing(w http.ResponseWriter, r *http.Request) {
	db, err = gorm.Open("postgres", Config.DbConnectionString)
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
	if ring.ID > 0 {
		ring.RingNumber = updatedRing.RingNumber
		ring.Age = updatedRing.Age
		ring.Rank = updatedRing.Rank
		ring.Division = updatedRing.Division
		ring.Gender = updatedRing.Division
		ring.Status = updatedRing.Status

		db.Save(&ring)

		json.NewEncoder(w).Encode(&ring)
	} else {
		w.WriteHeader(http.StatusUnprocessableEntity)
	}

}
