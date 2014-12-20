package sockio

import (
	"os"
	"syscall"
)

func RecvIo(fd int) (*os.File, error) {
	// # TODO: why 4?
	buf := make([]byte, syscall.CmsgSpace(4))
	_, _, _, _, err := syscall.Recvmsg(fd, nil, buf, syscall.MSG_WAITALL)

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
func SendIo(fd int, file *os.File) error {
	rights := syscall.UnixRights(int(file.Fd()))
	return syscall.Sendmsg(fd, nil, rights, nil, 0)
}
