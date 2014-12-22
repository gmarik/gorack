package gorack

import (
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"syscall"

	"../ipcio"
)

var gorackRunner = "./ruby/gorack"

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

func (s *RackHandler) sendIo() (*os.File, *os.File, error) {

	req_reader, req_writer, err := os.Pipe()

	if err != nil {
		return nil, nil, err
	}

	res_reader, res_writer, err := os.Pipe()

	if err != nil {
		return nil, nil, err
	}

	err = ipcio.SendIo(s.local_fd, req_reader)

	if err != nil {
		return nil, nil, err
	}

	err = ipcio.SendIo(s.local_fd, res_writer)

	if err != nil {
		return nil, nil, err
	}

	// Once sent a the process - close
	// they'll still be open in the process
	// to read/write
	req_reader.Close()
	res_writer.Close()

	return res_reader, req_writer, nil
}

func (s *RackHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	res_reader, req_writer, err := s.sendIo()

	if err != nil {
		log.Fatal(err)
	}

	rackReq := NewRackRequest(r, "server", "port")

	req_writer.Write(rackReq.Bytes())

	if _, err = io.Copy(req_writer, r.Body); err != nil {
		log.Println("Error writing request body:", err)
	}

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
