package gorack

import (
	"io"
	"log"
)

type LogWriter struct {
	logger *log.Logger
	prefix string
}

func NewLogWriter(w io.Writer, prefix string, flags int) *LogWriter {
	return &LogWriter{log.New(w, "", flags), prefix}
}

func (l *LogWriter) Write(data []byte) (int, error) {
	l.logger.Printf("%s%s", l.prefix, data)
	return len(data), nil
}
