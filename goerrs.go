package goerrs

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"runtime"
	"strings"
)

type basicErr struct {
	isNew  bool
	msg    string
	caller trace
	err    error
}

func PrintCallStack(err error) {
	fmt.Print(CallStack(err))
}

func CallStack(err error) string {
	var errTmp = err
	var buff bytes.Buffer
	i := 1
	for errTmp != nil {
		if val, ok := errTmp.(*basicErr); ok {
			msg := strings.TrimSpace(val.msg)
			if msg != "" {
				msg = msg + ", "
			}
			buff.WriteString(fmt.Sprintf("%d) %s%s %s:%d\n", i, msg, val.caller.funcName, val.caller.fileName, val.caller.line))
		} else {
			buff.WriteString(fmt.Sprintf("%d) %v\n", i, errTmp))
		}
		i++
		errTmp = errors.Unwrap(errTmp)
	}
	return buff.String()
}

func (b *basicErr) Error() string {
	return b.msg
}

func (b *basicErr) Unwrap() error {
	return b.err
}

func (b *basicErr) Format(f fmt.State, verb rune) {
	io.WriteString(f, b.Error())
}

func New(msg string) error {
	return &basicErr{
		isNew:  true,
		msg:    msg,
		caller: caller(),
		err:    nil,
	}
}

//Wrapf wrap error message with call stack
//(format wrap error message use %v instead of %w)
func Wrapf(format string, args ...interface{}) error {

	var err error = nil
	for _, a := range args {
		if val, ok := a.(error); ok {
			err = val
			break
		}
	}
	msg := fmt.Sprintf(format, args...)
	return &basicErr{
		isNew:  false,
		msg:    msg,
		caller: caller(),
		err:    err,
	}
}

func caller() trace {
	fpcs := make([]uintptr, 1)
	n := runtime.Callers(3, fpcs)
	if n == 0 {
		return emptyTrace()
	}

	caller := runtime.FuncForPC(fpcs[0] - 1)
	if caller == nil {
		return emptyTrace()
	}

	file, line := caller.FileLine(fpcs[0] - 1)
	return trace{
		valid:    true,
		funcName: caller.Name(),
		fileName: file,
		line:     line,
	}
}

type trace struct {
	valid    bool
	funcName string
	fileName string
	line     int
}

func emptyTrace() trace {
	return trace{
		valid: false,
	}
}
