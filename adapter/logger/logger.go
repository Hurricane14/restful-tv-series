package logger

type Logger interface {
	Debugf(string, ...any)
	Errorf(string, ...any)
	Fatalf(string, ...any)
	Infof(string, ...any)
	Printf(string, ...any)
	Warnf(string, ...any)
}
