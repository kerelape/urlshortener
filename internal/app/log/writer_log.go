package log

import (
	"fmt"
	"io"
)

// WriterLog is a log that writes its messages to the io.Writer.
type WriterLog struct {
	info io.Writer
	err  io.Writer
}

// NewWriterLog returns a new WriterLog.
func NewWriterLog(info io.Writer, err io.Writer) *WriterLog {
	return &WriterLog{
		info: info,
		err:  err,
	}
}

// WriteInfo writes message to the Log with Info level.
func (log *WriterLog) WriteInfo(message string) {
	_, err := fmt.Fprint(log.info, message)
	if err != nil {
		panic(err)
	}
}

// WriteFailure writes message to the Log with Fail level.
func (log *WriterLog) WriteFailure(message string) {
	_, err := fmt.Fprint(log.err, message)
	if err != nil {
		panic(err)
	}
}
