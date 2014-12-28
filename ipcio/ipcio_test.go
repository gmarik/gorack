package ipcio

// run with -v to see debug output:
// go test -v gorack/sock_reader_test.go

import (
	"fmt"
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

func currentPath() string {
	_, filename, _, _ := runtime.Caller(0)
	return path.Dir(filename)
}

func TestIpcIo(t *testing.T) {
	log.Println("Creating Socket")
	pair, err := syscall.Socketpair(syscall.AF_LOCAL, syscall.SOCK_STREAM, 0)

	if err != nil {
		t.Error(err)
	}

	path := path.Join(currentPath(), "ipcio_test.rb")

	var remote, local = pair[0], pair[1]

	fmt.Println("Running:", path)

	cmd := exec.Command("ruby", path)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	// child process' FDs start from 3 (0, 1, 2)
	cmd.ExtraFiles = []*os.File{os.NewFile(uintptr(remote), "reader")}

	ch := make(chan string)
	expected := "hello"

	// TODO: investigate: sometimes takes too long
	go ipcEcho(local, ch, t)
	go func() { t.Error(cmd.Run()) }()

	// send data to proces
	ch <- expected

	// read data from process
	received := <-ch

	if !reflect.DeepEqual(received, expected) {
		t.Errorf("\nGot: %s\nExp: %s", received, expected)
	}
}

func ipcEcho(fd int, ch chan string, t *testing.T) {
	r, w, err := os.Pipe()

	if err != nil {
		t.Error(err)
	}

	err = SendIo(fd, r)
	w.Write([]byte(<-ch))
	w.Close()

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
