package main

import (
	"courier-app/config"
	"courier-app/db"
	"courier-app/server"
	"fmt"
)

func main() {
	config.InitENV()
	dbErr := db.Connect()
	if dbErr != nil {
		fmt.Println("error connecting to DB")
	}
	server.Init()

}
