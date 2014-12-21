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

func SendIo(fd int) (*os.File, *os.File, error) {

	req_reader, req_writer, err := os.Pipe()

	if err != nil {
		return nil, nil, err
	}

	res_reader, res_writer, err := os.Pipe()

	if err != nil {
		return nil, nil, err
	}

	err = ipcio.SendIo(fd, req_reader)

	if err != nil {
		return nil, nil, err
	}

	err = ipcio.SendIo(fd, res_writer)

	if err != nil {
		return nil, nil, err
	}

	req_reader.Close()
	res_writer.Close()

	return res_reader, req_writer, nil
}

func ServeHttp(local_fd int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		jsonData, err := json.Marshal(gorack.NewRackRequest(r, "server", "port"))

		if err != nil {
			log.Fatal(err)
		}

		res_reader, req_writer, err := SendIo(local_fd)

		if err != nil {
			log.Fatal(err)
		}

		req_writer.Write(jsonData)
		req_writer.Close()

		// resp := gorack.NewResponse(io.TeeReader(res_reader, os.Stdout))
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
			log.Fatal(err)
		}
	}
}

type RackHandler struct {
	local_fd int
}

func (s *RackHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ServeHttp(s.local_fd)(w, r)
}

func main() {

	pair, err := syscall.Socketpair(syscall.AF_LOCAL, syscall.SOCK_STREAM, 0)

	if err != nil {
		log.Fatal(err)
	}

	remote, local := pair[0], pair[1]

	if len(os.Args) == 1 {
		log.Fatal("specify path to config.ru file")
	}

	config_path := os.Args[1]

	// child process' FDs start from 3 (0, 1, 2)
	master_io := os.NewFile(uintptr(remote), "master_io")
	go runProcessMaster(master_io, "./ruby/gorack", config_path)

	address := "localhost:3001"
	log.Print("Starting on:", address)
	log.Fatal(http.ListenAndServe(address, &RackHandler{local}))
}

func runProcessMaster(fd *os.File, bin_path string, args ...string) {
	cmd := exec.Command(bin_path, args...)

	cmd.ExtraFiles = []*os.File{fd}

	var err error
	var out, outerr io.Reader

	if out, err = cmd.StdoutPipe(); err != nil {
		log.Fatal(err)
	}

	if outerr, err = cmd.StderrPipe(); err != nil {
		log.Fatal(err)
	}

	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	go io.Copy(NewLogWriter(os.Stdout, "", log.LstdFlags), out)
	go io.Copy(NewLogWriter(os.Stderr, "[StdErr]", log.LstdFlags), outerr)

	if err = cmd.Wait(); err != nil {
		log.Fatal(err)
	}
}

type LogWriter struct {
	logger *log.Logger
	prefix string
}

func NewLogWriter(w io.Writer, prefix string, flags int) *LogWriter {
	return &LogWriter{log.New(w, "", flags), prefix}
}

func (l *LogWriter) Write(data []byte) (int, error) {
	l.logger.Printf("%s%s", l.prefix, data)
	return len(data), nil
}
