package log

import (
	"fmt"
	"time"
)

type FormattedLog struct {
	Origin     Log
	TimeLayout string
}

func NewFormattedLog(origin Log, timeLayout string) *FormattedLog {
	return &FormattedLog{
		Origin:     origin,
		TimeLayout: timeLayout,
	}
}

const Format = "[%s] (%s) %s\r\n"

func (log *FormattedLog) WriteInfo(message string) {
	time := time.Now().Format(log.TimeLayout)
	log.Origin.WriteInfo(fmt.Sprintf(Format, time, "INFO", message))
}

func (log *FormattedLog) WriteFailure(message string) {
	time := time.Now().Format(log.TimeLayout)
	log.Origin.WriteFailure(fmt.Sprintf(Format, time, "FAIL", message))
}
