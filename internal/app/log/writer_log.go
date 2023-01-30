package log

import (
	"fmt"
	"io"
)

type WriterLog struct {
	Info io.Writer
	Err  io.Writer
}

func NewWriterLog(info io.Writer, err io.Writer) *WriterLog {
	return &WriterLog{
		Info: info,
		Err:  err,
	}
}

func (log *WriterLog) WriteInfo(message string) {
	var _, err = fmt.Fprint(log.Info, message)
	if err != nil {
		panic(err)
	}
}

func (log *WriterLog) WriteFailure(message string) {
	var _, err = fmt.Fprint(log.Err, message)
	if err != nil {
		panic(err)
	}
}
