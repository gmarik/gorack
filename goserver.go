package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
)

type Handler struct {
	writer *os.File
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	cmd := exec.Command("./gorack", "./config.ru")

	out, err := cmd.StdoutPipe()

	cmd.ExtraFiles = []*os.File{h.writer}

	if err != nil {
		fmt.Println(err)
	}

	err = cmd.Start()
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		io.Copy(os.Stdout, out)
	}()

	log.Printf("Waiting for command to finish...")
	err = cmd.Wait()
	log.Printf("Command finished with error: %v", err)

}

func main() {

	// writer, err := os.OpenFile("/tmp/123.sock", os.O_RDWR|os.O_CREATE, 0777)
	//
	reader, writer, err := os.Pipe()

	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			writer.Write([]byte("it works\n"))
		}
	}()

	handler := &Handler{reader}
	http.ListenAndServe("localhost:3001", handler)
}
