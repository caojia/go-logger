package log

import (
	"fmt"
	"log"
	"os"
	"runtime/debug"

	"context"
)

var (
	loggerFlag = log.Ldate | log.Ltime
)

type Level int

/**
Java-Like Logger
*/
const (
	TRACE = Level(0)
	DEBUG = Level(1)
	INFO  = Level(2)
	WARN  = Level(3)
	ERROR = Level(4)

	CONTEXT_REQUEST_ID = "REQUEST_ID"
)

type Logger struct {
	*log.Logger
	level   Level
	context context.Context
}

func (gl *Logger) SetLevel(l Level) {
	gl.level = l
}

func (gl *Logger) WithContext(ctx context.Context) *Logger {
	ngl := new(Logger)
	*ngl = *gl
	ngl.context = ctx
	return ngl
}

func (gl *Logger) formatContext() string {
	if gl.context != nil {
		requestId, ok := gl.context.Value(CONTEXT_REQUEST_ID).(string)
		if !ok {
			requestId = "-"
		}
		return requestId
	} else {
		return "-"
	}
}

func (gl *Logger) msg(level Level, format string, args ...interface{}) string {
	str := ""
	if len(args) == 0 {
		str = format
	} else {
		str = fmt.Sprintf(format, args...)
	}
	return fmt.Sprintf("[%d] %s ", int(level), gl.formatContext()) + str
}

func (gl *Logger) Trace(format string, args ...interface{}) {
	if gl.level <= TRACE {
		gl.Logger.Println(gl.msg(TRACE, format, args...))
	}
}

func (gl *Logger) Debug(format string, args ...interface{}) {
	if gl.level <= DEBUG {
		gl.Logger.Println(gl.msg(DEBUG, format, args...))
	}
}

func (gl *Logger) Info(format string, args ...interface{}) {
	if gl.level <= INFO {
		gl.Logger.Println(gl.msg(INFO, format, args...))
	}
}

func (gl *Logger) Warn(format string, args ...interface{}) {
	if gl.level <= WARN {
		gl.Logger.Println(gl.msg(WARN, format, args...))
	}
}

// NOTE: 只有ERROR的时候会把Stack Trace打出来
func (gl *Logger) Error(format string, args ...interface{}) {
	if gl.level <= ERROR {
		stack := debug.Stack()
		text := fmt.Sprintf(format, args)
		text = text + "\nstack trace: " + string(stack)
		text = gl.msg(ERROR, text)
		gl.Logger.Println(text)
	}
}

func (gl *Logger) ErrorWithoutStack(format string, args ...interface{}) {
	if gl.level <= ERROR {
		text := fmt.Sprintf(format, args...)
		text = gl.msg(ERROR, text)
		gl.Logger.Println(text)
	}
}

func (gl *Logger) SetPrefix(prefix string) {
	// do nothing
}

func (gl *Logger) Println(args ...interface{}) {
	gl.Debug("", args...)
}

func (gl *Logger) Printf(format string, args ...interface{}) {
	gl.Debug(format, args...)
}

var CustomLogger = &Logger{
	Logger:  log.New(os.Stdout, "", loggerFlag),
	level:   DEBUG,
	context: context.Background(),
}
