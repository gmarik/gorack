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

log "Creating socket"
r,w = IO.pipe

log "sending some data"
sock.send_io(r)

w.write("hello")
w.close
r.close
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

	fd := os.NewFile(uintptr(pair[0]), "writer")

	// child process' FDs start from 3 (0, 1, 2)
	cmd.ExtraFiles = []*os.File{fd}

	if err := cmd.Start(); err != nil {
		log.Fatal("Error running process", err)
	}

	fd.Close()

	go io.Copy(os.Stdout, out)
	go io.Copy(os.Stderr, errout)

	quit := make(chan struct{})

	go func() {
		defer close(quit)
		reader := os.NewFile(uintptr(pair[1]), "reader")

		// from: https://github.com/ftrvxmtrx/fd/blob/master/fd.go
		// recvmsg
		buf := make([]byte, syscall.CmsgSpace(4))
		_, _, _, _, err = syscall.Recvmsg(int(reader.Fd()), nil, buf, syscall.MSG_WAITALL)

		if err != nil {
			log.Fatal(err)
		}

		// parse control msgs
		var msgs []syscall.SocketControlMessage
		msgs, err = syscall.ParseSocketControlMessage(buf)

		if err != nil {
			log.Fatal(err)
		}

		filenames := []string{"result"}

		// convert fds to files
		res := make([]*os.File, 0, len(msgs))

		for _, msg := range msgs {
			var fds []int
			fds, err = syscall.ParseUnixRights(&msg)

			if err != nil {
				log.Println(err)
				continue
			}

			for fi, fd := range fds {
				var filename string
				if fi < len(filenames) {
					filename = filenames[fi]
				}

				res = append(res, os.NewFile(uintptr(fd), filename))
			}
		}

		file := res[0]

		log.Println(res)

		fmt.Println("Reading data")

		data := make([]byte, 100)

		n, err := file.Read(data)

		if err != nil {
			log.Fatal(err)
		}

		received, expected := data[0:n], []byte("hello")

		if !reflect.DeepEqual(received, expected) {
			t.Errorf("\nGot: %s\nExp: %s", received, expected)
		}
	}()

	err = cmd.Wait()

	log.Println("Program exited with", err)

	// writer.Close()
	// reader.Close()
	<-quit
}
