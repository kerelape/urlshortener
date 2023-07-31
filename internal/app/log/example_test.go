package log_test

import (
	"os"

	"github.com/kerelape/urlshortener/internal/app/log"
)

func ExampleWriterLog_stdout() {
	l := log.NewWriterLog(os.Stdout, os.Stderr)
	l.WriteInfo("stdout info")
	l.WriteFailure("stderr failure")
}
