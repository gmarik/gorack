package sockio

import (
	"os"
	"syscall"
)

type SockIo struct {
	Fd int
}

func (s *SockIo) RecvIo() (*os.File, error) {
	// # TODO: why 4?
	buf := make([]byte, syscall.CmsgSpace(4))
	_, _, _, _, err := syscall.Recvmsg(s.Fd, nil, buf, syscall.MSG_WAITALL)

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
func (s *SockIo) SendIo(file *os.File) error {
	rights := syscall.UnixRights(int(file.Fd()))
	return syscall.Sendmsg(s.Fd, nil, rights, nil, 0)
}
