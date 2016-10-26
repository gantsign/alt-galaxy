package logging

type Log interface {
	Progressf(format string, a ...interface{})

	Errorf(format string, a ...interface{})
}
