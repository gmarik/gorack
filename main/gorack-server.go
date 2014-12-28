package main

import (
	"flag"
	"log"
	"net/http"

	"gmarik/gorack"
)

func main() {

	var (
		config_path    *string = flag.String("config", "./config.ru", "rack config file")
		listen_address *string = flag.String("address", "localhost:3000", "address to listen at")
	)

	flag.Parse()

	log.Print("Listening at ", *listen_address)
	log.Fatal(http.ListenAndServe(*listen_address, gorack.NewRackHandler(*config_path)))
}
