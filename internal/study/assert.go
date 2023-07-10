package study

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/go-test/deep"
)

type AssertFailure struct {
	Expected interface{}
	Actual   interface{}
}

var errAssert = errors.New("")

func Assert(expected interface{}, actual interface{}) error {
	if actualErr, ok := actual.(error); ok {
		actual = actualErr.Error()
	}

	if expectedErr, ok := expected.(error); ok {
		expected = expectedErr.Error()
	}

	expectedString, expectedIsString := expected.(string)
	actualString, actualIsString := actual.(string)
	if expectedIsString && actualIsString {
		if expectedString != actualString {
			b, _ := json.Marshal(AssertFailure{
				Expected: expectedString,
				Actual:   actualString,
			})

			return errors.New(string(b))
		}

		return nil
	}

	if diff := deep.Equal(expected, actual); diff != nil {
		return fmt.Errorf("%w%v", errAssert, diff)
	}

	return nil
}
