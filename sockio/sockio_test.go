package gorack

// run with -v to see debug output:
// go test -v gorack/sock_reader_test.go

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"reflect"
	"syscall"
	"testing"
)

const rubyProcess = `#!/usr/bin/env ruby
require 'socket'

def log(msg)
	puts "[RUBY] #{msg}"
end

# passed from parent process
sock = UNIXSocket.for_fd(3)

log "receiving socket"
r = sock.recv_io

log "creating proxy pipe"
ior, iow = IO.pipe

log "sending the pipe"
sock.send_io(ior)

log "copying stream"
IO.copy_stream(r, iow)
`

var rubyProcessFile *os.File

func writeContent() {
	var err error

	rubyProcessFile, err = ioutil.TempFile("/tmp/", "rubySock")

	if err != nil {
		log.Fatal(err)
	}

	if _, err := rubyProcessFile.Write([]byte(rubyProcess)); err != nil {
		log.Fatal(err)
	}
}

func TestSocketReading(t *testing.T) {
	writeContent()
	defer rubyProcessFile.Close()

	fmt.Println("Creating Socket")
	pair, err := syscall.Socketpair(syscall.AF_LOCAL, syscall.SOCK_STREAM, 0)

	var reader, writer = pair[0], pair[1]

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Reading from:", rubyProcessFile.Name())

	cmd := exec.Command("ruby", rubyProcessFile.Name())

	out, _ := cmd.StdoutPipe()
	errout, _ := cmd.StderrPipe()

	if err != nil {
		log.Fatal(err)
	}

	fd := os.NewFile(uintptr(reader), "writer")

	// child process' FDs start from 3 (0, 1, 2)
	cmd.ExtraFiles = []*os.File{fd}

	if err := cmd.Start(); err != nil {
		log.Fatal("Error running process", err)
	}

	fd.Close()

	go io.Copy(os.Stdout, out)
	go io.Copy(os.Stderr, errout)

	response := make(chan string)

	expected := "hello"

	go processEcho(writer, response, expected, t)

	err = cmd.Wait()

	log.Println("Program exited with", err)

	received := <-response

	if !reflect.DeepEqual(received, expected) {
		t.Errorf("\nGot: %s\nExp: %s", received, expected)
	}
}

func processEcho(writer int, ch chan string, str string, t *testing.T) {
	r, w, err := os.Pipe()

	if err != nil {
		t.Error(err)
	}

	err = SendIo(writer, r)

	w.Write([]byte(str))

	if err != nil {
		t.Error(err)
	}

	defer close(ch)

	file, err := RecvIo(writer)

	if err != nil {
		t.Error(err)
	}

	data, err := ioutil.ReadAll(file)

	if err != nil {
		t.Error(err)
	}

	ch <- string(data)
}

func RecvIo(socket_fd int) (*os.File, error) {
	// # TODO: why 4?
	buf := make([]byte, syscall.CmsgSpace(4))
	_, _, _, _, err := syscall.Recvmsg(socket_fd, nil, buf, syscall.MSG_WAITALL)

	if err != nil {
		return nil, err
	}

	msgs, err := syscall.ParseSocketControlMessage(buf)

	if err != nil {
		return nil, err
	}

	fds, err := syscall.ParseUnixRights(&msgs[0])

	return os.NewFile(uintptr(fds[0]), ""), nil
}

// from: https://github.com/ftrvxmtrx/fd/blob/master/fd.go
func SendIo(socket_fd int, file *os.File) error {
	rights := syscall.UnixRights(int(file.Fd()))
	return syscall.Sendmsg(socket_fd, nil, rights, nil, 0)
}
