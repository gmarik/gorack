package main

import (
	"log"
	"net/http"
	"os/exec"
)

type Handler struct{}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	cmd := exec.Command("./gorack", "./config.ru", "/tmp/123.sock")
	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Waiting for command to finish...")
	err = cmd.Wait()
	log.Printf("Command finished with error: %v", err)

}

func main() {

	handler := &Handler{}
	http.ListenAndServe("localhost:3001", handler)
}
