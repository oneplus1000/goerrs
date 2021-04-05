package goerrs

import (
	"fmt"
	"strconv"
	"testing"
)

var errX = New("error xxxx not found")

func TestCallers(t *testing.T) {
	err := sub3()
	fmt.Printf("%+v\n", err)
	//PrintCallStack(err)
	/*var errtmp = err
	i := 0
	for errtmp != nil {
		fmt.Printf("%d) %s\n", i, errtmp.Error())
		errtmp = errors.Unwrap(errtmp)
		i++
	}*/

}

func sub3() error {
	err := sub2()
	if err != nil {
		return Wrapf("%v sub2 fail", err)
		//return fmt.Errorf("%w", err)
	}
	return nil
}

func sub2() error {
	err := sub1()
	if err != nil {
		//return Wrapf("%+v in sub2", err)
		return fmt.Errorf("%w cccc", err)
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
