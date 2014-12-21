package gorack

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"syscall"

	"../ipcio"
)

var gorackRunner = "./rack/gorack"

type RackHandler struct {
	local_fd   int
	configPath string
}

func NewRackHandler(configPath string) *RackHandler {
	pair, err := syscall.Socketpair(syscall.AF_LOCAL, syscall.SOCK_STREAM, 0)

	if err != nil {
		log.Fatal(err)
	}

	local, remote := pair[0], pair[1]

	// child process' FDs start from 3 (0, 1, 2)
	fd := os.NewFile(uintptr(remote), "master_io")

	cmd := exec.Command(gorackRunner, configPath)
	cmd.ExtraFiles = []*os.File{fd}

	go runProcessMaster(cmd)

	return &RackHandler{
		local_fd:   local,
		configPath: configPath,
	}
}

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

func (s *RackHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	jsonData, err := json.Marshal(NewRackRequest(r, "server", "port"))

	if err != nil {
		log.Fatal(err)
	}

	res_reader, req_writer, err := SendIo(s.local_fd)

	if err != nil {
		log.Fatal(err)
	}

	req_writer.Write(jsonData)
	req_writer.Close()

	// resp := NewResponse(io.TeeReader(res_reader, os.Stdout))
	resp := NewResponse(res_reader)

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

func runProcessMaster(cmd *exec.Cmd) {
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
