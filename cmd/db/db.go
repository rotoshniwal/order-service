package main

import (
	_ "github.com/lib/pq"
	"order/app"
	"log"
)

//main is invoked when this executable runs.
//Establishes a connection to desired DB and initializes the DB
func main() {

	// initialize the DB and create required tables
	app.InitDB()
	log.Println("Initialized the DB successfully")
}