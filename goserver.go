package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"syscall"

	"./gorack"
	"./ipcio"
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

func ServeHttp(local_fd int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

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
			log.Fatal(err)
		}

		req_reader, req_writer, err := os.Pipe()

		if err != nil {
			log.Fatal(err)
		}

		res_reader, res_writer, err := os.Pipe()

		if err != nil {
			log.Fatal(err)
		}

		err = ipcio.SendIo(local_fd, req_reader)
		if err != nil {
			log.Fatal(err)
		}

		err = ipcio.SendIo(local_fd, res_writer)

		if err != nil {
			log.Fatal(err)
		}

		req_writer.Write(jsonData)
		req_writer.Close()
		req_reader.Close()
		res_writer.Close()

		// resp := gorack.NewResponse(io.TeeReader(serverReader, os.Stdout))
		resp := gorack.NewResponse(res_reader)

		if err := resp.Parse(); err != nil {
			log.Println("Error:", err.Error())
			http.Error(w, err.Error(), http.StatusServiceUnavailable)
			return
		}

		for name, values := range resp.Headers {
			for _, val := range values {
				// fmt.Println(name, val)
				w.Header().Add(name, val)
			}
		}

		w.WriteHeader(resp.StatusCode)

		_, err = io.Copy(w, resp.Body)

		if err != nil {
			log.Println(err.Error())
		}
	}
}

func main() {

	pair, err := syscall.Socketpair(syscall.AF_LOCAL, syscall.SOCK_STREAM, 0)

	if err != nil {
		log.Fatal(err)
	}

	remote, local := pair[0], pair[1]

	go runProcessMaster(remote)

	http.HandleFunc("/", ServeHttp(local))
	http.ListenAndServe("localhost:3001", nil)
}

func runProcessMaster(remote_fd int) {
	cmd := exec.Command("./gorack.sh", "./config.ru")

	out, err := cmd.StdoutPipe()
	erro, err := cmd.StderrPipe()

	// child process' FDs start from 3 (0, 1, 2)
	master_io := os.NewFile(uintptr(remote_fd), "master_io")
	cmd.ExtraFiles = []*os.File{master_io}

	if err = cmd.Start(); err != nil {
		log.Fatal(err)
	}

	go io.Copy(os.Stdout, out)
	go io.Copy(os.Stderr, erro)

	err = cmd.Wait()
}
