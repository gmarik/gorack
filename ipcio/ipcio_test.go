package ipcio

// run with -v to see debug output:
// go test -v gorack/sock_reader_test.go

import (
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

	log.Println("Running:", path)

	cmd := exec.Command("ruby", path)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	// child process' FDs start from 3 (0, 1, 2)
	cmd.ExtraFiles = []*os.File{os.NewFile(uintptr(remote), "reader")}

	ch := make(chan string)
	expected := "hello"

	go ipcEcho(local, ch, t)           // sends data to a child process
	go func() { t.Error(cmd.Run()) }() // runs the child process

	// send data to proces
	ch <- expected

	// read data from process
	received := <-ch

	if !reflect.DeepEqual(received, expected) {
		t.Errorf("\nGot: %s\nExp: %s", received, expected)
	}
}

func ipcEcho(fd int, ch chan string, t *testing.T) {
	var r, w *os.File
	var err error

	if r, w, err = os.Pipe(); err != nil {
		t.Error(err)
	}

	// send pipe's reader to a process
	// so the process can read data from it
	if err = SendIo(fd, r); err != nil {
		t.Error(err)
	}

	w.Write([]byte(<-ch))
	w.Close()

	var file *os.File
	var data []byte

	// receives a reader to read reply from the process
	if file, err = RecvIo(fd); err != nil {
		t.Error(err)
	}

	if data, err = ioutil.ReadAll(file); err != nil {
		t.Error(err)
	}

	ch <- string(data)
	close(ch)
}
