package log

type Log interface {
	WriteInfo(message string)
	WriteFailure(message string)
}
