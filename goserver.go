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
	writer io.Writer
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	go func() {
		h.writer.Write([]byte("Hello"))

		// fmt.Println(ioutil.ReadAll(reader))
	}()

	cmd := exec.Command("./gorack", "./config.ru", "/tmp/123.sock")

	out, err := cmd.StdoutPipe()

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

	_, writer := io.Pipe()

	go func() {
		writer.Write([]byte("hello\n"))
		writer.Close()
	}()

	handler := &Handler{writer}
	http.ListenAndServe("localhost:3001", handler)
}
