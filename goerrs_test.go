package goerrs

import (
	"errors"
	"strconv"
	"strings"
	"testing"
)

var errX = errors.New("error xxxx not found")

func TestCallers(t *testing.T) {
	err := sub3()
	if val, ok := err.(WrapErrWithCallStacker); ok {
		if strings.TrimSpace(val.Caller().FileName) == "" || strings.TrimSpace(val.Caller().FuncName) == "" {
			t.Error("Caller is Empty")
		}
	} else {
		t.Error("error not WrapErrWithCallStacker")
	}
}

func sub3() error {
	err := sub2()
	if err != nil {
		return WrapCallStackf("%v sub2 fail", err)
		//return fmt.Errorf("%w", err)
	}
	return nil
}

func sub2() error {
	err := sub1()
	if err != nil {
		//return Wrapf("%+v in sub2", err)
		return WrapCallStackf("%v cccc", err)
	}
	return nil
}

func sub1() error {
	_, err := strconv.Atoi("x")
	if err != nil {
		return errX //Wrap(err, "")
	}
	return nil
}
