package config

import (
	"os"
)

//MySigningKey secret signature Key
var MySigningKey = []byte(os.Getenv("mySigningKey"))

//DbConnectionString database connection
var DbConnectionString = "host=" + os.Getenv("dbHost") + " port=" + os.Getenv("dbPort") + " user=" + os.Getenv("dbUser") + " dbname=" + os.Getenv("dbName") + " password=" + os.Getenv("dbPassword") + " sslmode=disable"
