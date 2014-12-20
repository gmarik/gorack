package sockio

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

	// child process' FDs start from 3 (0, 1, 2)
	cmd.ExtraFiles = []*os.File{os.NewFile(uintptr(reader), "writer")}

	if err := cmd.Start(); err != nil {
		log.Fatal("Error running process", err)
	}

	go io.Copy(os.Stdout, out)
	go io.Copy(os.Stderr, errout)

	response := make(chan string)

	expected := "hello"

	sio := &SockIo{writer}
	// TODO: investigate: sometimes takes too long
	go processEcho(sio, response, expected, t)

	err = cmd.Wait()

	log.Println("Program exited with", err)

	received := <-response

	if !reflect.DeepEqual(received, expected) {
		t.Errorf("\nGot: %s\nExp: %s", received, expected)
	}
}

func processEcho(s *SockIo, ch chan string, str string, t *testing.T) {
	r, w, err := os.Pipe()

	if err != nil {
		t.Error(err)
	}

	err = s.SendIo(r)

	w.Write([]byte(str))

	if err != nil {
		t.Error(err)
	}

	defer close(ch)

	file, err := s.RecvIo()

	if err != nil {
		t.Error(err)
	}

	data, err := ioutil.ReadAll(file)

	if err != nil {
		t.Error(err)
	}

	ch <- string(data)
}
