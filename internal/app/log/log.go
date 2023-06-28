package log

// Log is a log.
type Log interface {
	// WriteInfo writes message to the Log with Info level.
	WriteInfo(message string)

	// WriteFailure writes message to the Log with Fail level.
	WriteFailure(message string)
}
