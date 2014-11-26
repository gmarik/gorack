package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
)

type RackRequest struct {
	REQUEST_METHOD string
	SCRIPT_NAME    string
	PATH_INFO      string
	QUERY_STRING   string
	SERVER_NAME    string
	SERVER_PORT    string
	HTTP_vars      []string
}

type Handler struct {
	clientReader *os.File
	clientWriter *os.File
	serverWriter *os.File
	serverReader *os.File
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	clientReader, serverWriter, err := os.Pipe()

	if err != nil {
		log.Fatal(err)
	}

	serverReader, clientWriter, err := os.Pipe()

	if err != nil {
		log.Fatal(err)
	}

	rr := RackRequest{
		REQUEST_METHOD: r.Method,
		SCRIPT_NAME:    r.URL.Path,
		PATH_INFO:      r.URL.Path,
		QUERY_STRING:   r.URL.RawQuery,
		SERVER_NAME:    "hello",
		SERVER_PORT:    "80",
	}

	jsonData, err := json.Marshal(rr)

	if err != nil {
		w.WriteHeader(500)
		return
	}

	serverWriter.Write(jsonData)
	serverWriter.Close()

	cmd := exec.Command("./gorack", "./config.ru", strconv.Itoa(int(clientReader.Fd())), strconv.Itoa(int(clientWriter.Fd())))

	out, err := cmd.StdoutPipe()

	cmd.ExtraFiles = []*os.File{clientReader, clientWriter}

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

	io.Copy(w, serverReader)
}

func main() {

	// writer, err := os.OpenFile("/tmp/123.sock", os.O_RDWR|os.O_CREATE, 0777)
	//

	handler := &Handler{}
	http.ListenAndServe("localhost:3001", handler)
}
