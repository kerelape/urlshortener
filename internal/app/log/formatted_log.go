package log

import (
	"fmt"
	"time"
)

// FormattedLog is a Log that formats the messages.
type FormattedLog struct {
	origin     Log
	timeLayout string
}

// NewFormattedLog returns a new FormattedLog.
func NewFormattedLog(origin Log, timeLayout string) *FormattedLog {
	return &FormattedLog{
		origin:     origin,
		timeLayout: timeLayout,
	}
}

const format = "[%s] (%s) %s\r\n"

// WriteInfo writes message to the Log with Info level.
func (log *FormattedLog) WriteInfo(message string) {
	time := time.Now().Format(log.timeLayout)
	log.origin.WriteInfo(fmt.Sprintf(format, time, "INFO", message))
}

// WriteFailure writes message to the Log with Fail level.
func (log *FormattedLog) WriteFailure(message string) {
	time := time.Now().Format(log.timeLayout)
	log.origin.WriteFailure(fmt.Sprintf(format, time, "FAIL", message))
}
