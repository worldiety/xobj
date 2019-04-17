package xobj

import (
	"fmt"
	"time"
)

var logger Logger = &defaultLogger{}

// SetLogger sets the module/package level logger
func SetLogger(log Logger) {
	logger = log
}

// Fields is just a type alias to avoid verbosity while using
type Fields = map[string]interface{}

// A Logger is just a simple interface which can be easily satisfied by any implementor
type Logger interface {
	Info(fields Fields)
}

// the default logger just prints as json into stdout
type defaultLogger struct {
}

func (l *defaultLogger) Info(fields map[string]interface{}) {
	fields["ts"] = time.Now().String()
	fields["level"] = "INFO"
	fmt.Println(Object(fields).String())
}
