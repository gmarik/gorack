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
	"path"
	"reflect"
	"runtime"
	"syscall"
	"testing"
)

func TestSockIo(t *testing.T) {
	fmt.Println("Creating Socket")
	pair, err := syscall.Socketpair(syscall.AF_LOCAL, syscall.SOCK_STREAM, 0)

	var reader, writer = pair[0], pair[1]

	if err != nil {
		log.Fatal(err)
	}

	_, filename, _, _ := runtime.Caller(0)
	path := path.Join(path.Dir(filename), "ruby_process_test.rb")

	fmt.Println("Running:", path)

	cmd := exec.Command("ruby", path)

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

	// TODO: investigate: sometimes takes too long
	go processEcho(writer, response, expected, t)

	err = cmd.Wait()

	log.Println("Program exited with", err)

	received := <-response

	if !reflect.DeepEqual(received, expected) {
		t.Errorf("\nGot: %s\nExp: %s", received, expected)
	}
}

func processEcho(fd int, ch chan string, str string, t *testing.T) {
	r, w, err := os.Pipe()

	if err != nil {
		t.Error(err)
	}

	err = SendIo(fd, r)

	w.Write([]byte(str))

	if err != nil {
		t.Error(err)
	}

	defer close(ch)

	file, err := RecvIo(fd)

	if err != nil {
		t.Error(err)
	}

	data, err := ioutil.ReadAll(file)

	if err != nil {
		t.Error(err)
	}

	ch <- string(data)
}
