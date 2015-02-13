package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/gmarik/gorack"
)

func main() {

	gorack.GoRackExec = path.Join(path.Dir(os.Args[0]), "gorack-ruby")

	var (
		config_path    *string = flag.String("config", "./config.ru", "rack config file")
		listen_address *string = flag.String("address", "localhost:3000", "address to listen at")
	)

	flag.Parse()

	log.Print("Listening at ", *listen_address)
	handler := gorack.NewRackHandler(*config_path)
	if err := handler.StartRackProcess(); err != nil {
		log.Fatal(err)
	}
	log.Fatal(http.ListenAndServe(*listen_address, handler))
}
