package goerrs

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"runtime"
	"strings"
)

type WrapErrWithCallStacker interface {
	Caller() Trace
}

type wrapErrWithCallStack struct {
	msg    string
	caller Trace
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
		if val, ok := errTmp.(WrapErrWithCallStacker); ok {
			msg := strings.TrimSpace(errTmp.Error())
			if msg != "" {
				msg = msg + ", "
			}
			buff.WriteString(fmt.Sprintf("%d) %s%s %s:%d\n", i, msg, val.Caller().FuncName, val.Caller().FileName, val.Caller().Line))
		} else {
			buff.WriteString(fmt.Sprintf("%d) %v\n", i, errTmp))
		}
		i++
		errTmp = errors.Unwrap(errTmp)
	}
	return buff.String()
}

func (w *wrapErrWithCallStack) Caller() Trace {
	return w.caller
}

func (w *wrapErrWithCallStack) Error() string {
	return w.msg
}

func (w *wrapErrWithCallStack) Unwrap() error {
	return w.err
}

func (w *wrapErrWithCallStack) Format(f fmt.State, verb rune) {
	io.WriteString(f, w.Error())
}

//WrapCallStack wrap error message with call stack
func WrapCallStack(err error) error {
	return &wrapErrWithCallStack{
		msg:    err.Error(),
		caller: caller(),
		err:    err,
	}
}

//WrapCallStackf wrap error message with call stack
//(format wrap error message use %v instead of %w)
func WrapCallStackf(format string, args ...interface{}) error {

	var err error = nil
	for _, a := range args {
		if val, ok := a.(error); ok {
			err = val
			break
		}
	}
	msg := fmt.Sprintf(format, args...)
	return &wrapErrWithCallStack{
		msg:    msg,
		caller: caller(),
		err:    err,
	}
}

func caller() Trace {
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
	return Trace{
		Valid:    true,
		FuncName: caller.Name(),
		FileName: file,
		Line:     line,
	}
}

type Trace struct {
	Valid    bool
	FuncName string
	FileName string
	Line     int
}

func emptyTrace() Trace {
	return Trace{
		Valid: false,
	}
}
