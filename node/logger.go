package node

type Logger interface {
	Printf(format string, v ...interface{})
	Println(v ...interface{})
}

func NewNullLogger() Logger {
	return &nullLogger{}
}

type nullLogger struct {
}

func (log *nullLogger) Printf(format string, v ...interface{}) {
	// stub
}

func (log *nullLogger) Println(v ...interface{}) {
	// stub
}
