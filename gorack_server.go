package main

import (
	"log"
	"net/http"
	"os"

	"./gorack"
)

func main() {

	if len(os.Args) == 1 {
		log.Fatal("specify path to config.ru file")
	}

	config_path := os.Args[1]

	address := "localhost:3001"
	log.Print("Starting on:", address)
	log.Fatal(http.ListenAndServe(address, gorack.NewRackHandler(config_path)))
}
