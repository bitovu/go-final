package main

import (
	"final/pkg/db"
	"final/pkg/server"
	"log"
)

func main() {
	if err := db.Init("scheduler.db"); err != nil {
		log.Fatal(err)
	}
	server.Run()
}
