package errs

import (
	"fmt"
	"runtime"
)

type Logger interface {
	Log(message string)
}
type Error interface {
	error
	Log(logger Logger, additionalContext ...interface{}) Error
	WithStack() Error
	Unwrap() error
	Message() string
	Code() int
	Data() interface{}
	Errors() interface{}
	String() string
}

type appError struct {
	message string
	code    int
	stack   []byte
	err     error
	data    interface{}
	errors  interface{}
}

func New(message string, code ...int) Error {
	err := &appError{message: message}
	if len(code) > 0 {
		err.code = code[0]
	}
	return err
}

func Wrap(err error, message string, code ...int) Error {
	appErr := &appError{err: err, message: message}
	if len(code) > 0 {
		appErr.code = code[0]
	}
	return appErr
}

func (e *appError) Error() string {
	if e.err != nil {
		return fmt.Sprintf("%s: %v", e.message, e.err)
	}
	return e.message
}

func (e *appError) Log(logger Logger, additionalContext ...interface{}) Error {
	e = e.WithStack().(*appError)
	logMessage := fmt.Sprintf("[%s] %s", e.Error(), additionalContext)
	if len(additionalContext) > 0 {
		logMessage += fmt.Sprintf(" - Additional Context: %v", additionalContext)
	}
	if e.stack != nil {
		logMessage += "\nStack Trace:\n" + string(e.stack)
	}
	logger.Log(logMessage)
	return e
}

func (e *appError) WithStack() Error {
	buf := make([]byte, 1024)
	n := runtime.Stack(buf, false)
	e.stack = buf[:n]
	return e
}

func (e *appError) Unwrap() error {
	return e.err
}

func (e *appError) String() string {
	str := fmt.Sprintf("Error: %s", e.Error())
	if e.stack != nil {
		str += "\nStack Trace:\n" + string(e.stack)
	}
	return str
}

func (e *appError) Message() string {
	return e.message
}

func (e *appError) Code() int {
	return e.code
}
func (e *appError) Data() interface{} {
	return e.data
}

func (e *appError) Errors() interface{} {
	return e.errors
}
