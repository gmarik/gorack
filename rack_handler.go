package gorack

import (
	"log"
	"net/http"
	"os"
	"os/exec"
	"syscall"

	"gmarik/gorack/ipcio"
)

var gorackRunner = "./ruby/bin/gorack"

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

	go spawnRackProcess(cmd)

	return &RackHandler{
		local_fd:   local,
		configPath: configPath,
	}
}

func (s *RackHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	res_reader, req_writer, err := s.sendIo()

	if err != nil {
		log.Println("[Error] creating resp/request pipes", err.Error())
		return
	}

	rackReq := NewRackRequest(r, "server", "port")
	go func() {
		if err := rackReq.WriteTo(req_writer); err != nil {
			log.Println("[Error] writing request body:", err)
		}
	}()

	resp := NewRackResponse(res_reader)
	if err := resp.WriteTo(w); err != nil {
		log.Println("[Error] writing response:", err.Error())
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
	}
}

// TODO: make it safe for concurrency
func (s *RackHandler) sendIo() (*os.File, *os.File, error) {

	var req_reader, req_writer *os.File
	var res_reader, res_writer *os.File
	var err error

	if req_reader, req_writer, err = os.Pipe(); err != nil {
		return nil, nil, err
	}

	if res_reader, res_writer, err = os.Pipe(); err != nil {
		return nil, nil, err
	}

	if err = ipcio.SendIo(s.local_fd, req_reader); err != nil {
		return nil, nil, err
	}

	if err = ipcio.SendIo(s.local_fd, res_writer); err != nil {
		return nil, nil, err
	}

	// close pipes once FDs sent to a child process
	// they'll still be open in the child process
	req_reader.Close()
	res_writer.Close()

	return res_reader, req_writer, nil
}

func spawnRackProcess(cmd *exec.Cmd) {
	cmd.Stdin = nil
	cmd.Stdout = NewLogWriter(os.Stdout, "", log.LstdFlags)
	cmd.Stderr = NewLogWriter(os.Stderr, "[StdErr]", log.LstdFlags)

	if err := cmd.Run(); err != nil {
		log.Fatal("Process '", cmd.Path, "' - failed to run:", err)
	}
}
